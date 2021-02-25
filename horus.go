package horus

import (
	"encoding/json"
	"fmt"
	"github.com/ichtrojan/horus/models"
	"github.com/ichtrojan/horus/storage"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func Watch(next func(http.ResponseWriter, *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := storage.Connect()

		if err != nil {
			log.Fatal(err)
		}

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		header, err := json.Marshal(r.Header)

		if err != nil {
			log.Fatal(err)
		}

		requestBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Fatal(err)
		}

		req := models.Request{
			ResponseBody:  "",
			ResposeStatus: 200,
			RequestBody:   requestBody,
			Path:          r.RequestURI,
			Headers:       header,
			Method:        r.Method,
			Host:          r.Host,
			Ipadress:      ip,
		}

		write := request.Create(&req)

		if write.RowsAffected != 1 {
			log.Fatal("unable to log request")
		}

		next(w, r)
	}
}

func Serve(port string) error {
	http.HandleFunc("/horus", func(w http.ResponseWriter, r *http.Request) {
		var req models.Request

		request, err := storage.Connect()

		if err != nil {
			_ = fmt.Errorf("%v", err)
		}

		request.First(&req)
	})

	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}

package utils

import (
	"encoding/json"
	http "github.com/Noooste/fhttp"
	"io/ioutil"
	"strings"
)

func GetRequestBody(request *http.Request) []byte {
	var body []byte

	encoding := request.Header.Get("Content-Encoding")

	bodyBytes, err := ioutil.ReadAll(request.Body)

	if err != nil {
		body = []byte{}
	} else if encoding != "" {
		body = DecompressBody(bodyBytes, encoding)
	} else {
		body = bodyBytes
	}

	return body
}

func GetResponseBody(response *http.Response) []byte {
	defer response.Body.Close()

	encoding := response.Header.Get("content-encoding")

	bodyBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return []byte{}
	}

	return DecompressBody(bodyBytes, encoding)
}

type SuccessReturn struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type SuccessInitSensor struct {
	Success bool   `json:"success"`
	Id      string `json:"id"`
}

func ReturnWeb(success any) []byte {
	dumped, err := json.Marshal(success)

	if err != nil {
		return []byte("{}")
	}

	return dumped
}

func ReturnWebError(error string) []byte {
	return ReturnWeb(SuccessReturn{
		Success: false,
		Error:   error,
	})
}

func GetIP(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

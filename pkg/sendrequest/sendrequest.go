package sendrequest

import (
	"bytes"
	"log"
	"net/http"
)

func sendRequest(addr string, method string, body []byte) *http.Response {
	request, err := http.NewRequest(method, addr, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("exit: %s\n", err.Error())
		return nil
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("exit: %s\n", err.Error())
		return nil
	}

	return response
}

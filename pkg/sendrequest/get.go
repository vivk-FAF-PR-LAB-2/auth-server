package sendrequest

import "net/http"

func Get(addr string, body []byte) *http.Response {
	return sendRequest(addr, "GET", body)
}

package sendrequest

import "net/http"

func Post(addr string, body []byte) *http.Response {
	return sendRequest(addr, "POST", body)
}

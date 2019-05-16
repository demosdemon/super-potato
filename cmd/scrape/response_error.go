package scrape

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type ResponseError struct {
	StatusCode int
	Status     string
	Body       []byte
}

func (r ResponseError) Error() string {
	return fmt.Sprintf("HTTP %03d: %s - %s", r.StatusCode, r.Status, string(r.Body))
}

func NewResponseError(r *http.Response) error {
	body, _ := ioutil.ReadAll(r.Body)
	return ResponseError{
		StatusCode: r.StatusCode,
		Status:     r.Status,
		Body:       body,
	}
}

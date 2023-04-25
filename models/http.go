package models

// HttpResponse is a struct that is used to send
// a response to the client
type HttpResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Err     error  `json:"error"`
}

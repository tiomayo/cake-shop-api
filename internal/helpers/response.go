package helpers

type JSONResponse struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

type ErrorObject struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

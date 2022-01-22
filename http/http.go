package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	pa "github.com/Lambels/patrickarvatu.com"
)

// error code -> http error code.
var httpCodes = map[string]int{
	pa.ENOTFOUND:       http.StatusNotFound,
	pa.ECONFLICT:       http.StatusConflict,
	pa.EINTERNAL:       http.StatusInternalServerError,
	pa.EINVALID:        http.StatusBadRequest,
	pa.ENOTIMPLEMENTED: http.StatusNotImplemented,
	pa.EUNAUTHORIZED:   http.StatusUnauthorized,
}

// getErrorCode maps the code to an http code if possible or returns 500.
func getErrorCode(code string) int {
	v, ok := httpCodes[code]
	if !ok {
		return http.StatusInternalServerError
	}
	return v
}

// errorResponse represents an http error response used by send error.
type errorResponse struct {
	Code  int    `json:"code"`
	Trace string `json:"trace"`
}

// SendError sends an error over http.
// maps system error codes to http error codes.
func SendError(w http.ResponseWriter, r *http.Request, err error) {
	code, message := pa.ErrorCode(err), pa.ErrorMessage(err)

	if code == pa.EINTERNAL {
		LogError(r, err)
	}

	v := getErrorCode(code)
	w.WriteHeader(v)
	SendJSON(w, errorResponse{
		Code:  v,
		Trace: message,
	})
}

// LogError logs an error with format '[HTTP] error: {method} {path}: {error message}'
func LogError(r *http.Request, err error) {
	log.Printf("[HTTP] error: %s %s: %s\n", r.URL.Path, r.Method, err)
}

// SendJSON sends json: data over http.
func SendJSON(w io.Writer, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}

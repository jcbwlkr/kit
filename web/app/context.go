// Package app provides application support for context and MongoDB access.
// Current Status Codes:
//		200 OK           : StatusOK                  : Call is success and returning data.
//		204 No Content   : StatusNoContent           : Call is success and returns no data.
//		400 Bad Request  : StatusBadRequest          : Invalid post data (syntax or semantics).
//		401 Unauthorized : StatusUnauthorized        : Authentication failure.
//		404 Not Found    : StatusNotFound            : Invalid URL or identifier.
//		500 Internal     : StatusInternalServerError : Application specific beyond scope of user.
package app

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ardanlabs/kit/auth"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
)

// Invalid describes a validation error belonging to a specific field.
type Invalid struct {
	Fld string `json:"field_name"`
	Err string `json:"error"`
}

// jsonError is the response for errors that occur within the API.
type jsonError struct {
	Error  string    `json:"error"`
	Fields []Invalid `json:"fields,omitempty"`
}

//==============================================================================

// Context contains data associated with a single request.
type Context struct {
	DB *db.DB
	http.ResponseWriter
	Request   *http.Request
	Params    map[string]string
	SessionID string
	User      *auth.User
	Status    int
}

// Error handles all error responses for the API.
func (c *Context) Error(err error) {
	switch err {
	case ErrNotFound:
		c.RespondError(err.Error(), http.StatusNotFound)
	case ErrInvalidID:
		c.RespondError(err.Error(), http.StatusBadRequest)
	case ErrValidation:
		c.RespondError(err.Error(), http.StatusBadRequest)
	default:
		c.RespondError(err.Error(), http.StatusInternalServerError)
	}
}

// Respond sends JSON to the client.
// If code is StatusNoContent, v is expected to be nil.
func (c *Context) Respond(data interface{}, code int) {
	log.User(c.SessionID, "api : Respond", "Started : Code[%d]", code)

	c.Status = code

	if code == http.StatusNoContent {
		c.WriteHeader(http.StatusNoContent)
		return
	}

	// Set application default header values.
	c.Header().Set("Content-Type", "application/json")

	// Load any user defined header values.
	if app.userHeaders != nil {
		for key, value := range app.userHeaders {
			log.User("startup", "Init", "Setting user headers : %s:%s", key, value)
			c.Header().Set(key, value)
		}
	}

	c.WriteHeader(code)

	// Marshal the data into a JSON string.
	jsonData, err := json.Marshal(data)
	if err != nil {
		jsonData = []byte("{}")
	}

	// Look for a JSONP marker
	if cb := c.Request.URL.Query().Get("callback"); cb != "" {

		// We need to wrap the result in a function call.
		// callback_value({"data_1": "hello world", "data_2": ["the","sun","is","shining"]});
		io.WriteString(c, cb+"("+string(jsonData)+")")

		log.User(c.SessionID, "api : Respond", "Completed")
		return
	}

	// We can send the result straight through.
	io.WriteString(c, string(jsonData))

	log.User(c.SessionID, "api : Respond", "Completed")
}

// RespondInvalid sends JSON describing field validation errors.
func (c *Context) RespondInvalid(fields []Invalid) {
	v := jsonError{
		Error:  "field validation failure",
		Fields: fields,
	}
	c.Respond(v, http.StatusBadRequest)
}

// RespondError sends JSON describing the error
func (c *Context) RespondError(error string, code int) {
	c.Respond(jsonError{Error: error}, code)
}

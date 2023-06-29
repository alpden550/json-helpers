package json_helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// JSONResponse is the type used to send JSON back
type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool is the main type of the package with json methods
type Tool struct {
	MaxJSONSize          int
	AllowedUnknownFields bool
}

// ReadJSONBody tries to read the body of http request and converts it the JSON object.
// data parameter is expected to be a pointer to object to read json data into it.
func (tool *Tool) ReadJSONBody(writer http.ResponseWriter, request *http.Request, data any) error {
	if request.Header.Get("Content-Type") != "" {
		contentType := request.Header.Get("Content-Type")
		if strings.ToLower(contentType) != "application/json" {
			return errors.New("the Content-Type header is not application/json")
		}
	}

	var maxBytes = 1024 * 1024 // 1 megabyte
	if tool.MaxJSONSize != 0 {
		maxBytes = tool.MaxJSONSize
	}
	request.Body = http.MaxBytesReader(writer, request.Body, int64(maxBytes))

	dec := json.NewDecoder(request.Body)

	if !tool.AllowedUnknownFields {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(&data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("body contains incorrect JSON type for field %q at offset %d", unmarshalTypeError.Field, unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			return fmt.Errorf("error unmarshalling json: %s", err.Error())

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// WriteJSON send back json response with the given status code and headers to the client.
// The Content-Type header is set to application/json.
func (tool *Tool) WriteJSON(writer http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			writer.Header()[key] = value
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_, err = writer.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// WriteErrorJSON send back JSON error response with given error and status code, if passed.
// Default status code is 400.
func (tool *Tool) WriteErrorJSON(writer http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := JSONResponse{
		Error:   true,
		Message: err.Error(),
	}

	return tool.WriteJSON(writer, statusCode, payload)
}

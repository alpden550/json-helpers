# Go Json-Helpers

Go JSON helpers to read and write JSON  from request to response writer.


## How to install

```
go get -u github.com/alpden550/json-helpers
```

## Working with JSON

```go
package main

import helpers "github.com/alpden550/json-helpers"

// JSONPayload is the type for JSON data that we receive from post request
type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// read json into requestPayload
var requestPayload JSONPayload
if err := helpers.ReadJSONBody(writer, request, &requestPayload); err != nil {
    return err
}

payload := jsonResponse{
		Error:   false,
		Message: "message",
	}
	
// send error json message if error was happened
user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		err = helpers.WriteErrorJSON(writer, errors.New("not found user"), http.StatusBadRequest)
		return
	}
	
// write and send response back as JSON 	
_ = helpers.WriteJSON(writer, http.StatusOK, payload)
```

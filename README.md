# Go Json-Helpers

Go JSON helpers to read and write JSON  from request to response writer.


## How to install

```
go get -u github.com/alpden550/json-helpers
```

## Working with JSON

Using in a http handler, for example:

```go
package main

import helpers "github.com/alpden550/json-helpers"

// JSONPayload is the type for JSON data that we receive from post request
type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func Handler(writer http.ResponseWriter, request *http.Request) {
	// create tool variable and initialize it
	tool := helpers.Tool{}
	
	// read json into requestPayload
	var requestPayload JSONPayload
	if err := tool.ReadJSONBody(writer, request, &requestPayload); err != nil {
		return err
	}
	
	// send error json message if error was happened
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		err = tool.WriteErrorJSON(writer, errors.New("not found user"))
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "message",
	}
	// write and send response back as JSON 	
	_ = tool.WriteJSON(writer, http.StatusOK, responsePayload)
}
```

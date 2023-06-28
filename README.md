# json-helpers

Go JSON helpers to read and write JSON  from request to response writer.


## How to install

```
go get -u github.com/alpden550/json-helpers
```

## How to use

```
import helpers "github.com/alpden550/json-helpers"

if err := helpers.ReadJSONBody(writer, request, &data); err != nil {
    return err
}

payload := jsonResponse{
		Error:   false,
		Message: "message",
	}

user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		err = helpers.WriteErrorJSON(writer, errors.New("not found user"), http.StatusBadRequest)
		return
	}

_ = helpers.WriteJSON(writer, http.StatusOK, payload)
```

package json_helpers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkTool_ReadJSONBody(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var tool Tool
		var decodedJSON struct {
			Foo string `json:"foo"`
		}

		request, _ := http.NewRequest("POST", "/", bytes.NewReader([]byte(`{"foo": "bar"}`)))
		rr := httptest.NewRecorder()
		_ = tool.ReadJSONBody(rr, request, &decodedJSON)
	}
}

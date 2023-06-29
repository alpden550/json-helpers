package json_helpers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTool_ReadJSONBodySuccessfully(t *testing.T) {
	req := require.New(t)

	var tests = map[string]struct {
		json         string
		expected     any
		maxSize      int
		allowUnknown bool
		contentType  string
	}{
		"valuable json":                           {json: `{"foo": "bar"}`, allowUnknown: false, expected: "bar"},
		"allowed unknown fields":                  {json: `{"foo": "bar", "bar": "foo"}`, allowUnknown: true, expected: "bar"},
		"allowed unknown fields without valuable": {json: `{"bar": "foo"}`, allowUnknown: true, expected: ""},
	}

	for name, testCase := range tests {
		tool := Tool{
			MaxJSONSize:          testCase.maxSize,
			AllowedUnknownFields: testCase.allowUnknown,
		}

		var decodedJSON struct {
			Foo string `json:"foo"`
		}

		t.Run(name, func(t *testing.T) {
			request, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(testCase.json)))
			rr := httptest.NewRecorder()
			err = tool.ReadJSONBody(rr, request, &decodedJSON)

			req.Equal(testCase.expected, decodedJSON.Foo)
			req.NoError(err)
		})
	}
}

func TestTool_ReadJSONBodyWithErrors(t *testing.T) {
	req := require.New(t)

	var tests = map[string]struct {
		json         string
		expected     any
		maxSize      int
		allowUnknown bool
		contentType  string
	}{
		"poor formatted json":           {json: `{"foo": }`, allowUnknown: false},
		"incorrect json type":           {json: `{"foo": 1}`, allowUnknown: false},
		"incorrect json type and value": {json: `{1: 1}`, allowUnknown: false},
		"several json passed":           {json: `{"foo": "2"}{"foo": "1"}`, allowUnknown: false},
		"empty json passed":             {json: ``, allowUnknown: false},
		"unknown json passed":           {json: `{"bar": "2"}`, allowUnknown: false},
		"json too large":                {json: `{"foo": "bar"}`, allowUnknown: false, maxSize: 1},
		"not json passed":               {json: `hello`, allowUnknown: false},
		"wrong header passwd":           {json: `{"foo": "bar"}`, allowUnknown: false, contentType: "application/text"},
	}

	for name, testCase := range tests {
		tool := Tool{
			MaxJSONSize:          testCase.maxSize,
			AllowedUnknownFields: testCase.allowUnknown,
		}

		var decodedJSON struct {
			Foo string `json:"foo"`
		}

		t.Run(name, func(t *testing.T) {
			request, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(testCase.json)))
			if testCase.contentType != "" {
				request.Header.Add("Content-Type", testCase.contentType)
			} else {
				request.Header.Add("Content-Type", "application/json")
			}

			rr := httptest.NewRecorder()
			err = tool.ReadJSONBody(rr, request, &decodedJSON)

			req.Error(err)
		})
	}
}

func TestTool_WriteJSONSuccessfully(t *testing.T) {
	req := require.New(t)

	var tests = map[string]struct {
		payload        any
		expected       any
		header         string
		expectedHeader string
	}{
		"valid":              {payload: JsonResponse{Error: false, Message: "bar"}, expected: "bar"},
		"valid with headers": {payload: JsonResponse{Error: false, Message: "bar"}, header: "FOO", expectedHeader: "BAR", expected: "bar"},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			var tool Tool
			rr := httptest.NewRecorder()
			headers := make(http.Header)
			headers.Add("FOO", "BAR")

			err := tool.WriteJSON(rr, http.StatusOK, testCase.payload, headers)
			req.NoError(err)
			req.Equal(rr.Header().Get(testCase.header), testCase.expectedHeader)
			req.Contains(rr.Body.String(), testCase.expected)
		})
	}
}

func TestTool_WriteJSONWithError(t *testing.T) {
	req := require.New(t)

	var tests = map[string]struct {
		payload        any
		expected       any
		header         string
		expectedHeader string
	}{
		"invalid":   {payload: make(chan int)},
		"invalid 2": {payload: make(chan string)},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			var tool Tool
			rr := httptest.NewRecorder()
			headers := make(http.Header)
			headers.Add("FOO", "BAR")

			err := tool.WriteJSON(rr, http.StatusOK, testCase.payload, headers)
			req.Error(err)
		})
	}
}

func TestTool_WriteErrorJSON(t *testing.T) {
	req := require.New(t)

	var tests = map[string]struct {
		status   int
		message  string
		expected any
	}{
		"bad request 400": {message: "new error", expected: "new error", status: http.StatusBadRequest},
		"bad request 401": {message: "new error passed", expected: "new error passed", status: http.StatusUnauthorized},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			var tool Tool
			rr := httptest.NewRecorder()

			err := tool.WriteErrorJSON(rr, errors.New(testCase.message), testCase.status)
			req.NoError(err)
			req.Equal(testCase.status, rr.Code)
			req.Contains(rr.Body.String(), testCase.expected)
		})
	}
}

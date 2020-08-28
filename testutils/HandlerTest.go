package testutils

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type HandlerTestCase struct {
	ErrorMessage string
	Status       int
}

func HandlerTest(t *testing.T, handler http.Handler, request *http.Request, url string, testCase HandlerTestCase, positiveResponseTest func(t *testing.T, body []byte)) {
	rec := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(url, handler)
	router.ServeHTTP(rec, request)

	if rec.Code != testCase.Status {
		t.Errorf("wrong status expected %d, got %d", testCase.Status, rec.Code)
	}

	body, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("can not read read recodeded Response Body; cause:%v", err)
	}

	if len(testCase.ErrorMessage) <= 0 {
		var respBody struct {
			Message   string    `json:"message"`
			TimeStamp time.Time `json:"timestamp"`
		}
		if err = json.Unmarshal(body, &respBody); err != nil {
			t.Fatalf("can not unmarshall body; cause: %v", err)
		}
		if respBody.TimeStamp.IsZero() {
			t.Errorf("expected a non zero timestamp")
		}
		if respBody.Message != testCase.ErrorMessage {
			t.Errorf("wrong status expected \"%s\", got \"%s\"", testCase.ErrorMessage, respBody.Message)
		}
		return
	}

	positiveResponseTest(t, body)
}


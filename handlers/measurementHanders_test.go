package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ChristophBe/weather-data-server/config"
	"github.com/ChristophBe/weather-data-server/data/models"
	"github.com/ChristophBe/weather-data-server/handlers"
	"github.com/ChristophBe/weather-data-server/services"
	"github.com/ChristophBe/weather-data-server/testutils"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)


type RequestTestCase struct {
	err         string
	status      int
	auth        string
	requestBody string
}
type AddMeasurementTestCase struct {
	RequestTestCase
	name     string
	nodeId   string
	Expected models.Measurement
}

func TestMeasurementHandlersImpl_AddMeasurementHandler(t *testing.T) {

	err := config.GetConfigManager().LoadConfig("../config_testing.json")
	if err != nil {
		t.Fatalf("failed to load confing; cause: %v", err)
	}
	authTokenService := services.GetAuthTokenService()

	node, _:= testutils.GetSavedMeasuringNode(t)
	validToken, err := authTokenService.GenerateNodeAccessToken(node)

	if err != nil {
		t.Fatalf("can not create accesstoke; cause: %v", err)
	}
	measurementHandlers := handlers.GetMeasurementHandlers()

	addMeasurementHandler := measurementHandlers.GetAddMeasurementHandler()

	testCases := []AddMeasurementTestCase{
		{
			RequestTestCase: RequestTestCase{
				err:    "invalid or insufficient authorization",
				status: http.StatusForbidden,
			},
			name:   "no authentication",
			nodeId: strconv.FormatInt(node.Id,10),
		},
		{
			RequestTestCase: RequestTestCase{
				err:    "invalid or insufficient authorization",
				status: http.StatusForbidden,
				auth:   "invalid",
			},
			name:   "invalid auth",
			nodeId: strconv.FormatInt(node.Id,10),
		},
		{
			RequestTestCase: RequestTestCase{
				err:    "invalid or insufficient authorization",
				status: http.StatusForbidden,
				auth:   validToken,
			},
			name:     "auth for other node",
			nodeId:   "42",
		},
		{
			RequestTestCase: RequestTestCase{
				err:         "invalid body",
				status:      http.StatusBadRequest,
				auth:        validToken,
				requestBody: `{pressure: 12.1,temperature: 12.2,humidity:42.1}`,
			},
			name:   "invalid body",
			nodeId: strconv.FormatInt(node.Id,10),
		},
		{
			RequestTestCase: RequestTestCase{
				err:    "invalid body",
				auth:   validToken,
				status: http.StatusBadRequest,
			},
			name:   "invalid path param",
			nodeId: "q",
		},
		{
			RequestTestCase: RequestTestCase{
				err:         "",
				auth:        validToken,
				requestBody: `{"pressure": 12.1,"temperature": 12.2,"humidity":42.1}`,
				status:      http.StatusAccepted,
			},
			name:   "valid node",
			nodeId: "1",
			Expected: models.Measurement{
				Pressure:    12.1,
				Temperature: 12.2,
				Humidity:    42.1,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testAddMeasurementHandlerRequest(t, testCase, addMeasurementHandler)
		})
	}
}

func testAddMeasurementHandlerRequest(t *testing.T, testCase AddMeasurementTestCase, addMeasurementHandler http.Handler) {
	bodyBytes := []byte(testCase.requestBody)
	reqBody := bytes.NewBuffer(bodyBytes)
	urlFormat := "/nodes/%v/measurements"
	url := fmt.Sprintf(urlFormat, testCase.nodeId)
	request, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		t.Fatalf("can not create request; cause: %v", err)
	}

	testRequest(t, addMeasurementHandler, request, fmt.Sprintf(urlFormat, "{nodeId}"), testCase.RequestTestCase, func(t *testing.T, body []byte) {
		var respBody models.Measurement
		if err = json.Unmarshal(body, &respBody); err != nil {
			t.Fatalf("can not unmarshall body; cause: %v", err)
		}
		if respBody.Id != 0 {
			t.Errorf("expected mesurement Id to be not zero")
		}

		if respBody.TimeStamp.IsZero() {
			t.Errorf("expected timestamp to be not zero")
		}

		if respBody.Humidity != testCase.Expected.Humidity {
			t.Errorf("wrong humidity expected %v, got %v", testCase.Expected.Humidity, respBody.Humidity)
		}

		if respBody.Temperature != testCase.Expected.Temperature {
			t.Errorf("wrong humidity expected %v, got %v", testCase.Expected.Temperature, respBody.Temperature)
		}

		if respBody.Pressure != testCase.Expected.Pressure {
			t.Errorf("wrong humidity expected %v, got %v", testCase.Expected.Pressure, respBody.Pressure)
		}
	})

}
func testRequest(t *testing.T, handler http.Handler, request *http.Request, url string, testCase RequestTestCase, positiveResponseTest func(t *testing.T, body []byte)) {

	if len(testCase.auth) > 0 {
		authValue := fmt.Sprintf("Bearer: %s", testCase.auth)
		request.Header.Add("Authorization", authValue)
	}

	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.Handle(url, handler)
	router.ServeHTTP(rec, request)

	if rec.Code != testCase.status {
		t.Errorf("wrong status expected %d, got %d", testCase.status, rec.Code)
	}

	body, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("can not read read recodeded Response Body; cause:%v", err)
	}

	if len(testCase.err) <= 0 {
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
		if respBody.Message != testCase.err {
			t.Errorf("wrong status expected \"%s\", got \"%s\"", testCase.err, respBody.Message)
		}
		return
	}

	positiveResponseTest(t, body)
}

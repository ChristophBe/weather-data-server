package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ChristophBe/weather-data-server/data/models"
	"github.com/ChristophBe/weather-data-server/data/repositories"
	"github.com/ChristophBe/weather-data-server/services"
	"github.com/ChristophBe/weather-data-server/testing/mockRepositories"
	"github.com/ChristophBe/weather-data-server/testing/mockServices"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)
func getMockAuthTokenService() services.AuthTokenService {
	return mockServices.MockAuthTokenService{}
}

func getMockMeasurementRepository() repositories.MeasuringRepository {
	return mockRepositories.MockMeasuringRepository{
		CreateMeasurementFunc: func(stationId int64, measurement models.Measurement) (savedMeasurement models.Measurement, err error) {
			savedMeasurement = measurement
			savedMeasurement.Id = 1
			return
		},
	}
}

func getMockNodeRepository() repositories.MeasuringNodeRepository {
	return mockRepositories.MockMeasuringNodeRepository{}
}

type AddMeasurementTestCase struct {
	name        string
	err         string
	nodeId      string
	requestBody string
	auth        string
	Expected    models.Measurement
	status      int
}

func TestMeasurementHandlersImpl_AddMeasurementHandler(t *testing.T) {
	validToken, err := getMockAuthTokenService().GenerateNodeAccessToken(models.MeasuringNode{Id: 1})
	if err != nil {
		t.Fatalf("can not create accesstoke; cause: %v", err)
	}
	measurementHandlers := measurementHandlersImpl{
		authTokenService:      getMockAuthTokenService(),
		measurementRepository: getMockMeasurementRepository(),
		nodeRepository:        getMockNodeRepository(),
	}
	addMeasurementHandler := measurementHandlers.GetAddMeasurementHandler()

	testCases := []AddMeasurementTestCase{
		{
			name:   "no authentication",
			err:    "invalid or insufficient authorization",
			auth:   "",
			nodeId: "1",
			status: http.StatusForbidden,
		},
		{
			name:   "invalid auth",
			err:    "invalid or insufficient authorization",
			auth:   "invalid",
			nodeId: "1",
			status: http.StatusForbidden,
		},
		{
			name:     "auth for other node",
			err:      "invalid or insufficient authorization",
			auth:     validToken,
			nodeId:   "42",
			Expected: models.Measurement{},
			status:   http.StatusForbidden,
		},
		{
			name:        "invalid body",
			err:         "invalid body",
			auth:        validToken,
			nodeId:      "1",
			requestBody: `{pressure: 12.1,temperature: 12.2,humidity:42.1}`,
			status:      http.StatusBadRequest,
		},
		{
			name:   "invalid path param",
			err:    "invalid body",
			auth:   validToken,
			nodeId: "q",
			status: http.StatusBadRequest,
		},
		{
			name:        "valid node",
			err:         "",
			auth:        validToken,
			nodeId:      "1",
			requestBody: `{"pressure": 12.1,"temperature": 12.2,"humidity":42.1}`,
			Expected: models.Measurement{
				Pressure:    12.1,
				Temperature: 12.2,
				Humidity:    42.1,
			},
			status: http.StatusAccepted,
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

	if len(testCase.auth) > 0 {
		authValue := fmt.Sprintf("Bearer: %s", testCase.auth)
		request.Header.Add("Authorization", authValue)
	}

	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.Handle(fmt.Sprintf(urlFormat, "{nodeId}"), addMeasurementHandler)
	router.ServeHTTP(rec, request)

	if rec.Code != testCase.status {
		t.Errorf("wrong status expected %d, got %d", testCase.status, rec.Code)
	}

	body, err := ioutil.ReadAll(rec.Body)
	if err != nil {

		return
	}

	if len(testCase.err) <= 0 {
		var respBody struct {
			Message string `json:"message"`
		}
		if err = json.Unmarshal(body, &respBody); err != nil {
			t.Fatalf("can not unmarshall body; cause: %v", err)
		}

		if respBody.Message != testCase.err {
			t.Errorf("wrong status expected \"%s\", got \"%s\"", testCase.err, respBody.Message)
		}
		return
	}

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
}

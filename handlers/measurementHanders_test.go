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
	"net/http"
	"strconv"
	"testing"
)



type AddMeasurementTestCase struct {
	testutils.HandlerTestCase
	Name        string
	RequestBody string
	AuthToken   string
	NodeId      string
	Expected    models.Measurement
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
			HandlerTestCase: testutils.HandlerTestCase{
				ErrorMessage: "invalid or insufficient authorization",
				Status:       http.StatusForbidden,
			},
			Name:   "no authentication",
			NodeId: strconv.FormatInt(node.Id,10),
		},
		{
			HandlerTestCase: testutils.HandlerTestCase{
				ErrorMessage:    "invalid or insufficient authorization",
				Status: http.StatusForbidden,
			},
			AuthToken: "invalid",
			Name:      "invalid auth",
			NodeId:    strconv.FormatInt(node.Id,10),
		},
		{
			HandlerTestCase: testutils.HandlerTestCase{
				ErrorMessage:    "invalid or insufficient authorization",
				Status: http.StatusForbidden,
			},
			AuthToken: validToken,
			Name:      "auth for other node",
			NodeId:    "42",
		},
		{
			HandlerTestCase: testutils.HandlerTestCase{
				ErrorMessage:         "invalid body",
				Status:      http.StatusBadRequest,
			},
			RequestBody: `{pressure: 12.1,temperature: 12.2,humidity:42.1}`,
			AuthToken:   validToken,
			Name:        "invalid body",
			NodeId:      strconv.FormatInt(node.Id,10),
		},
		{
			HandlerTestCase: testutils.HandlerTestCase{
				ErrorMessage:    "invalid body",
				Status: http.StatusBadRequest,
			},
			AuthToken: validToken,
			Name:      "invalid path param",
			NodeId:    "q",
		},
		{
			HandlerTestCase: testutils.HandlerTestCase{
				Status:      http.StatusAccepted,
			},
			AuthToken:   validToken,
			RequestBody: `{"pressure": 12.1,"temperature": 12.2,"humidity":42.1}`,
			Name:        "valid node",
			NodeId:      strconv.FormatInt(node.Id,10),
			Expected: models.Measurement{
				Pressure:    12.1,
				Temperature: 12.2,
				Humidity:    42.1,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			urlFormat := "/nodes/%v/measurements"
			url := fmt.Sprintf(urlFormat, testCase.NodeId)

			bodyBytes := []byte(testCase.RequestBody)
			request, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
			
			if err != nil {
				t.Fatalf("can not create request; cause: %v", err)
			}

			if len(testCase.AuthToken) > 0 {
				authValue := fmt.Sprintf("Bearer: %s", testCase.AuthToken)
				request.Header.Add("Authorization", authValue)
			}

			testutils.HandlerTest(t, addMeasurementHandler, request, fmt.Sprintf(urlFormat, "{nodeId}"), testCase.HandlerTestCase, func(t *testing.T, body []byte) {
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
		})
	}
}

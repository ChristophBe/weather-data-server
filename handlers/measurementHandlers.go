package handlers

import (
	"github.com/ChristophBe/weather-data-server/data/database"
	"github.com/ChristophBe/weather-data-server/services"
	"net/http"
)

type MeasurementHandlers interface {
	GetAddMeasurementHandler() http.Handler
	GetMeasurementsByNodeHandler() http.Handler
}

func GetMeasurementHandlers() MeasurementHandlers {
	return measurementHandlersImpl{
		authTokenService:      services.GetAuthTokenService(),
		measurementRepository: database.GetMeasurementRepository(),
		nodeRepository:        database.GetMeasuringNodeRepository(),
	}
}

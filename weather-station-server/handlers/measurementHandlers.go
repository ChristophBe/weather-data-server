package handlers

import (
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/services"
	"net/http"
)

type MeasurementHandlers interface{
	GetAddMeasurementHandler() http.Handler
	GetMeasurementsByNodeHandler() http.Handler
}

func GetMeasurementHandlers()  MeasurementHandlers{
	return measurementHandlersImpl{
		authTokenService:      services.GetAuthTokenService(),
		measurementRepository: database.GetMeasurementRepository(),
		nodeRepository: 	   database.GetMeasuringNodeRepository(),
	}
}
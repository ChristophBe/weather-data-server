package handlers

import (
	"errors"
	"fmt"
	"github.com/ChristophBe/weather-data-server/data/models"
	"github.com/ChristophBe/weather-data-server/data/repositories"
	"github.com/ChristophBe/weather-data-server/handlers/httpHandler"
	"github.com/ChristophBe/weather-data-server/services"
	"net/http"
	"strconv"
	"time"
)

type measurementHandlersImpl struct {
	authTokenService      services.AuthTokenService
	measurementRepository repositories.MeasuringRepository
	nodeRepository        repositories.MeasuringNodeRepository
}

func (m measurementHandlersImpl) GetAddMeasurementHandler() http.Handler {
	return httpHandler.AuthorizedAppHandler(m.authTokenService.VerifyNodeAccessToken, m.addMeasurementHandler)
}

func (m measurementHandlersImpl) GetMeasurementsByNodeHandler() http.Handler {
	return httpHandler.AuthorizedAppHandlerWithUnauthorisedFallback(m.authTokenService.VerifyUserAccessToken, m.fetchMeasurementsByNodeAuthorized, m.fetchMeasurementsByNodePublic)
}

func (m measurementHandlersImpl) addMeasurementHandler(nodeId int64, r *http.Request) (response httpHandler.HandlerResponse, err error) {

	nodeIdPath, err := httpHandler.ReadPathVariableInt(r, "nodeId")
	if err != nil {
		message := fmt.Sprintf(httpHandler.ErrorMessageParameterf, "nodeId")
		err = httpHandler.BadRequest(message, err)
		return
	}

	if nodeId != nodeIdPath {
		err = httpHandler.Forbidden(httpHandler.ErrorMessageNotAuthorized, errors.New("authorization is not valid for this node"))
		return
	}

	var measuring models.Measurement
	err = httpHandler.ReadJsonBody(r, measuring)
	if err != nil {
		err = httpHandler.BadRequest(httpHandler.ErrorMessageInvalidBody, err)
		return
	}

	//TODO: Check validity of measured data.

	measuring.TimeStamp = time.Now()

	measuring, err = m.measurementRepository.CreateMeasurement(nodeId, measuring)
	if err != nil {
		err = httpHandler.InternalError(err)
		return
	}

	response.Data = measuring
	response.Status = http.StatusAccepted
	return
}

func (m measurementHandlersImpl) fetchMeasurementsByNodeAuthorized(userId int64, r *http.Request) (httpHandler.HandlerResponse, error) {
	return m.fetchMeasurementsByNode(r, func(node models.MeasuringNode) bool {
		if node.IsPublic {
			return true
		}
		relations, err := m.nodeRepository.FetchAllMeasuringNodeUserRelations(node.Id, userId)
		return err == nil && len(relations) > 0
	})
}

func (m measurementHandlersImpl) fetchMeasurementsByNodePublic(r *http.Request) (httpHandler.HandlerResponse, error) {
	return m.fetchMeasurementsByNode(r, func(node models.MeasuringNode) bool {
		return node.IsPublic
	})
}

func (m measurementHandlersImpl) fetchMeasurementsByNode(r *http.Request, checkPermission func(node models.MeasuringNode) bool) (response httpHandler.HandlerResponse, err error) {

	nodeId, err := httpHandler.ReadPathVariableInt(r, "nodeId")
	if err != nil {
		message := fmt.Sprintf(httpHandler.ErrorMessageParameterf, "nodeId")
		err = httpHandler.BadRequest(message, err)
		return
	}

	node, err := m.nodeRepository.FetchMeasuringNodeById(nodeId)
	if err != nil || node.Id != nodeId {
		message := fmt.Sprintf(httpHandler.ErrorMessageNotFoundf, "node")
		err = httpHandler.NotFound(message, err)
	}

	if !checkPermission(node) {
		err = httpHandler.Forbidden(httpHandler.ErrorMessageNotAuthorized, errors.New("user is not authorized for this node"))
		return
	}

	var measurements = make([]models.Measurement, 0)

	limitValue := r.FormValue("limit")
	if len(limitValue) > 0 {
		var limit int64
		limit, err = strconv.ParseInt(limitValue, 10, 64)

		if err != nil {
			message := fmt.Sprintf(httpHandler.ErrorMessageParameterf, "nodeId")
			err = httpHandler.BadRequest(message, err)

			return
		}
		measurements, err = m.measurementRepository.FetchLastMeasuringsByNodeId(nodeId, limit)
	} else {
		measurements, err = m.measurementRepository.FetchAllMeasuringsByNodeId(nodeId)
	}

	if err != nil {
		err = httpHandler.InternalError(err)
		return
	}

	response.Data = measurements
	response.Status = http.StatusOK
	return
}

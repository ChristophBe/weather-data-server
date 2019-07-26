import {
    MEASUREMENTS_FETCH, MEASUREMENTS_RECEIVE, MEASUREMENTS_REQUEST_FAILED,
} from "./ActionTypes";

export default class MeasurementsActions {

    static fechMeasurements= (nodeId) => ({
        type: MEASUREMENTS_FETCH,
        nodeId
    });

    static receiveMeasurements = (nodeId, measurements) => ({
        type: MEASUREMENTS_RECEIVE,
        measurements,
        nodeId
    });

    static measurementsRequestFailed = (nodeId) => ({
        type: MEASUREMENTS_REQUEST_FAILED,
        nodeId
    });


}

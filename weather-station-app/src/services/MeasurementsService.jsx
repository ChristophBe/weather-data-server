import RequestService from "./RequestService";
import {Config} from "../config";


export class MeasurementsService {

    static fetchMeasurementsByNodeId(nodeId,authToken){
        return RequestService.doGet(Config.apiBaseUrl + "/nodes/" + nodeId + "/measurements?limit=24" , {},authToken);
    }


}
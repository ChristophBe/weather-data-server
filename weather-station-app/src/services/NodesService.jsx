import RequestService from "./RequestService";
import {Config} from "../config";


export class NodesService {

    static fetchNodes(authToken){
        return RequestService.doGet(Config.apiBaseUrl + "/nodes" , {},authToken);
    }


}
import RequestService from "./RequestService";
import {Config} from "../config";
import {NODES_RECEIVE} from "../redux/actions/ActionTypes";


export class NodesService {

    static fetchNodes(authToken){
        return RequestService.doGet(Config.apiBaseUrl + "/nodes" , {},authToken);
    }

    static fetchNodesByOwnersUserId(userId, authToken){
        return RequestService.doGet(Config.apiBaseUrl + "/users/" +  userId + "/nodes" , {},authToken);
    }

    static shareNode(nodeId, email, authToken){
        const payLoad = {email};
        return RequestService.doPost(Config.apiBaseUrl + "/nodes/" +  nodeId + "/share" , payLoad,true,{},authToken);
    }


}
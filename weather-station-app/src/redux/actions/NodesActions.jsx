import {
    NODES_FETCH,NODES_RECEIVE,NODES_REQUEST_FAILED,NODES_SELECT
} from "./ActionTypes";

export default class NodesActions {

    static fechNodes = () => ({
        type: NODES_FETCH,
    });

    static receiveNodes = (nodes) => ({
        type: NODES_RECEIVE,
        nodes
    });

    static nodesRequestFailed = () => ({
        type: NODES_REQUEST_FAILED,
    });

    static selectNode = (nodeId) =>({
        type: NODES_SELECT,
        nodeId
    });

}

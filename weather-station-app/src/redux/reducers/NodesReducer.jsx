import {
    NODES_FETCH, NODES_OWN_NODES_RECEIVE, NODES_RECEIVE, NODES_REQUEST_FAILED, NODES_SELECT
} from "../actions/ActionTypes";


function mapNodeListToMap(inputList) {

    const reducer = (accumulator, currentValue) => {
        if(!accumulator.hasOwnProperty(currentValue.id)){
            accumulator[currentValue.id] = currentValue;
        }
        return accumulator;
    };

    return inputList.reduce(reducer,{})
}

const INITIAL_AUTHENTICATION_STATE = {
    isFailed: false,
    isLoading: false,
    map: {},
    ownNodes:[],
    selectedNode: null,
};

const NodesReducer = (state = INITIAL_AUTHENTICATION_STATE, action) => {
    switch(action.type) {
        case NODES_FETCH:
            return Object.assign({},state,{
                isLoading:true,
                isFailed: false,
            });
        case NODES_RECEIVE:
            return Object.assign({},state,{
                map: mapNodeListToMap(action.nodes),
                isFailed: false,
                isLoading: false,
                selectedNode: state.selectedNode == null ? action.nodes[0].id : state.selectedNode
            });
        case NODES_REQUEST_FAILED:
            return Object.assign({},state,{

                isFailed: true,
                isLoading: false
            });

        case NODES_SELECT:
            return Object.assign({},state,{
                selectedNode: action.nodeId
            });
        case NODES_OWN_NODES_RECEIVE:
            return Object.assign({},state,{
                ownNodes: action.nodes
            });
        default:
            return state
    }
};
export default NodesReducer;
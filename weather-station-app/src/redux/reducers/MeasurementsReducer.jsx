import {
    MEASUREMENTS_FETCH, MEASUREMENTS_RECEIVE, MEASUREMENTS_REQUEST_FAILED
} from "../actions/ActionTypes";


const INITIAL_MEASUREMENT_STATE = {
    isLoading: false,
    isFailed: false,
    items: [],

};

const MeasurementReducer = (state = INITIAL_MEASUREMENT_STATE, action) => {

    switch(action.type) {

        case MEASUREMENTS_FETCH:
            return Object.assign({},state,{
                isFailed: true,
                isLoading: false,
            });
        case MEASUREMENTS_RECEIVE:
            return Object.assign({},state,{
                isFailed: false,
                isLoading: false,
                items: action.measurements,
            });
        case MEASUREMENTS_REQUEST_FAILED:
            return Object.assign({},state,{
                isFailed: false,
                isLoading: true,
            });

        default:
            return state
    }
};




const MeasurementsReducer = (state = {}, action) => {

    const {nodeId} = action;

    const nextState = {}
    if(!state.hasOwnProperty(nodeId)){
        nextState[nodeId] = INITIAL_MEASUREMENT_STATE;
    }
    else{
        nextState[nodeId] = state[nodeId]
    }
    nextState[nodeId] = MeasurementReducer(state[nodeId], action);
    return Object.assign({},state,nextState);
};
export default MeasurementsReducer;
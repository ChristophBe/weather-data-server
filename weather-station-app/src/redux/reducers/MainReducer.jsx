import {combineReducers} from "redux";
import AuthReducer from "./AuthReducer";
import NodesReducer from "./NodesReducer";
import MeasurementsReducer from "./MeasurementsReducer";

export const MainReducer= combineReducers({
    auth: AuthReducer,
    nodes: NodesReducer,
    measurements: MeasurementsReducer,
});
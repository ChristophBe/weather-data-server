import {call, put, select} from "redux-saga/effects";
import NodesActions from "../actions/NodesActions";
import {NodesService} from "../../services/NodesService";


const getAuth = (state) => state.auth;




export function* fetchNodes() {
    try {
        const {token} = yield select(getAuth);
        const nodes = yield call(NodesService.fetchNodes, token);
        yield put(NodesActions.receiveNodes(nodes))
    }
    catch (e) {
        yield put(NodesActions.nodesRequestFailed())
    }
}

export function* fetchOwnNodes() {
    try {
        const {token,currentUser} = yield select(getAuth);
        const nodes = yield call(NodesService.fetchNodesByOwnersUserId,currentUser.id, token);
        yield put(NodesActions.receiveOwnNodes(nodes))
    }
    catch (e) {
        yield put(NodesActions.nodesRequestFailed())
    }
}
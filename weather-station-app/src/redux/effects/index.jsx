import {fetchAuthTokenEffect, fetchUserEffect, logoutAuthCookieEffect, tryCookieAuthEffect} from "./AuthEffects";

import {takeEvery, takeLatest} from 'redux-saga/effects'
import {
    AUTH_DO_LOGIN, AUTH_DO_LOGOUT,
    AUTH_LOGIN_SUCCESS, AUTH_RECEIVE_USER, AUTH_TRY_COOKIE_AUTH, MEASUREMENTS_FETCH, NODES_FETCH,
} from "../actions/ActionTypes";
import {fetchNodes, fetchOwnNodes} from "./NodesEffects";
import {fetchMeasurements} from "./MeasurementsEffects";

export function* rootSaga(){
    yield takeLatest(AUTH_DO_LOGIN, fetchAuthTokenEffect);
    yield takeEvery(AUTH_LOGIN_SUCCESS, fetchUserEffect);
    yield takeEvery(AUTH_TRY_COOKIE_AUTH, tryCookieAuthEffect);
    yield takeEvery(AUTH_DO_LOGOUT, logoutAuthCookieEffect);


    yield takeEvery(AUTH_DO_LOGOUT, fetchNodes);
    yield takeEvery(AUTH_LOGIN_SUCCESS, fetchNodes);

    yield takeEvery(NODES_FETCH, fetchNodes);
    yield takeEvery(AUTH_RECEIVE_USER, fetchOwnNodes);

    yield takeEvery(MEASUREMENTS_FETCH, fetchMeasurements);


}
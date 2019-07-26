import {AuthenticationService} from "../../services/AuthenticationService";
import {call, put} from "redux-saga/effects";
import AuthActions from "../actions/AuthActions";
import {UserService} from "../../services/UserService";
import CookieAuthenticationService from "../../services/CookieAuthenticationService";
import NodesActions from "../actions/NodesActions";

export function* fetchAuthTokenEffect(action) {

    try {
        const token = yield call(AuthenticationService.authenticateEmailPassword, action.authDetails);
        yield put(AuthActions.loginSuccessLogin(token));

    } catch (e) {
        yield put(AuthActions.loginFailedLogin(e));
    }
}

export function* fetchUserEffect(action) {
    try {
        const user = yield call(UserService.fetchUserMe, action.token);
        yield put(AuthActions.receiveUser(user))
    }
    catch (e) {
        yield put(AuthActions.doLogout())
    }
}

export function* tryCookieAuthEffect() {
    try {
        console.log("try Cookie Auth!");
        const token = yield call(CookieAuthenticationService.restoreToken);
        yield put(AuthActions.loginSuccessLogin(token));
    } catch (e) {
        yield put(AuthActions.cookieAuthFailed())
        console.log(e)

    }
}

export function* logoutAuthCookieEffect(){
    yield call(CookieAuthenticationService.logout);
}

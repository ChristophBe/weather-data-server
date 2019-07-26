import {
    AUTH_COOKIE_AUTH_FAILED,
    AUTH_DO_LOGIN,
    AUTH_DO_LOGOUT,
    AUTH_LOGIN_FAILED,
    AUTH_LOGIN_SUCCESS,
    AUTH_RECEIVE_USER, AUTH_TRY_COOKIE_AUTH,
} from "./ActionTypes";

export default class AuthActions {

    static tryCookieAuth = () => ({
        type: AUTH_TRY_COOKIE_AUTH,
    });

    static doLogin = (authDetails) => ({
        type: AUTH_DO_LOGIN,
        authDetails
    });

    static loginSuccessLogin = (token) => ({
        type: AUTH_LOGIN_SUCCESS,
        token
    });

    static loginFailedLogin = (error) =>({
        type: AUTH_LOGIN_FAILED,
        error
    });

    static doLogout = ()=> ({
        type: AUTH_DO_LOGOUT
    });

    static receiveUser  = (user) => ({
        type: AUTH_RECEIVE_USER,
        user
    });

    static cookieAuthFailed = () => ({
        type: AUTH_COOKIE_AUTH_FAILED
    })
}

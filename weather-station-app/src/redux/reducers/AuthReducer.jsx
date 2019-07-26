import {
    AUTH_DO_LOGIN,
    AUTH_DO_LOGOUT,
    AUTH_LOGIN_FAILED,
    AUTH_LOGIN_SUCCESS, AUTH_RECEIVE_USER,
    AUTH_TRY_COOKIE_AUTH,AUTH_COOKIE_AUTH_FAILED
} from "../actions/ActionTypes";


const INITIAL_AUTHENTICATION_STATE = {
    isAuthenticationFailed: false,
    isLoading: false,
    token: null,
    currentUser: null
};

const AuthReducer = (state = INITIAL_AUTHENTICATION_STATE, action) => {
    switch(action.type) {
        case AUTH_DO_LOGIN:
            return Object.assign({},state,{
                isLoading:true
            });


        case AUTH_TRY_COOKIE_AUTH:
            return Object.assign({},state,{
                isLoading:true
            });

        case AUTH_LOGIN_SUCCESS:
            return Object.assign({},state,{

                isAuthenticationFailed:false,
                isLoading:false,
                token:action.token
            });
        case AUTH_LOGIN_FAILED:
            return Object.assign({},state,{
                isAuthenticationFailed:true,
                isLoading:false,
            });

        case AUTH_COOKIE_AUTH_FAILED:
            return Object.assign({},state,{
                isLoading:false,
            });


        case AUTH_RECEIVE_USER:
            return Object.assign({},state,{
                currentUser: action.user
            });

        case AUTH_DO_LOGOUT:
            return INITIAL_AUTHENTICATION_STATE;
        default:
            return state;
    }
};
export default AuthReducer;
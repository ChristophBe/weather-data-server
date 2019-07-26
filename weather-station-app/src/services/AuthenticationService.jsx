import {Config} from "../config";
import RequestService from "./RequestService";
import CookieAuthenticationService from "./CookieAuthenticationService";


export class AuthenticationService {

    static authenticateEmailPassword(authDetails){
        const data = {
            email: authDetails.email,
            password: authDetails.password,
        };

        const req =AuthenticationService.fetchAuthToken(data);
        req.then(token=>CookieAuthenticationService.saveToken(token));
        return req
    }

    static authenticateRefreshToken(token){
        const data = {
            refresh_token: token.refresh_token
        };
        return AuthenticationService.fetchAuthToken(data)
    }

    static generateAuthHeader(token){
        return {"Authorization": token.token_type + " "+ token.access_token }
    }

    static fetchAuthToken(data){
        const url = Config.apiBaseUrl + "/users/login";

        return RequestService.doPost(url,data,true,{});
    }
}
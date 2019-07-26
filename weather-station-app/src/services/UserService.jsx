import RequestService from "./RequestService";
import {Config} from "../config";


export class UserService {

    static fetchUserMe(authToken){
        return RequestService.doGet(Config.apiBaseUrl + "/users/me" , {},authToken);
    }


}
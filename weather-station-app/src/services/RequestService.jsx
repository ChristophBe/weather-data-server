import * as $ from "jquery";
import {AuthenticationService} from "./AuthenticationService";


export default class RequestService{

    static doPost(url, data, json, headers, authToken){
        headers = RequestService.assignAuthToken(headers,authToken);
        let extras = {};

        if(json){
            data = JSON.stringify(data);
            extras = {
                contentType: "application/json; charset=utf-8",
                dataType: "json"
            }
        }

        const request = $.post({
            data,
            headers,
            url,
            ...extras
        });

        return this.handleRequest(request);
    }
    static doGet(url, headers, authToken){


        headers = RequestService.assignAuthToken(headers,authToken);
        const request = $.getJSON({
            headers,
            url
        });

        return RequestService.handleRequest(request);
    }
    static doDelete(url, headers, authToken){
        headers = RequestService.assignAuthToken(headers,authToken);
        const request = $.ajax({
            type:"delete",
            headers,
            url
        });

        return RequestService.handleRequest(request);
    }

    static handleRequest(request) {
        return new Promise((resolve, reject) => {
            request.done((resp)=>{
                window.console.log(resp);

                resolve(resp)
            });
            request.fail(reject)
        });

    }

    static assignAuthToken(headers, authToken){
        if(authToken){
            const authHeaders = AuthenticationService.generateAuthHeader(authToken);
            return Object.assign({}, headers, authHeaders);
        }
        return headers;
    }


}

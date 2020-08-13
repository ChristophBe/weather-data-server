package httpHandler

import (
	"errors"
	"net/http"
	"strings"
)

func AuthorizedAppHandler(
	tokenVerifier func(token string) (sub int64, err error),
	handler func(sub int64, r *http.Request) (HandlerResponse, error),
) JsonHandler {
	return AuthorizedAppHandlerWithUnauthorisedFallback(tokenVerifier, handler, func(r *http.Request) (HandlerResponse, error) {
		return HandlerResponse{}, Forbidden(ErrorMessageNotAuthorized, nil)
	})
}
func AuthorizedAppHandlerWithUnauthorisedFallback(
	tokenVerifier func(token string) (sub int64, err error),
	handlerAuthorized func(sub int64, r *http.Request) (HandlerResponse, error),
	handlerUnauthorized JsonHandler,
) JsonHandler {
	return func(r *http.Request) (response HandlerResponse, err error) {
		_, tokenString, err := readTokenFormRequest(r)

		if err != nil {
			return handlerUnauthorized(r)
		}

		sub, err := tokenVerifier(tokenString)

		if err != nil {
			err =  Forbidden(ErrorMessageNotAuthorized, err)
			return
		}

		return handlerAuthorized(sub, r)
	}
}

func readTokenFormRequest(r *http.Request) (tokenType string, tokenString string, err error) {
	auth := r.Header.Get("Authorization")

	authHeaderParts := strings.Split(auth, " ")
	if len(authHeaderParts) < 2 {
		err = errors.New("token not found")
		return
	}
	tokenType = authHeaderParts[0]
	tokenString = authHeaderParts[1]
	return
}

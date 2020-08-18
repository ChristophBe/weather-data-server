package transitory

import "regexp"

func isValidEmail(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return len(email) < 254 && rxEmail.MatchString(email)
}

type UserCreateBody struct {
	InvitationToken string `json:"invitation_token"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
}

func (u UserCreateBody) IsValid() bool {
	return len(u.Username) > 4 && len(u.Password) > 4 && isValidEmail(u.Email)
}

package controllers

import (
	"strings"

	"github.com/kataras/iris/v12/sessions"
)

const userIDKey = "UserID"

func isLoggedIn(session *sessions.Session) (string, bool) {
	userId := session.GetStringDefault(userIDKey, "")
	return userId, userId != ""
}

func getUnitHost(v, host string) string {
	if strings.HasPrefix(v, ":") {
		baseHost := host
		if idx := strings.Index(baseHost, ":"); idx > -1 {
			baseHost = baseHost[:idx]
		}
		return "http://" + baseHost + v
	}

	if strings.HasPrefix(v, "/") {
		return "http://" + host + v
	}

	return v
}

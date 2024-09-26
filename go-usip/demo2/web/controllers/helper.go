package controllers

import "github.com/kataras/iris/v12/sessions"

const userIDKey = "UserID"

func isLoggedIn(session *sessions.Session) (string, bool) {
	userId := session.GetStringDefault(userIDKey, "")
	return userId, userId != ""
}

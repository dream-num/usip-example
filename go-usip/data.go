package main

var (
	Users = map[string]*User{
		"1": {UserID: "1", Name: "user1", Avatar: "http://avatar.com/1"},
		"2": {UserID: "2", Name: "user2", Avatar: "http://avatar.com/2"},
		"3": {UserID: "3", Name: "user3", Avatar: "http://avatar.com/3"},
	}

	UnitCollaborators = map[string][]*Collaborator{
		"unit1": {
			{UserID: "1", Role: RoleOwner},
			{UserID: "2", Role: RoleEditor},
		},
		"unit2": {
			{UserID: "2", Role: RoleOwner},
			{UserID: "3", Role: RoleReader},
		},
		"unit3": {
			{UserID: "3", Role: RoleOwner},
			{UserID: "1", Role: RoleEditor},
		},
	}
)

func VerifyToken(token string) (userID string, ok bool) {
	switch token {
	case "token:1":
		return "1", true
	case "token:2":
		return "2", true
	case "token:3":
		return "3", true
	default:
		return "", false
	}
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/userinfo", BatchGetUserInfo)
	http.HandleFunc("/collaborators", BatchGetCollaborators)
	http.HandleFunc("/role", GetUnitCollaboratorRole)
	http.HandleFunc("/credential", CredentialVerify)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

const (
	RoleOwner  = "owner"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

type User struct {
	UserID string `json:"userID,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type Collaborator struct {
	UserID string `json:"userID,omitempty"`
	Role   string `json:"role,omitempty"`
}

// handles the batch get user info request
func BatchGetUserInfo(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Users []*User `json:"users"`
	}

	_ = r.ParseForm()
	userIDs := r.Form["userIDs"]
	result := response{}
	for _, userID := range userIDs {
		if u, ok := Users[userID]; ok {
			result.Users = append(result.Users, u)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// handles the batch get collaborators request
func BatchGetCollaborators(w http.ResponseWriter, r *http.Request) {
	type subject struct {
		ID     string `json:"id,omitempty"`
		Name   string `json:"name,omitempty"`
		Avatar string `json:"avatar,omitempty"`
		// Type   string `json:"type,omitempty"`
	}
	type response struct {
		Collaborators []*struct {
			UnitID   string `json:"unitID,omitempty"`
			Subjects []*struct {
				Subject *subject `json:"subject,omitempty"`
				Role    string   `json:"role,omitempty"`
			} `json:"subjects,omitempty"`
		} `json:"collaborators"`
	}

	_ = r.ParseForm()
	unitIDs := r.Form["unitIDs"]
	result := response{}
	for _, unitID := range unitIDs {
		if cs, ok := UnitCollaborators[unitID]; ok {
			item := struct {
				UnitID   string `json:"unitID,omitempty"`
				Subjects []*struct {
					Subject *subject `json:"subject,omitempty"`
					Role    string   `json:"role,omitempty"`
				} `json:"subjects,omitempty"`
			}{UnitID: unitID}
			for _, c := range cs {
				s := &subject{}
				if u, ok := Users[c.UserID]; ok {
					s.ID = u.UserID
					s.Name = u.Name
					s.Avatar = u.Avatar
				}
				item.Subjects = append(item.Subjects, &struct {
					Subject *subject `json:"subject,omitempty"`
					Role    string   `json:"role,omitempty"`
				}{Subject: s, Role: c.Role})
			}
			result.Collaborators = append(result.Collaborators, &item)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// handles the get unit collaborator role request
func GetUnitCollaboratorRole(w http.ResponseWriter, r *http.Request) {
	type response struct {
		UserID string `json:"userID,omitempty"`
		Role   string `json:"role,omitempty"`
	}

	_ = r.ParseForm()
	userID := r.FormValue("userID")
	unitID := r.FormValue("unitID")

	role := ""
	if cs, ok := UnitCollaborators[unitID]; ok {
		for _, c := range cs {
			if c.UserID == userID {
				role = c.Role
				break
			}
		}
	}

	result := response{Role: role, UserID: userID}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// handles the credential verify request
func CredentialVerify(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User *User `json:"user,omitempty"`
	}

	userID, ok := VerifyToken(r.Header.Get("x-authorization"))
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result := response{User: Users[userID]}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

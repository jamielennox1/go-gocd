package gocd

import (
	"reflect"
	"strings"
)

type User struct {
	LoginName      string   `json:"login_name"`
	DisplayName    string   `json:"display_name"`
	Enabled        bool     `json:"enabled"`
	Email          string   `json:"email"`
	EmailMe        bool     `json:"email_me"`
	CheckinAliases []string `json:"checkin_aliases"`
}

func NewUser() *User {
	return &User{CheckinAliases: make([]string, 0)}
}

func (p User) Diff(user *User) map[string]interface{} {
	result := make(map[string]interface{})
	if strings.Compare(p.LoginName, user.LoginName) != 0 {
		result["login_name"] = user.LoginName
	}
	if strings.Compare(p.DisplayName, user.DisplayName) != 0 {
		result["display_name"] = user.DisplayName
	}
	if p.Enabled == user.Enabled {
		result["enabled"] = user.Enabled
	}
	if strings.Compare(p.Email, user.Email) != 0 {
		result["email"] = user.Email
	}
	if p.EmailMe == user.EmailMe {
		result["email_me"] = user.EmailMe
	}
	if reflect.DeepEqual(p.CheckinAliases, user.CheckinAliases) {
		result["checkin_aliases"] = user.CheckinAliases
	}
	return result
}

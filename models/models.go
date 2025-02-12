package models

import "github.com/dgrijalva/jwt-go"

type Gun struct {
	Gun_Id int `json:"gun_id", db:"gun_id"`
	Type string `json:"type", db:"type"`
	Name string `json:"name", db:"name"`
}

type Attachment struct {
	Attachment_Id int `json:"attachment_id", db:"attachment_id"`
	Gun_Id int `json:"gun_id", db:gun_id`
	Name string `json:"name", db:"name"`
	SubAttachments []string `json:"subattachments", db:subattachments`
}

type Credentials struct {
	User_Id int `json:"user_id", db:"user_id"`
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
	Token string `json:"token", db:"token"`
}

type Loadout struct {
	Loadout_Id int `json:"loadout_id", db:"loadout_id"`
	User_Id int `json:"user_id", db:"user_id"`
	Type string `json:"type", db:"type"`
	Gun string `json:"gun", db:"gun"`
	Attachments []string `json:"attachments", db:"attachments"`
	SubAttachments []string `json:"subattachments", db:"subattachments"`
	Description string `json:"description", db:"description"`
}

type User struct {
	Username string `json:"username", db:"username"`
	User_Id int `json:"user_id", db:"user_id"`
	Logged_In bool `json:"logged_in"`
}

type Claims struct {
	Username string `json:"username"`
	User_Id int `json:"user_id"`
	jwt.StandardClaims
}

type Error struct {
	Status string `json:"status"`
	Text string `json:"text"`
	Message string `json:"message"`
}

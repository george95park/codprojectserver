package models

type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
}
 type Loadout struct {
	 Id string `json:"id", db:"id"`
	 Username string `json:"username", db:"username"`
	 Gunname string `json:"gunname", db:"gunname"`
	 Attachments []string `json:"attachments", db:"attachments"`
	 Description string `json:"description", db:"description"`
 }

type User struct {
	Username string `json:"username", db:"username"`
}

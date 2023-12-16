package models

type User struct {
	ID       int    `json:"id", gorm:"primary_key;auto_increment"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

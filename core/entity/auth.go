package entity

type User struct {
	Email    string `json:"email"`
	FullName string `json:"fullname" gorm:"column:fullname"`
	Password string `json:"password"`
}

type Auth struct {
	Email       string `json:"email"`
	FullName    string `json:"fullname"`
	AccessToken string `json:"accessToken"`
}

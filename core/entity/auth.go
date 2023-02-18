package entity

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullname" gorm:"column:fullname"`
	Password string `json:"password"`
}

type Auth struct {
	Email       string `json:"email"`
	FullName    string `json:"fullname"`
	AccessToken string `json:"accessToken"`
}

type CredentialClaim struct {
	ID       int
	Email    string
	FullName string
}

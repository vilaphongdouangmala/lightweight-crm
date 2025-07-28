package domain

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email" gorm:"unique"`
	Password     string `json:"password"`
	BuisnessName string `json:"buisness_name"`
}

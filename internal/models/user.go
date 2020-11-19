package models

type User struct {
	Nickname string `json:"nickname"` // Имя пользователя (уникальное поле). Данное поле допускает только латиницу, цифры и знак подчеркивания.
	FullName string `json:"fullname"`
	About    string `json:"about"` // Описание пользователя.
	Email    string `json:"email"` // Почтовый адрес пользователя (уникальное поле).
	
}

type UsersUpdate struct {
	Nickname string `json:"-"`
	FullName string `json:"fullname"`
	About    string `json:"about"` // Описание пользователя.
	Email    string `json:"email"` // Почтовый адрес пользователя (уникальное поле).
}

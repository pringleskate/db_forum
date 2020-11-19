package models

type Post struct {
	SlagOrID string `json:"-"`
	Id       int64  `json:"id"`       // Идентификатор данного сообщения.
	Parent   int64  `json:"parent"`   // Идентификатор родительского сообщения (0 - корневое сообщение обсуждения).
	Author   string `json:"author"`   // Автор, написавший данное сообщение.
	Message  string `json:"message"`  // Собственно сообщение форума.
	IsEdited bool   `json:"isEdited"` // Истина, если данное сообщение было изменено.
	Forum    string `json:"forum"`    // Идентификатор форума (slug) данного сообещния.
	Thread   int32  `json:"thread"`   // Идентификатор ветви (id) обсуждения данного сообещния.
	Created  string `json:"created"`  // Дата создания сообщения на форуме.

}

type PostUpdate struct {
	ID int64 `json:"-"`
	Message  string `json:"message"` // Собственно сообщение форума.
}

type PostFull struct {
	Post   Post   `json:"post"`
	Author User   `json:"author"`
	Thread Thread `json:"thread"`
	Forum  Forum  `json:"forum"`
}

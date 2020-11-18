package models

type Forum struct {
	Title string `json:"title"` // Название форума.
	User string `json:"user"` // Nickname пользователя, который отвечает за форум.
	Slug string `json:"slug"` // Человекопонятный URL
	Posts int64 `json:"posts"` // Общее кол-во сообщений в данном форуме.
	Threads int32 `json:"threads"` // Общее кол-во ветвей обсуждения в данном форуме.
}

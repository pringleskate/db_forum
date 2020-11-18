package models

type Thread struct {
	Id int32 `json:"id"` // Идентификатор ветки обсуждения.
	Title string `json:"title"` // Заголовок ветки обсуждения.
	Author string `json:"author"` //  Пользователь, создавший данную тему.
	Forum string `json:"forum"` // Форум, в котором расположена данная ветка обсуждения.
	Message string `json:"message"` // Описание ветки обсуждения.
	Votes int32 `json:"votes"` // Кол-во голосов непосредственно за данное сообщение форума.
	Slag string `json:"slag"` // Человекопонятный URL
	Created string `json:"created"` // Дата создания ветки на форуме.
}

type ThreadUpdate struct {
	Title string `json:"title"` // Заголовок ветки обсуждения.
	Message string `json:"message"` // Описание ветки обсуждения.
}

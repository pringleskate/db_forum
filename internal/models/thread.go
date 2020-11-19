package models

type Thread struct {
	ForumSlug string `json:"-"`
	Id        int32  `json:"id"`      // Идентификатор ветки обсуждения.
	Title     string `json:"title"`   // Заголовок ветки обсуждения.
	Author    string `json:"author"`  // Пользователь, создавший данную тему.
	Forum     string `json:"forum"`   // Форум, в котором расположена данная ветка обсуждения.
	Message   string `json:"message"` // Описание ветки обсуждения.
	Votes     int32  `json:"votes"`   // Кол-во голосов непосредственно за данное сообщение форума.
	Slag      string `json:"slag"`    // Человекопонятный URL
	Created   string `json:"created"` // Дата создания ветки на форуме.
}

type ThreadUpdate struct {
	SlagOrID string `json:"-"`
	Title    string `json:"title"`   // Заголовок ветки обсуждения.
	Message  string `json:"message"` // Описание ветки обсуждения.
}

type ThreadQueryParams struct {
	Limit int    // Максимальное кол-во возвращаемых записей.
	Since string // Идентификатор поста, после которого будут выводиться записи (пост с данным идентификатором в результат не попадает)
	Sort  string
	Desc  bool   // Флаг сортировки по убыванию
}

type Vote struct {
	SlagOrID string `json:"-"`
	Nickname string `json:"nickname"` // Идентификатор пользователя.
	Voice    int32  `json:"voice"`    // Отданный голос.
}

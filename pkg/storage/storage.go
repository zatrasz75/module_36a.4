package db

// Post Публикация, получаемая из RSS.
type Post struct {
	ID      int    // Номер записи
	Title   string // Заголовок публикации
	Content string // Содержание публикации
	PubTime int64  // Время публикации
	Link    string // Ссылка на источник
}

// Interface задаёт контракт на работу с БД
type Interface interface {
	Posts(n int) ([]Post, error) // Получение последних публикаций
	AddPost(p Post) error        // Добавление новой публикации в базу
}

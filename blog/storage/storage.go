package storage

type Post struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	Author      string `db:"author"`
	Description string `db:"description"`
	Category    string `db:"category"`
}

type Category struct {
	CategoryID   int64  `db:"category_id"`
	CategoryName string `db:"category_name"`
}

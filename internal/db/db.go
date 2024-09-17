package db

type Database interface {
	SaveLink(id, link string)
	GetLink(id string) (string, bool)
}

type DB struct {
	Links map[string]string
}

func NewDB() *DB {
	return &DB{
		Links: make(map[string]string),
	}
}

func (db *DB) SaveLink(id, link string) {
	db.Links[id] = link
}

func (db *DB) GetLink(id string) (string, bool) {
	link, exists := db.Links[id]
	return link, exists
}

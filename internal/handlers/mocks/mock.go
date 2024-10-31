package mocks

type Database interface {
	SaveLink(id, link string)
	GetLink(id string) (string, bool)
}

type DB struct{}

func (m *DB) SaveLink(id, link string) {
}

func (m *DB) GetLink(id string) (link string, exists bool) {
	if id == "abc123" {
		link, exists = "https://example.com", true
	} else {
		link, exists = "", false
	}
	return
}

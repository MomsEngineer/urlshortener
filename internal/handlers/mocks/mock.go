package mocks

type Storage struct{}

type LinkStorage interface {
	SaveLink(id, link string)
	GetLink(id string) (string, bool)
}

func (s *Storage) SaveLink(id, link string) {
}

func (s *Storage) GetLink(id string) (link string, exists bool) {
	if id == "abc123" {
		link, exists = "https://example.com", true
	} else {
		link, exists = "", false
	}
	return
}

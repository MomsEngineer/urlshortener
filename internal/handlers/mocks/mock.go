package mocks

type Storage struct{}

type LinkStorage interface {
	Ping() error
	SaveLink(id, link string) error
	GetLink(id string) (string, bool, error)
}

func (s *Storage) Ping() error {
	return nil
}

func (s *Storage) SaveLink(id, link string) error {
	return nil
}

func (s *Storage) GetLink(id string) (link string, exists bool, err error) {
	if id == "abc123" {
		link, exists = "https://example.com", true
	} else {
		link, exists = "", false
	}
	return
}

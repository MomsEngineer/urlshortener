package memory

type LinksMap struct {
	Links map[string]string
}

func NewLinksMap() *LinksMap {
	return &LinksMap{
		Links: make(map[string]string),
	}
}

func (lm *LinksMap) SaveLink(id, link string) {
	lm.Links[id] = link
}

func (lm *LinksMap) GetLink(id string) (string, bool) {
	link, exists := lm.Links[id]
	return link, exists
}

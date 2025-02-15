package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
}

type ArticleStatus uint8

const (
	ArticleStatusUnknown = iota
	ArticleStatusUnPublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

func (status *ArticleStatus) String() string {
	switch *status {
	case ArticleStatusUnPublished:
		return "UnPublished"
	case ArticleStatusPublished:
		return "Published"
	case ArticleStatusPrivate:
		return "Private"
	default:
		return "Unknown"
	}
}

func (status *ArticleStatus) isPublished() bool {
	return *status == ArticleStatusPublished
}

type Author struct {
	Id   int64
	Name string
}

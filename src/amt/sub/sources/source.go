package sources

type Source interface {
	Search(domain string, timeOut int) ([]string, error)
	GetName() string
}
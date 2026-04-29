package sources

type Source interface {
	Search(domain string) ([]string, error)
	GetName() string
}
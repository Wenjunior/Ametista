package sources

import (
	"time"
)

type Source interface {
	Search(domain string, timeOut time.Duration) ([]string, error)
	GetName() string
}
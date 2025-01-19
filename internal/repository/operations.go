package repository

type RateLimitRepository interface {
	IncrementCounter(key string) (int, error)
	GetCounter(key string) (int, error)
	SetExpiration(key string, seconds int) error
	KeyExists(key string) (bool, error)
	Close() error
}

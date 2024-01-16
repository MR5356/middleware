package cache

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error

	TryLock(key string) error
	Unlock(key string) error
	Subscribe(key string, callback interface{}) error
	Publish(key string, data interface{})
}

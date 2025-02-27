package ktypes

// IKernelStoreKv -- интерфейс к локальному быстрому key-value хранилищу ядра
type IKernelStoreKv interface {
	// Get -- возвращает значение по ключу
	Get(key string) ([]byte, error)
	// Set -- устанавливает значение по ключу
	Set(key string, val []byte) error
	// Delete -- удаляет значение по ключу
	Delete(key string) error
	// Log -- возвращает локальный лог
	Log() ILogBuf
}

package kernel_types

// IKernelStore -- интерфейс к локальному хранилищу ядра
type IKernelStore interface {
	// Get -- возвращает значение по ключу
	Get(key string) ([]byte, error)
	// Set -- устанавливает значение по ключу
	Set(key string, val []byte) error
	// Delete -- удаляет значение по ключу
	Delete(key string) error
}

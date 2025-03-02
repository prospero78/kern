package ktypes

// IRWMutex -- интерфейс RW-мьютекса
type IRWMutex interface {
	// Lock -- блокирует контекст ядра
	Lock()
	// Unlock -- разблокирует контекст ядра
	Unlock()
	// RLock -- блокирует контекст ядра только для чтения
	RLock()
	// RUnlock -- разблокирует контекст ядра только на чтение
	RUnlock()
}

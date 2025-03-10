// package safe_string -- потокобезопасная строка
package safe_string

import (
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// safeString -- потокобезопасная строка
type safeString struct {
	sync.RWMutex
	val string
}

// NewSafeString -- возвращает новую потокобезопасную строку
func NewSafeString() ISafeString {
	sf := &safeString{}
	return sf
}

// Byte -- возвращает байтовое представление строки
func (sf *safeString) Byte() []byte {
	sf.RLock()
	defer sf.RUnlock()
	return []byte(sf.val)
}

// Get -- возвращает хранимую строку
func (sf *safeString) Get() string {
	sf.RLock()
	defer sf.RUnlock()
	return sf.val
}

// Set -- устанавливает строку
func (sf *safeString) Set(val string) {
	sf.Lock()
	defer sf.Unlock()
	sf.val = val
}

// Reset -- сбрасывает строку
func (sf *safeString) Reset() {
	sf.Lock()
	defer sf.Unlock()
	sf.val = ""
}

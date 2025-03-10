// package safe_int -- потокобезопасный целое
package safe_int

import (
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// safeInt -- потокобезопасное целое
type safeInt struct {
	sync.RWMutex
	val int
}

// NewSafeInt -- возвращает новое потокобезопасное целое
func NewSafeInt() ISafeInt {
	sf := &safeInt{}
	return sf
}

// Get -- возвращает хранимое целое
func (sf *safeInt) Get() int {
	sf.RLock()
	defer sf.RUnlock()
	return sf.val
}

// Set -- устанавливает целое
func (sf *safeInt) Set(val int) {
	sf.Lock()
	defer sf.Unlock()
	sf.val = val
}

// Reset -- сбрасывает целое
func (sf *safeInt) Reset() {
	sf.Lock()
	defer sf.Unlock()
	sf.val = 0
}

// package safe_bool -- потокобезопасный булевый признак
package safe_bool

import (
	"sync"

	. "github.com/prospero78/kern/krn/ktypes"
)

// safeBool -- потокобезопасный булевый признак
type safeBool struct {
	sync.RWMutex
	val bool
}

// NewSafeBool -- возвращает новый потокобезопасный булевый признак
func NewSafeBool() ISafeBool {
	sf := &safeBool{}
	return sf
}

// Get -- возвращает хранимый булевый признак
func (sf *safeBool) Get() bool {
	sf.RLock()
	defer sf.RUnlock()
	return sf.val
}

// Set -- устанавливает булевый признак
func (sf *safeBool) Set() {
	sf.Lock()
	defer sf.Unlock()
	sf.val = true
}

// Reset -- сбрасывает булевый признак
func (sf *safeBool) Reset() {
	sf.Lock()
	defer sf.Unlock()
	sf.val = false
}

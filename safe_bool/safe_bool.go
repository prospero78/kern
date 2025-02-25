// package safe_bool -- потокобезопасный булевый признак
package safe_bool

import (
	"sync"

	. "github.com/prospero78/kern/kernel_types"
)

// safeBool -- потокобезопасный булевый признак
type safeBool struct {
	val   bool
	block sync.RWMutex
}

// NewSafeBool -- возвращает новый потокобезопасный булевый признак
func NewSafeBool() ISafeBool {
	sf := &safeBool{}
	return sf
}

// Get -- возвращает хранимый булевый признак
func (sf *safeBool) Get() bool {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// Set -- устанавливает булевый признак
func (sf *safeBool) Set() {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.val = true
}

// Reset -- сбрасывает булевый признак
func (sf *safeBool) Reset() {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.val = false
}

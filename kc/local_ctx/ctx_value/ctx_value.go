// package ctx_value -- потокобезопасное значение локального контекста
package ctx_value

import (
	"sync"
	"time"

	. "github.com/prospero78/kern/krn/kalias"
	. "github.com/prospero78/kern/krn/ktypes"
)

// ctxValue -- потокобезопасное значение локального контекста
type ctxValue struct {
	val      any
	createAt ATime
	updateAt ATime
	comment  string
	block    sync.RWMutex
}

// NewCtxValue -- возвращает новое потокобезопасное значение локального контекста
func NewCtxValue(val any, comment string) ICtxValue {
	sf := &ctxValue{
		val:      val,
		comment:  comment,
		createAt: ATime(time.Now().Local().String()),
	}
	return sf
}

// Update -- обновляет хранимое значение
func (sf *ctxValue) Update(val any, comment string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.val = val
	sf.comment = comment
	sf.updateAt = ATime(time.Now().Local().String())
}

// Val -- возвращает хранимое значение
func (sf *ctxValue) Val() any {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// UpdateAt -- возвращает время обновления значения
func (sf *ctxValue) UpdateAt() ATime {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.updateAt
}

// CreateAt -- возвращает время создания значения
func (sf *ctxValue) CreateAt() ATime {
	return sf.createAt
}

// Comment -- возвращает комментарий значения
func (sf *ctxValue) Comment() string {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.comment
}

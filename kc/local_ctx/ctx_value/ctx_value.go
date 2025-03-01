// package ctx_value -- потокобезопасное значение локального контекста
package ctx_value

import (
	"fmt"
	"sync"

	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/kalias"
	. "github.com/prospero78/kern/krn/ktypes"
)

// ctxValue -- потокобезопасное значение локального контекста
type ctxValue struct {
	key      string
	val      any
	createAt ATime
	updateAt ATime
	comment  string
	block    sync.RWMutex
}

// NewCtxValue -- возвращает новое потокобезопасное значение локального контекста
func NewCtxValue(key string, val any, comment string) ICtxValue {
	Hassert(key != "", "NewCtxValue(): key is empty")
	sf := &ctxValue{
		key:      key,
		val:      val,
		comment:  comment,
		createAt: TimeNow(),
	}
	return sf
}

// Update -- обновляет хранимое значение
func (sf *ctxValue) Update(val any, comment string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.val = val
	sf.comment = comment
	sf.updateAt = TimeNow()
}

// Key -- возвращает ключ значения
func (sf *ctxValue) Key() string {
	return sf.key
}

// ValStr -- возвращает строковое представление значения
func (sf *ctxValue) ValStr() string {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return fmt.Sprint(sf.val)
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

// package ctx_value -- потокобезопасное значение локального контекста
package ctx_value

import (
	"sync"

	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/kalias"
	. "github.com/prospero78/kern/krn/ktypes"
)

// CtxValue -- потокобезопасное значение локального контекста
type CtxValue struct {
	sync.RWMutex
	key       string
	createAt  ATime
	Val_      any
	UpdateAt_ ATime
	Comment_  string
}

// NewCtxValue -- возвращает новое потокобезопасное значение локального контекста
func NewCtxValue(key string, val any, comment string) ICtxValue {
	Hassert(key != "", "NewCtxValue(): key is empty")
	sf := &CtxValue{
		key:      key,
		createAt: TimeNowStr(),
		Val_:     val,
		Comment_: comment,
	}
	return sf
}

// Key -- возвращает ключ значения
func (sf *CtxValue) Key() string {
	return sf.key
}

// Val -- возвращает хранимое значение
func (sf *CtxValue) Val() any {
	sf.RLock()
	defer sf.RUnlock()
	return sf.Val_
}

// UpdateAt -- возвращает время обновления значения
func (sf *CtxValue) UpdateAt() ATime {
	sf.RLock()
	defer sf.RUnlock()
	return sf.UpdateAt_
}

// CreateAt -- возвращает время создания значения
func (sf *CtxValue) CreateAt() ATime {
	return sf.createAt
}

// Comment -- возвращает комментарий значения
func (sf *CtxValue) Comment() string {
	sf.RLock()
	defer sf.RUnlock()
	return sf.Comment_
}

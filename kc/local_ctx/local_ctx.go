// package local_ctx -- локальный контекст
package local_ctx

import (
	"context"
	"sync"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/local_ctx/ctx_value"
	. "github.com/prospero78/kern/krn/ktypes"
)

// localCtx -- локальный контекст
type localCtx struct {
	ctx      context.Context      // Отменяемый контекст
	fnCancel func()               // Функция отмены контекста
	dictVal  map[string]ICtxValue // Словарь различных значений
	block    sync.RWMutex
}

// NewLocalCtx -- возвращает новый локальный контекст
func NewLocalCtx(ctx context.Context) ILocalCtx {
	Hassert(ctx != nil, "NewLocalCtx(): ctx==nil")
	_ctx, fnCancel := context.WithCancel(ctx)
	sf := &localCtx{
		ctx:      _ctx,
		fnCancel: fnCancel,
		dictVal:  map[string]ICtxValue{},
	}
	return sf
}

// Get -- возвращает хранимое значение
func (sf *localCtx) Get(key string) ICtxValue {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.dictVal[key]
}

// Cancel -- отменяет контекст
func (sf *localCtx) Cancel() {
	sf.fnCancel()
}

// Del -- удаляет значение из контекста
func (sf *localCtx) Del(key string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	delete(sf.dictVal, key)
}

// Set -- добавляет значение в контекст
func (sf *localCtx) Set(key string, val any, comment string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	_val, isOk := sf.dictVal[key]
	if isOk {
		_val.Update(val, comment)
		return
	}
	_val = ctx_value.NewCtxValue(val, comment)
	sf.dictVal[key] = _val
}

// Done -- блокирующий вызов ожидания отмены контекста
func (sf *localCtx) Done() {
	<-sf.ctx.Done()
}

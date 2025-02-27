// package local_ctx -- локальный контекст
package local_ctx

import (
	"context"
	"sync"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/local_ctx/ctx_value"
	"github.com/prospero78/kern/kc/log_buf"
	. "github.com/prospero78/kern/krn/ktypes"
)

// localCtx -- локальный контекст
type localCtx struct {
	ctx      context.Context      // Отменяемый контекст
	fnCancel func()               // Функция отмены контекста
	dictVal  map[string]ICtxValue // Словарь различных значений
	log      ILogBuf              // Локальный буфер
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
		log:      log_buf.NewLogBuf(),
	}
	return sf
}

// Log -- возвращает локальный буферный лог
func (sf *localCtx) Log() ILogBuf {
	return sf.log
}

// Get -- возвращает хранимое значение
func (sf *localCtx) Get(key string) ICtxValue {
	sf.block.RLock()
	defer sf.block.RUnlock()
	sf.log.Debug("localCtx.Get(): key='%v'", key)
	return sf.dictVal[key]
}

// Cancel -- отменяет контекст
func (sf *localCtx) Cancel() {
	sf.log.Warn("localCtx.Cancel()")
	sf.fnCancel()
}

// Del -- удаляет значение из контекста
func (sf *localCtx) Del(key string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.log.Debug("localCtx.Del(): key='%v'", key)
	delete(sf.dictVal, key)
}

// Set -- добавляет значение в контекст
func (sf *localCtx) Set(key string, val any, comment string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.log.Debug("localCtx.Set(): key='%v'", key)
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
	sf.log.Debug("localCtx.Done(): done")
}

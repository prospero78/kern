// package local_ctx -- локальный контекст
package local_ctx

import (
	"context"
	"sync"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/local_ctx/ctx_value"
	"github.com/prospero78/kern/kc/local_ctx/lst_sort"
	"github.com/prospero78/kern/kc/log_buf"
	. "github.com/prospero78/kern/krn/ktypes"
)

// LocalCtx -- локальный контекст
type LocalCtx struct {
	sync.RWMutex
	ctx      context.Context // Отменяемый контекст
	fnCancel func()          // Функция отмены контекста

	dictVal map[string]ICtxValue // Словарь различных значений
	lstSort *lst_sort.LstSort    // Сортированный список значений
	log     ILogBuf              // Локальный буфер
}

// NewLocalCtx -- возвращает новый локальный контекст
func NewLocalCtx(ctx context.Context) ILocalCtx {
	Hassert(ctx != nil, "NewLocalCtx(): ctx==nil")
	_ctx, fnCancel := context.WithCancel(ctx)
	sf := &LocalCtx{
		ctx:      _ctx,
		fnCancel: fnCancel,
		dictVal:  map[string]ICtxValue{},
		lstSort:  lst_sort.NewLstSort(),
		log:      log_buf.NewLogBuf(),
	}
	return sf
}

// Size -- возвращает размер контекста
func (sf *LocalCtx) Size() int {
	sf.RLock()
	defer sf.RUnlock()
	return len(sf.dictVal)
}

// SortedList -- возвращает сортированный список значений
func (sf *LocalCtx) SortedList() []ICtxValue {
	return sf.lstSort.List()
}

// Log -- возвращает локальный буферный лог
func (sf *LocalCtx) Log() ILogBuf {
	return sf.log
}

// Get -- возвращает хранимое значение
func (sf *LocalCtx) Get(key string) ICtxValue {
	sf.RLock()
	defer sf.RUnlock()
	Hassert(key != "", "localCtx.Get(): key is empty")
	sf.log.Debug("localCtx.Get(): key='%v'", key)
	return sf.dictVal[key]
}

// Del -- удаляет значение из контекста
func (sf *LocalCtx) Del(key string) {
	sf.Lock()
	defer sf.Unlock()
	sf.log.Debug("localCtx.Del(): key='%v'", key)
	val := sf.dictVal[key]
	delete(sf.dictVal, key)
	sf.lstSort.Del(val)
}

// Set -- добавляет значение в контекст
func (sf *LocalCtx) Set(key string, val any, comment string) {
	sf.Lock()
	defer sf.Unlock()
	sf.log.Debug("localCtx.Set(): key='%v'", key)
	_val, isOk := sf.dictVal[key]
	if isOk {
		val0 := _val.(*ctx_value.CtxValue)
		val0.Lock()
		val0.UpdateAt_ = TimeNow()
		val0.Val_ = val
		val0.Unlock()
		return
	}
	_val = ctx_value.NewCtxValue(key, val, comment)
	sf.dictVal[key] = _val
	sf.lstSort.Add(_val)
}

// Done -- блокирующий вызов ожидания отмены контекста
func (sf *LocalCtx) Done() {
	<-sf.ctx.Done()
}

// Cancel -- отменяет контекст
func (sf *LocalCtx) Cancel() {
	sf.log.Warn("localCtx.Cancel()")
	sf.fnCancel()
}

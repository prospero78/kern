// package local_ctx -- локальный контекст
package local_ctx

import (
	"context"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/local_ctx/ctx_value"
	"github.com/prospero78/kern/kc/local_ctx/lst_sort"
	"github.com/prospero78/kern/kc/log_buf"
	. "github.com/prospero78/kern/krn/ktypes"
)

// LocalCtx -- локальный контекст
type LocalCtx struct {
	ctx      context.Context // Отменяемый контекст
	fnCancel func()          // Функция отмены контекста

	chGetIn  chan string
	chGetOut chan ICtxValue

	chDelIn  chan string
	chDelOut chan int

	chSetIn  chan triple
	chSetOut chan int

	chSizeIn  chan int
	chSizeOut chan int

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

		chGetIn:  make(chan string, 2),
		chGetOut: make(chan ICtxValue, 2),

		chDelIn:  make(chan string, 2),
		chDelOut: make(chan int, 2),

		chSetIn:  make(chan triple, 2),
		chSetOut: make(chan int, 2),

		chSizeIn:  make(chan int, 2),
		chSizeOut: make(chan int, 2),

		dictVal: map[string]ICtxValue{},
		lstSort: lst_sort.NewLstSort(),
		log:     log_buf.NewLogBuf(),
	}
	go sf.run()
	return sf
}

// Size -- возвращает размер контекста
func (sf *LocalCtx) Size() int {
	sf.chSizeIn <- 1
	return <-sf.chSizeOut
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
	Hassert(key != "", "localCtx.Get(): key is empty")
	sf.log.Debug("localCtx.Get(): key='%v'", key)
	sf.chGetIn <- key
	return <-sf.chGetOut
}

// Del -- удаляет значение из контекста
func (sf *LocalCtx) Del(key string) {
	sf.log.Debug("localCtx.Del(): key='%v'", key)
	sf.chDelIn <- key
	<-sf.chDelOut
}

type triple struct {
	key     string
	val     any
	comment string
}

// Set -- добавляет значение в контекст
func (sf *LocalCtx) Set(key string, val any, comment string) {
	sf.log.Debug("localCtx.Set(): key='%v'", key)
	valIn := triple{
		key:     key,
		val:     val,
		comment: comment,
	}
	sf.chSetIn <- valIn
	<-sf.chSetOut
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

func (sf *LocalCtx) run() {
	for {
		select {
		case <-sf.ctx.Done():
			sf.log.Debug("localCtx.run(): done")
			return
		case key := <-sf.chGetIn:
			sf.chGetOut <- sf.dictVal[key]
		case key := <-sf.chDelIn:
			val := sf.dictVal[key]
			delete(sf.dictVal, key)
			sf.lstSort.Del(val)
			sf.chDelOut <- 1
		case valIn := <-sf.chSetIn:
			sf.set(valIn)
			sf.chSetOut <- 1
		case <-sf.chSizeIn:
			sf.chSizeOut <- len(sf.dictVal)
		}
	}
}

func (sf *LocalCtx) set(val triple) {
	_val, isOk := sf.dictVal[val.key]
	if isOk {
		_val.Update(val.val, val.comment)
		return
	}
	_val = ctx_value.NewCtxValue(val.key, val.val, val.comment)
	sf.dictVal[val.key] = _val
	sf.lstSort.Add(_val)
}

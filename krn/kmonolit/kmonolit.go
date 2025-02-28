// package kmonolit -- модульный монолит на основе ядра
package kmonolit

import (
	"sync"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/local_ctx"
	"github.com/prospero78/kern/kc/safe_bool"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

// kMonolit -- объект модульного монолита
type kMonolit struct {
	kCtx    IKernelCtx
	ctx     ILocalCtx
	log     ILogBuf
	name    string
	isLocal bool
	isWork  ISafeBool
	isEnd   bool
	dict    map[AModuleName]IKernelModule // Словарь модулей монолита
	block   sync.Mutex
}

var (
	mon *kMonolit
)

// GetMonolit -- возвращает монолит
func GetMonolit(name string) IKernelMonolit {
	if mon != nil {
		return mon
	}
	Hassert(name != "", "NewMonolit(): name is empty")
	kCtx := kctx.GetKernelCtx()
	sf := &kMonolit{
		kCtx:    kCtx,
		ctx:     local_ctx.NewLocalCtx(kCtx.BaseCtx()),
		name:    name,
		dict:    map[AModuleName]IKernelModule{},
		isWork:  safe_bool.NewSafeBool(),
		isLocal: kCtx.Get("isLocal").Val().(bool),
	}
	sf.log = sf.ctx.Log()
	kCtx.Set("monolitName", name, "name of monolit")
	mon = sf
	return sf
}

// Log -- возвращает лог монолита
func (sf *kMonolit) Log() ILogBuf {
	return sf.ctx.Log()
}

// Name -- возвращает имя монолита
func (sf *kMonolit) Name() string {
	return sf.name
}

// Add -- добавляет модуль в монолит
func (sf *kMonolit) Add(module IKernelModule) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(module != nil, "kMonolit.Add(): module==nil")
	_, isOk := sf.dict[module.Name()]
	Hassert(!isOk, "kMonolit.Add(): module(%v) already exists", module.Name())
	sf.dict[module.Name()] = module
	sf.log.Debug("kMonolit.Add(): module='%v'", module.Name())
	if sf.isWork.Get() {
		go module.Run()
		sf.log.Debug("kMonolit.Add(): module='%v' is run", module.Name())
	}
}

// Run -- запускает монолит в работу
func (sf *kMonolit) Run() {
	sf.block.Lock()
	defer sf.block.Unlock()
	if sf.isEnd {
		return
	}
	sf.isWork.Set()
	for _, module := range sf.dict {
		go module.Run()
	}
	sf.log.Debug("kMonolit.Run()")
}

// IsLocal -- возвращает признак локальной шины
func (sf *kMonolit) IsLocal() bool {
	return sf.isLocal
}

// IsWork -- возвращает признак работы монолита
func (sf *kMonolit) IsWork() bool {
	return sf.isWork.Get()
}

// Ожидание завершения работы монолита
func (sf *kMonolit) Wait() {
	sf.kCtx.Done()
	sf.kCtx.Wg().Wait()
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.isWork.Reset()
	sf.isEnd = true
	sf.log.Debug("kMonolit.close(): end")
}

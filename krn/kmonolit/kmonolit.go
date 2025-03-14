// package kmonolit -- модульный монолит на основе ядра
package kmonolit

import (
	"fmt"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	"gitp78su.ipnodns.ru/svi/kern/kc/local_ctx"
	"gitp78su.ipnodns.ru/svi/kern/kc/safe_bool"
	. "gitp78su.ipnodns.ru/svi/kern/krn/kalias"
	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// kMonolit -- объект модульного монолита
type kMonolit struct {
	kCtx    IKernelCtx
	ctx     ILocalCtx
	log     ILogBuf
	name    string
	isLocal bool
	isWork  ISafeBool
	isEnd   ISafeBool
	dict    map[AModuleName]IKernelModule // Словарь модулей монолита
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
		isEnd:   safe_bool.NewSafeBool(),
		isLocal: kCtx.Get("isLocal").Val().(bool),
	}
	sf.log = sf.ctx.Log()
	sf.kCtx.Set("monolitName", name, "name of monolit")
	sf.kCtx.Set("monolit", sf, "monolit-app")
	sf.ctx.Set("monolitName", name, "name of monolit")
	mon = sf
	return sf
}

// Ctx -- возвращает контекст монолита
func (sf *kMonolit) Ctx() ILocalCtx {
	return sf.ctx
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
	sf.kCtx.RLock()
	defer sf.kCtx.RUnlock()
	Hassert(module != nil, "kMonolit.Add(): module==nil")
	_, isOk := sf.dict[module.Name()]
	Hassert(!isOk, "kMonolit.Add(): module(%v) already exists", module.Name())
	sf.dict[module.Name()] = module
	sf.log.Debug("kMonolit.Add(): module='%v'", module.Name())
	if sf.isWork.Get() {
		go module.Run()
		sf.log.Debug("kMonolit.Add(): module='%v' is run", module.Name())
	}
	key := fmt.Sprintf("module_%v", len(sf.dict))
	moduleName := string(module.Name())
	sf.ctx.Set(key, module, "kMonolit.Add(): module="+moduleName)
}

// Run -- запускает монолит в работу
func (sf *kMonolit) Run() {
	sf.kCtx.RLock()
	defer sf.kCtx.RUnlock()
	if sf.isEnd.Get() {
		return
	}
	if sf.isWork.Get() {
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
	sf.kCtx.Lock()
	defer sf.kCtx.Unlock()
	sf.isWork.Reset()
	sf.isEnd.Set()
	sf.log.Debug("kMonolit.close(): end")
}

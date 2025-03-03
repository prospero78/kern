// package kmodule -- модуль на основе ядра
package kmodule

import (
	"time"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/local_ctx"
	"github.com/prospero78/kern/kc/safe_int"
	"github.com/prospero78/kern/kc/safe_string"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kbus/kbus_local"
	"github.com/prospero78/kern/krn/kctx"
	"github.com/prospero78/kern/krn/kmodule/mod_stat"
	. "github.com/prospero78/kern/krn/ktypes"
)

// kModule -- модуль на основе ядра
type kModule struct {
	kCtx      IKernelCtx
	ctx       ILocalCtx
	name      AModuleName
	bus       IKernelBus
	timePhase ISafeInt
	strLive   ISafeString
	stat      IModuleStat
}

// NewKernelModule -- возвращает новый модуль на основе ядра
func NewKernelModule(name AModuleName) IKernelModule {
	Hassert(name != "", "NewKernelModule(): name is empty")
	kCtx := kctx.GetKernelCtx()
	sf := &kModule{
		kCtx:      kCtx,
		ctx:       local_ctx.NewLocalCtx(kCtx.BaseCtx()),
		name:      name,
		bus:       kbus_local.GetKernelBusLocal(),
		timePhase: safe_int.NewSafeInt(),
		strLive:   safe_string.NewSafeString(),
		stat:      mod_stat.NewModStat(name),
	}
	sf.timePhase.Set(1000) // 1000 msec
	go sf.sigLive()
	return sf
}

// Stat -- возвращает статистику модуля
func (sf *kModule) Stat() IModuleStat {
	return sf.stat
}

// Log -- возвращает буферный лог
func (sf *kModule) Log() ILogBuf {
	return sf.ctx.Log()
}

// Ctx -- возвращает контекст модуля
func (sf *kModule) Ctx() ILocalCtx {
	return sf.ctx
}

// Run -- запускает модуль в работу
func (sf *kModule) Run() {
	Hassert(false, "kModule.Run(): module='%v', parent not realised this method", sf.name)
}

// Name -- возвращает уникальное имя модуля
func (sf *kModule) Name() AModuleName {
	return sf.name
}

// IsWork -- возвращает признак состояния работы
func (sf *kModule) IsWork() bool {
	Hassert(false, "kModule.IsWork(): module='%v', parent not realised this method", sf.name)
	return false
}

// Live -- возвращает индикатор жизни модуля
func (sf *kModule) Live() string {
	return sf.strLive.Get()
}

// Сигнал жизни, каждые 5 сек публикует в шину метку
func (sf *kModule) sigLive() {
	var (
		topic  = sf.name + "_live"
		iPhase = 0
		err    error
	)
	fnPhase := func() {
		switch iPhase {
		case 0:
			sf.strLive.Set("|")
			err = sf.bus.Publish(ATopic(topic), sf.strLive.Byte())
		case 1:
			sf.strLive.Set("/")
			err = sf.bus.Publish(ATopic(topic), sf.strLive.Byte())
		case 2:
			sf.strLive.Set("-")
			err = sf.bus.Publish(ATopic(topic), sf.strLive.Byte())
		case 3:
			sf.strLive.Set("\\")
			err = sf.bus.Publish(ATopic(topic), sf.strLive.Byte())
			iPhase = -1
		}
		sf.recErr(err)
		iPhase++
		sf.stat.Add(1)
		time.Sleep(time.Millisecond * time.Duration(sf.timePhase.Get()))
	}
	for {
		select {
		case <-sf.kCtx.BaseCtx().Done():
			return
		default:
			fnPhase()
		}
	}
}

// Регистрирует ошибку обработчика при публикации лайв сигнала, если была
func (sf *kModule) recErr(err error) {
	if err != nil {
		sf.Log().Err("kModule.recErr(): name=%v, in publish live, err=\n\t%v", err)
	}
}

// package kernel_keeper -- сторож системных сигналов
package kernel_keeper

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	"gitp78su.ipnodns.ru/svi/kern/kc/log_buf"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// kernelKeeper -- сторож системных сигналов
type kernelKeeper struct {
	ctx      context.Context
	fnCancel func()
	wg       IKernelWg
	log      ILogBuf
	chSys_   chan os.Signal
}

var (
	kernKeep *kernelKeeper
	block    sync.Mutex
)

// GetKernelKeeper -- возвращает новый сторож системных сигналов
func GetKernelKeeper(ctx context.Context, fnCancel func(), wg IKernelWg) *kernelKeeper {
	block.Lock()
	defer block.Unlock()
	if kernKeep != nil {
		kernKeep.log.Debug("GetKernelKeeper()")
		return kernKeep
	}
	Hassert(ctx != nil, "NewKernelCtx(): ctx==nil")
	Hassert(wg != nil, "NewKernelCtx(): IKernelWg==nil")
	Hassert(fnCancel != nil, "NewKernelCtx(): fnCancel==nil")
	sf := &kernelKeeper{
		ctx:      ctx,
		fnCancel: fnCancel,
		wg:       wg,
		log:      log_buf.NewLogBuf(),
		chSys_:   make(chan os.Signal, 2),
	}
	sf.log.Debug("GetKernelKeeper(): first run")
	err := sf.wg.Add("kernel_keeper")
	Hassert(err == nil, "NewKernelCtx(): in add stream kernel keeper in IKernelWg, err=\n\t%v,err")

	go sf.run(sf.chSys_)
	kernKeep = sf
	_ = IKernelKeeper(sf)
	return sf
}

// Log -- возвращает лог сторожа системных сигналов
func (sf *kernelKeeper) Log() ILogBuf {
	return sf.log
}

// Работает в отдельном потоке и ждёт сигналов прерываний работы
func (sf *kernelKeeper) run(chSys chan os.Signal) {
	sf.log.Debug("kernelKeeper.run()")

	// Регистрируем сигналы SIGINT (Ctrl+C) и SIGTERM (завершение процесса)
	// syscall.SIGHUP: Сигнал, отправляемый при закрытии терминала.
	// syscall.SIGQUIT: Сигнал, отправляемый при нажатии **Ctrl+**.
	signal.Notify(chSys, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	select {
	case sig := <-chSys: // системный сигнал
		sf.log.Debug("kernelKeeper.run(): system signal, sig=%v\n", sig)
		sf.fnCancel()
	case <-sf.ctx.Done(): // сигнал от приложения
		sf.log.Debug("kernelKeeper.run(): cancel app context, err=\n\t%v\n", sf.ctx.Err())
	}
	sf.wg.Done("kernel_keeper")
	sf.log.Debug("kernelKeeper.run(): end")
}

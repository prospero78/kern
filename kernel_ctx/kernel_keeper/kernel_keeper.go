// package kernel_keeper -- сторож системных сигналов
package kernel_keeper

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_types"
)

// kernelKeeper -- сторож системных сигналов
type kernelKeeper struct {
	ctx      context.Context
	fnCancel func()
	wg       IKernelWg
	chSys_   chan os.Signal
}

var (
	kernKeep *kernelKeeper
)

// GetKernelKeeper -- возвращает новый сторож системных сигналов
func GetKernelKeeper(ctx context.Context, fnCancel func(), wg IKernelWg) *kernelKeeper {
	if kernKeep != nil {
		return kernKeep
	}
	Hassert(ctx != nil, "NewKernelCtx(): ctx==nil")
	Hassert(wg != nil, "NewKernelCtx(): IKernelWg==nil")
	Hassert(fnCancel != nil, "NewKernelCtx(): fnCancel==nil")
	sf := &kernelKeeper{
		ctx:      ctx,
		fnCancel: fnCancel,
		wg:       wg,
		chSys_:   make(chan os.Signal, 2),
	}
	err := sf.wg.Add("kernel_keeper")
	Hassert(err == nil, "NewKernelCtx(): in add stream kernel keeper in IKernelWg, err=\n\t%v,err")

	go sf.run(sf.chSys_)
	kernKeep = sf
	_ = IKernelKeeper(sf)
	return sf
}

// Работает в отдельном потоке и ждёт сигналов прерываний работы
func (sf *kernelKeeper) run(chSys chan os.Signal) {
	log.Println("kernelKeeper.run()")

	// Регистрируем сигналы SIGINT (Ctrl+C) и SIGTERM (завершение процесса)
	// syscall.SIGHUP: Сигнал, отправляемый при закрытии терминала.
	// syscall.SIGQUIT: Сигнал, отправляемый при нажатии **Ctrl+**.
	signal.Notify(chSys, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	select {
	case sig := <-chSys: // системный сигнал
		log.Printf("kernelKeeper.run(): system signal, sig=%v\n", sig)
		sf.fnCancel()
	case <-sf.ctx.Done(): // сигнал от приложения
		log.Printf("kernelKeeper.run(): cancel app context, err=\n\t%v\n", sf.ctx.Err())
	}
	sf.wg.Done("kernel_keeper")
	log.Printf("kernelKeeper.run(): end")
}

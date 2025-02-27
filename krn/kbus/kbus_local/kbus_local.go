// package kbus_local -- реализация локальной шины сообщений
package kbus_local

import (
	"github.com/prospero78/kern/krn/kbus/kbus_base"
	. "github.com/prospero78/kern/krn/ktypes"
)

// Локальная шина данных
type kernelBusLocal struct {
	*kbus_base.KernelBusBase
}

var (
	bus *kernelBusLocal
)

// GetKernelBusLocal -- возвращает локальную шину сообщений
func GetKernelBusLocal() IKernelBus {
	if bus != nil {
		return bus
	}
	bus = &kernelBusLocal{
		KernelBusBase: kbus_base.GetKernelBusBase(),
	}
	return bus
}

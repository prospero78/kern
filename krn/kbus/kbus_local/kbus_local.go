// package kbus_local -- реализация локальной шины сообщений
package kbus_local

import (
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/kbus_base"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// Локальная шина данных
type kernelBusLocal struct {
	*kbus_base.KBusBase
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
		KBusBase: kbus_base.GetKernelBusBase(),
	}
	return bus
}

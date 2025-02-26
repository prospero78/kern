// package kernel_bus_local -- реализация локальной шины сообщений
package kernel_bus_local

import (
	"github.com/prospero78/kern/kernel_bus/kernel_bus_base"
	. "github.com/prospero78/kern/kernel_types"
)

// Локальная шина данных
type kernelBusLocal struct {
	*kernel_bus_base.KernelBusBase
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
		KernelBusBase: kernel_bus_base.GetKernelBusBase(),
	}
	return bus
}

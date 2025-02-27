// package client_bus_local -- клиент локальной шины
package client_bus_local

import (
	"github.com/prospero78/kern/krn/kbus/kbus_local"
	. "github.com/prospero78/kern/krn/ktypes"
)

// ClientBusLocal -- клиент локальной шины
type ClientBusLocal struct {
	IKernelBus
}

// NewClientBusLocal -- клиент локальной шины
func NewClientBusLocal() IClientBus {
	sf := &ClientBusLocal{
		IKernelBus: kbus_local.GetKernelBusLocal(),
	}
	return sf
}

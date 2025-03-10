// package client_bus_local -- клиент локальной шины
package client_bus_local

import (
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/kbus_local"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// ClientBusLocal -- клиент локальной шины
type ClientBusLocal struct {
	IKernelBus
}

// NewClientBusLocal -- клиент локальной шины
func NewClientBusLocal() IBusClient {
	sf := &ClientBusLocal{
		IKernelBus: kbus_local.GetKernelBusLocal(),
	}
	return sf
}

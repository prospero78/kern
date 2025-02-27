// package client_bus_http -- клиент HTTP-шины
package client_bus_http

import (
	"fmt"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/local_ctx"
	"github.com/prospero78/kern/kc/safe_bool"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kbus/kbus_http"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

// ClientBusHttp -- клиент HTTP-шины
type ClientBusHttp struct {
	bus    IKernelBus
	ctx    ILocalCtx
	log    ILogBuf
	isWork ISafeBool
}

// NewClientBusHttp - -возвращает новый клиент HTTP-шины
func NewClientBusHttp(url string) *ClientBusHttp {
	Hassert(url != "", "NewClientBusHttp(): url is empty")
	kCtx := kctx.GetKernelCtx()
	sf := &ClientBusHttp{
		ctx:    local_ctx.NewLocalCtx(kCtx.BaseCtx()),
		bus:    kbus_http.GetKernelBusHttp(),
		isWork: safe_bool.NewSafeBool(),
	}
	sf.log = sf.ctx.Log()
	_ = IClientBus(sf)
	return sf
}

// Unsubscribe -- отписывается от топика в дистанционной шине
func (sf *ClientBusHttp) Unsubscribe(handler IBusHandlerSubscribe) {}

// Subscribe -- подписывается на топик в дистанционной шине
func (sf *ClientBusHttp) Subscribe(handler IBusHandlerSubscribe) error {
	return fmt.Errorf("not realised")
}

// SendRequest -- отправляет в дистанционную шину запрос
func (sf *ClientBusHttp) SendRequest(topic ATopic, binReq []byte) ([]byte, error) {
	return nil, fmt.Errorf("not realised")
}

// RegisterServe -- регистрирует в локальной шине обработчик
func (sf *ClientBusHttp) RegisterServe(handler IBusHandlerServe) {}

// Publish -- публикует сообщение в дистанционной шину
func (sf *ClientBusHttp) Publish(topic ATopic, binMsg []byte) error {
	return fmt.Errorf("not realised")
}

// Log -- возвращает локальный лог клиента
func (sf *ClientBusHttp) Log() ILogBuf {
	return sf.log
}

// IsWork -- возвращает признак работы
func (sf *ClientBusHttp) IsWork() bool {
	return sf.isWork.Get()
}

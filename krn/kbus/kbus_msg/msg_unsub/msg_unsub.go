// package msg_unsub -- сообщения отписки
package msg_unsub

import (
	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/kalias"
)

// UnsubReq -- запрос на отписку от топика
type UnsubReq struct {
	Name_ AHandlerName `json:"name"` // Уникальная метка подписки
	Uuid_ string       `json:"uuid"`
}

// SelfCheck -- проверка запроса на правильность полей
func (sf *UnsubReq) SelfCheck() {
	Hassert(sf.Name_ != "", "UnsubReq.SelfCheck(): name is empty")
	Hassert(sf.Uuid_ != "", "UnsubReq.SelfCheck(): uuid is empty")
}

// UnsubResp -- ответ на запрос отписки
type UnsubResp struct {
	Status_ string `json:"status"`
	Uuid_   string `json:"uuid"`
}

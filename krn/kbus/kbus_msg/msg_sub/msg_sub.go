// package msg_sub -- сообщения для подписки
package msg_sub

import (
	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/kalias"
)

// SubscribeReq -- входящий запрос на подписку
type SubscribeReq struct {
	Topic_   ATopic `json:"topic"` // Топик, на который надо подписаться
	Uuid_    string `json:"uuid"`
	WebHook_ string `json:"web_hook"` // Веб-хук для обратного вызова
}

// SelfCheck -- проверяет поля на правильность
func (sf *SubscribeReq) SelfCheck() {
	Hassert(sf.Topic_ != "", "SubscribeReq.SelfCheck(): topic is empty")
	Hassert(sf.Uuid_ != "", "SubscribeReq.SelfCheck(): uuid is empty")
	Hassert(sf.WebHook_ != "", "SubscribeReq.SelfCheck(): WebHook_ is empty")
}

// SubscribeResp -- ответ на запрос подписки
type SubscribeResp struct {
	Status_ string       `json:"status"`
	Uuid_   string       `json:"uuid"`
	Name_   AHandlerName `json:"name"` // Уникальное имя подписки
}

// SelfCheck -- проверяет правильность своих полей
func (sf *SubscribeResp) SelfCheck() {
	Hassert(sf.Status_ != "", "SubscribeResp.SelfCheck(): status is empty")
}

// package msg_serve -- сообщения на обслуживание входящих запросов
package msg_serve

import (
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/kalias"
)

// ServeReq -- входящий запрос на обслуживание
type ServeReq struct {
	Topic_  ATopic `json:"topic"`
	Uuid_   string `json:"uuid"`
	BinReq_ []byte `json:"req"`
}

// SelfCheck -- проверяет структуру на правильность полей
func (sf *ServeReq) SelfCheck() {
	Hassert(sf.Topic_ != "", "ServeReq.SelfCheck(): topic is empty")
	Hassert(sf.Uuid_ != "", "ServeReq.SelfCheck(): uuid is empty")
}

// ServeResp -- ответ на входящий запрос
type ServeResp struct {
	Status_  string `json:"status"`
	Uuid_    string `json:"uuid"`
	BinResp_ []byte `json:"resp"`
}

// SelfCheck -- проверяет правильность своих полей
func (sf *ServeResp) SelfCheck() {
	Hassert(sf.Status_ != "", "ServeResp.SelfCheck(): status is empty")
}

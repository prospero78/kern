// package msg_pub -- сообщения публикации
package msg_pub

import (
	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/kalias"
)

// PublishReq -- запрос на публикацию
type PublishReq struct {
	Topic_  ATopic `json:"topic"`
	Uuid_   string `json:"uuid"`
	BinMsg_ []byte `json:"msg"`
}

// SelfCheck -- проверяет правильность своих полей
func (sf *PublishReq) SelfCheck() {
	Hassert(sf.Topic_ != "", "PublishReq.SelfCheck(): topic is empty")
	Hassert(sf.Uuid_ != "", "PublishReq.SelfCheck(): uuid is empty")
}

// PublishResp -- ответ на запрос публикации
type PublishResp struct {
	Status_ string `json:"status"`
	Uuid_   string `json:"uuid"`
}

// SelfCheck -- проверяет правильность своих полей
func (sf *PublishResp) SelfCheck() {
	Hassert(sf.Status_ != "", "PublishResp.SelfCheck(): status is empty")
}

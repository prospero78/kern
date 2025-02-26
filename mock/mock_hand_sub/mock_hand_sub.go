// package mock_hand_sub -- мок-обработчик подписки
package mock_hand_sub

import (
	"crypto/rand"
	"sync"

	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_alias"
	. "github.com/prospero78/kern/kernel_types"
)

type MockHandlerSub struct {
	Msg_   []byte       // Для обратного вызова
	Name_  AHandlerName // Уникальное имя мок-обработчика подписки
	Topic_ ATopic       // Имя топика подписки
	block  sync.RWMutex
}

// NewMockHandlerSub -- возвращает новый обработчик подписки
func NewMockHandlerSub(topic ATopic, webHook string) *MockHandlerSub {
	Hassert(topic != "", "NewMockHandlerSub(): topic is empty")
	Hassert(webHook != "", "NewMockHandlerSub(): name is empty")
	sf := &MockHandlerSub{
		Topic_: topic,
		Name_:  AHandlerName(webHook + "_" + rand.Text()),
	}
	_ = IBusHandlerSubscribe(sf)
	return sf
}

// Возвращает хранимое значение
func (sf *MockHandlerSub) Msg() string {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return string(sf.Msg_)
}

// Функция обратного вызова подписки
func (sf *MockHandlerSub) FnBack(binMsg []byte) {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.Msg_ = binMsg
}

// Возвращает уникальное имя обработчика подписки
func (sf *MockHandlerSub) Topic() ATopic {
	return sf.Topic_
}

// Возвращает топик для обработчика подписки
func (sf *MockHandlerSub) Name() AHandlerName {
	return sf.Name_
}

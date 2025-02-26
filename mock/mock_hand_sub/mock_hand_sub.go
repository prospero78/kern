// package mock_hand_sub -- мок-обработчик подписки
package mock_hand_sub

import (
	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_alias"
	. "github.com/prospero78/kern/kernel_types"
)

type MockHandlerSub struct {
	Msg_   []byte // Для обратного вызова
	Name_  string // Уникальное имя мок-обработчика подписки
	Topic_ ATopic // Имя топика подписки
}

// NewMockHandlerSub -- возвращает новый обработчик подписки
func NewMockHandlerSub(topic ATopic, name string) *MockHandlerSub {
	Hassert(topic != "", "NewMockHandlerSub(): topic is empty")
	Hassert(name != "", "NewMockHandlerSub(): name is empty")
	sf := &MockHandlerSub{
		Topic_: topic,
		Name_:  name,
	}
	_ = IBusHandlerSubscribe(sf)
	return sf
}

// Функция обратного вызова подписки
func (sf *MockHandlerSub) FnBack(binMsg []byte) {
	sf.Msg_ = binMsg
}

// Возвращает уникальное имя обработчика подписки
func (sf *MockHandlerSub) Topic() ATopic {
	return sf.Topic_
}

// Возвращает топик для обработчика подписки
func (sf *MockHandlerSub) Name() string {
	return sf.Name_
}

// package mock_hand_sub -- мок-обработчик подписки
package mock_hand_sub_local

import (
	"crypto/rand"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/kalias"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

type MockHandlerSub struct {
	Msg_   []byte       // Для обратного вызова
	Name_  AHandlerName // Уникальное имя мок-обработчика подписки
	Topic_ ATopic       // Имя топика подписки
	block  sync.RWMutex
}

// NewMockHandlerSub -- возвращает новый обработчик подписки
func NewMockHandlerSub(topic ATopic, localHook string) *MockHandlerSub {
	Hassert(topic != "", "NewMockHandlerSub(): topic is empty")
	Hassert(localHook != "", "NewMockHandlerSub(): name is empty")
	sf := &MockHandlerSub{
		Topic_: topic,
		Name_:  AHandlerName(localHook + "_" + rand.Text()),
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

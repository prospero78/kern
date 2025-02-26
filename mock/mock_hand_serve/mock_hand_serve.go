// package mock_hand_serve -- мок-обработчик входящих запросов
package mock_hand_serve

import (
	"fmt"
	"sync"

	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_alias"
	. "github.com/prospero78/kern/kernel_types"
	"github.com/prospero78/kern/safe_bool"
)

// MockHandlerServe -- мок-обработчик входящих запросов
type MockHandlerServe struct {
	IsBad_ ISafeBool    // Признак сбоя при вызове
	Msg_   []byte       // Для обратного вызова
	Name_  AHandlerName // Уникальное имя мок-обработчика подписки
	Topic_ ATopic       // Имя топика подписки
	block  sync.Mutex
}

// NewMockHandlerServe -- возвращает новый обработчик подписки
func NewMockHandlerServe(topic ATopic, name string) *MockHandlerServe {
	Hassert(topic != "", "NewMockHandlerServe(): topic is empty")
	Hassert(name != "", "NewMockHandlerServe(): name is empty")
	sf := &MockHandlerServe{
		Topic_: topic,
		Name_:  AHandlerName(name),
		IsBad_: safe_bool.NewSafeBool(),
	}
	_ = IBusHandlerServe(sf)
	return sf
}

// Функция обратного вызова подписки
func (sf *MockHandlerServe) FnBack(binMsg []byte) ([]byte, error) {
	sf.block.Lock()
	defer sf.block.Unlock()
	if sf.IsBad_.Get() {
		return nil, fmt.Errorf("FnBack(): isBad==true")
	}
	sf.Msg_ = binMsg
	return sf.Msg_, nil
}

// Возвращает уникальное имя обработчика подписки
func (sf *MockHandlerServe) Topic() ATopic {
	return sf.Topic_
}

// Возвращает топик для обработчика подписки
func (sf *MockHandlerServe) Name() AHandlerName {
	return sf.Name_
}

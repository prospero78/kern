package kbus_local

import (
	"testing"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/mock/mock_hand_serve"
	"gitp78su.ipnodns.ru/svi/kern/mock/mock_hand_sub_local"
)

type tester struct {
	t        *testing.T
	bus      IKernelBus
	handSub  *mock_hand_sub_local.MockHandlerSub
	handServ *mock_hand_serve.MockHandlerServe
}

func TestKernelBusLocal(t *testing.T) {
	sf := &tester{
		t:        t,
		handSub:  mock_hand_sub_local.NewMockHandlerSub("topic_hand_sub", "mock_hand_sub"),
		handServ: mock_hand_serve.NewMockHandlerServe("topic_hand_serv", "mock_hand_serv"),
	}
	sf.new()
}

// Создание локальной шины
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newGood1()

}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("newGood1(): panic=%v", _panic)
		}
	}()
	sf.bus = GetKernelBusLocal()
	sf.bus = GetKernelBusLocal()
	if sf.bus == nil {
		sf.t.Fatalf("newGood1(): IKernelBus==nil")
	}
	if !sf.bus.IsWork() {
		sf.t.Fatalf("newGood1(): bus not work")
	}
}

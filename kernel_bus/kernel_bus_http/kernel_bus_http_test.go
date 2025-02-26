package kernel_bus_http

import (
	"testing"

	. "github.com/prospero78/kern/kernel_types"
	"github.com/prospero78/kern/mock/mock_env"
	"github.com/prospero78/kern/mock/mock_hand_serve"
)

type tester struct {
	t        *testing.T
	bus      IKernelBus
	handServ *mock_hand_serve.MockHandlerServe
}

func TestKernelBusHttp(t *testing.T) {
	sf := &tester{
		t:        t,
		handServ: mock_hand_serve.NewMockHandlerServe("topic_serv", "name_serv"),
	}
	sf.get()
}

// Получает локальную шину
func (sf *tester) get() {
	sf.t.Log("get")
	_ = mock_env.MakeEnv()
	sf.bus = GetKernelBusHttp()
	if sf.bus == nil {
		sf.t.Fatalf("get(): bus==nil")
	}
	_ = GetKernelBusHttp()
}

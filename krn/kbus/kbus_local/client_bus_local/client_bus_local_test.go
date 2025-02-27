package client_bus_local

import (
	"testing"

	. "github.com/prospero78/kern/krn/ktypes"
)

type tester struct {
	t  *testing.T
	cl IBusClient
}

func TestClientBusLocal(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
}

// Создание нового клиента
func (sf *tester) new() {
	sf.t.Log("new")
	sf.cl = NewClientBusLocal()
	if sf.cl == nil {
		sf.t.Fatalf("new(): client==nil")
	}
	err := sf.cl.Publish("local_topic", []byte("test_msg"))
	if err != nil {
		sf.t.Fatalf("new(): err=%v", err)
	}
}

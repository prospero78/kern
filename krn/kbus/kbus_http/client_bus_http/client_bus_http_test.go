package client_bus_http

import (
	"testing"

	"github.com/prospero78/kern/mock/mock_env"
)

type tester struct {
	t *testing.T
}

func TestClientBusHttp(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
}

// Создание нового клиента
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newBad2()
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("newGood1(): panic=%v", _panic)
		}
	}()
	_ = mock_env.MakeEnv()
	_ = NewClientBusHttp("http://localhost:18200/")
}

func (sf *tester) newBad2() {
	sf.t.Log("newBad2")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad2(): panic==nil")
		}
	}()
	_ = mock_env.MakeEnv()
	_ = NewClientBusHttp("")
}

// Нет ничего для HTTP-шины
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic==nil")
		}
	}()
	_ = NewClientBusHttp("url")
}

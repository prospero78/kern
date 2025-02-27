package msg_unsub

import (
	"testing"
)

type tester struct {
	t *testing.T
}

func TestPublishMsg(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.req()
	sf.resp()
}

// Работа с ответом
func (sf *tester) resp() {
	sf.t.Log("resp")
	sf.respBad1()
	resp := &UnsubResp{
		Status_: "test_ok",
	}
	resp.SelfCheck()
}

// Кривые поля ответа
func (sf *tester) respBad1() {
	sf.t.Log("respBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("respBad1(): panic==nil")
		}
	}()
	resp := &UnsubResp{}
	resp.SelfCheck()
}

// Работа с запросом
func (sf *tester) req() {
	sf.t.Log("req")
	sf.reqBad1()
	req := &UnsubReq{
		Name_: "test_topic",
		Uuid_: "test_uuid",
	}
	req.SelfCheck()
}

// Кривые поля
func (sf *tester) reqBad1() {
	sf.t.Log("reqBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("reqBad1(): panic==nil")
		}
	}()
	req := &UnsubReq{}
	req.SelfCheck()
}

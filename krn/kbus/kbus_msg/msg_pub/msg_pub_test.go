package msg_pub

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
	resp := &PublishResp{
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
	resp := &PublishResp{}
	resp.SelfCheck()
}

// Работа с запросом
func (sf *tester) req() {
	sf.t.Log("req")
	sf.reqBad1()
	req := &PublishReq{
		Topic_:  "test_topic",
		Uuid_:   "test_uuid",
		BinMsg_: []byte("test msg"),
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
	req := &PublishReq{}
	req.SelfCheck()
}

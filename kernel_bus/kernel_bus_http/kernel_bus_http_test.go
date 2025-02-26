package kernel_bus_http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/prospero78/kern/kernel_serv_http"
	"github.com/prospero78/kern/mock/mock_env"
	"github.com/prospero78/kern/mock/mock_hand_serve"
)

type tester struct {
	t        *testing.T
	handServ *mock_hand_serve.MockHandlerServe
}

func TestKernelBusHttp(t *testing.T) {
	sf := &tester{
		t:        t,
		handServ: mock_hand_serve.NewMockHandlerServe("topic_serv", "name_serv"),
	}
	sf.get()
	sf.req()
}

// Входящий запрос
func (sf *tester) req() {
	sf.t.Log("req")
	/*
		fibApp := kernel_serv_http.GetKernelServHttp().Fiber()
		req := &ServeReq{
			Topic_:  sf.handServ.Topic_,
			Uuid_:   "test_uud",
			BinReq_: []byte("test_req"),
		}
		var body io.Reader
		hReq := http.NewRequest("POST", "/bus/request", body)
		fibApp.Test(hReq)
		if err != nil {
			sf.t.Fatalf("redirect(): after request, err=%v", err)
		}
		if resp.StatusCode != 303 {
			sf.t.Fatalf("redirect(): statusCode(%v)!=303", resp.StatusCode)
		}
	*/
	sf.reqBad1()
	sf.reqBad2()
	sf.reqBad3()
	sf.reqGood1()
	sf.reqBad4()
	sf.subBad1()
	sf.subBad2()
}

// Проверка кривых полей запроса в процессе подписки
func (sf *tester) subBad2() {
	sf.t.Log("subBad2")
	req := &SubscribeReq{
		Topic_: "",
	}
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("subBad2(): panic==nil")
		}
	}()
	_ = bus.processSubscribe(req)
}

// Кривой запрос
func (sf *tester) subBad1() {
	sf.t.Log("subBad1")
	req := "tra-ta-ta"
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))
	fibApp := kernel_serv_http.GetKernelServHttp().Fiber()
	hReq, err := http.NewRequest("POST", "/bus/sub", body)
	hReq.Header.Add("Content-Type", "application/json")
	if err != nil {
		sf.t.Fatalf("subBad1(): err=%v", err)
	}
	_resp, err := fibApp.Test(hReq)
	if err != nil {
		sf.t.Fatalf("subBad1(): after request, err=%v", err)
	}
	if _resp.StatusCode != 400 {
		sf.t.Fatalf("subBad1(): statusCode(%v)!=400", _resp.StatusCode)
	}
}

// Что-то с обработчиком
func (sf *tester) reqBad4() {
	sf.t.Log("reqBad4")
	sf.handServ.IsBad_.Set()
	defer sf.handServ.IsBad_.Reset()
	req := &ServeReq{
		Topic_:  sf.handServ.Topic_,
		Uuid_:   "test_uuid",
		BinReq_: []byte("test_msg"),
	}
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))
	fibApp := kernel_serv_http.GetKernelServHttp().Fiber()
	hReq, err := http.NewRequest("POST", "/bus/request", body)
	hReq.Header.Add("Content-Type", "application/json")
	if err != nil {
		sf.t.Fatalf("reqBad4(): err=%v", err)
	}
	_resp, err := fibApp.Test(hReq)
	if err != nil {
		sf.t.Fatalf("reqBad4(): after request, err=%v", err)
	}
	if _resp.StatusCode != 200 {
		sf.t.Fatalf("reqBad4(): statusCode(%v)!=200", _resp.StatusCode)
	}
	defer _resp.Body.Close()
	binBody, _ := io.ReadAll(_resp.Body)
	resp := &ServeResp{}
	err = json.Unmarshal(binBody, resp)
	if err != nil {
		sf.t.Fatalf("reqBad4(): err=%v", err)
	}
	if string(resp.Status_) == "ok" {
		sf.t.Fatalf("reqBad4(): status(%v)=='ok'", string(resp.Status_))
	}
}

func (sf *tester) reqGood1() {
	sf.t.Log("reqGood1")
	req := &ServeReq{
		Topic_:  sf.handServ.Topic_,
		Uuid_:   "test_uuid",
		BinReq_: []byte("test_msg"),
	}
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))
	fibApp := kernel_serv_http.GetKernelServHttp().Fiber()
	hReq, err := http.NewRequest("POST", "/bus/request", body)
	hReq.Header.Add("Content-Type", "application/json")
	if err != nil {
		sf.t.Fatalf("reqGood1(): err=%v", err)
	}
	_resp, err := fibApp.Test(hReq)
	if err != nil {
		sf.t.Fatalf("reqGood1(): after request, err=%v", err)
	}
	if _resp.StatusCode != 200 {
		sf.t.Fatalf("reqGood1(): statusCode(%v)!=200", _resp.StatusCode)
	}
	defer _resp.Body.Close()
	binBody, _ := io.ReadAll(_resp.Body)
	resp := &ServeResp{}
	err = json.Unmarshal(binBody, resp)
	if err != nil {
		sf.t.Fatalf("reqGood1(): err=%v", err)
	}
	if string(resp.BinResp_) != "test_msg" {
		sf.t.Fatalf("reqGood1(): resp(%v)!='test_msg'", string(resp.BinResp_))
	}
}

// Нет такого топика для запросов
func (sf *tester) reqBad3() {
	sf.t.Log("reqBad3")
	req := &ServeReq{
		Topic_:  "bad_topic",
		Uuid_:   "test_uuid",
		BinReq_: []byte("test_msg"),
	}
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))
	fibApp := kernel_serv_http.GetKernelServHttp().Fiber()
	hReq, err := http.NewRequest("POST", "/bus/request", body)
	if err != nil {
		sf.t.Fatalf("reqBad3(): err=%v", err)
	}
	resp, err := fibApp.Test(hReq)
	if err != nil {
		sf.t.Fatalf("reqBad3(): after request, err=%v", err)
	}
	if resp.StatusCode != 400 {
		sf.t.Fatalf("reqBad3(): statusCode(%v)!=400", resp.StatusCode)
	}
}

// Нет тела запроса
func (sf *tester) reqBad2() {
	sf.t.Log("reqBad1")
	body := strings.NewReader("test_msg")
	fibApp := kernel_serv_http.GetKernelServHttp().Fiber()
	hReq, err := http.NewRequest("POST", "/bus/request", body)
	if err != nil {
		sf.t.Fatalf("reqBad1(): err=%v", err)
	}
	resp, err := fibApp.Test(hReq)
	if err != nil {
		sf.t.Fatalf("reqBad1(): after request, err=%v", err)
	}
	if resp.StatusCode != 400 {
		sf.t.Fatalf("reqBad1(): statusCode(%v)!=400", resp.StatusCode)
	}
}

// Отсутствуют поля в запросе
func (sf *tester) reqBad1() {
	sf.t.Log("reqBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("reqBad1(): panic==nil")
		}
	}()
	bus.processSendRequest(nil)
}

// Получает локальную шину
func (sf *tester) get() {
	sf.t.Log("get")
	_ = mock_env.MakeEnv()
	_ = GetKernelBusHttp()
	if bus == nil {
		sf.t.Fatalf("get(): bus==nil")
	}
	_ = GetKernelBusHttp()
	bus.RegisterServe(sf.handServ)
}

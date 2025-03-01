package kbus_http

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/krn/kbus/kbus_msg/msg_pub"
	"github.com/prospero78/kern/krn/kbus/kbus_msg/msg_serve"
	"github.com/prospero78/kern/krn/kbus/kbus_msg/msg_sub"
	"github.com/prospero78/kern/krn/kbus/kbus_msg/msg_unsub"
	"github.com/prospero78/kern/krn/kctx"
	"github.com/prospero78/kern/krn/kserv_http"
	"github.com/prospero78/kern/mock/mock_env"
	"github.com/prospero78/kern/mock/mock_hand_serve"
	"github.com/prospero78/kern/mock/mock_hand_sub_local"
)

type tester struct {
	t        *testing.T
	handSub  *mock_hand_sub_local.MockHandlerSub
	handServ *mock_hand_serve.MockHandlerServe
}

func TestKernelBusHttp(t *testing.T) {
	sf := &tester{
		t:        t,
		handSub:  mock_hand_sub_local.NewMockHandlerSub("topic_sub", "http://localhost:18200/bus/pub"),
		handServ: mock_hand_serve.NewMockHandlerServe("topic_serv", "name_serv"),
	}
	ctx := kctx.GetKernelCtx()
	ctx.Set("monolitName", "test_monolit", "comment")
	sf.get()
	sf.req()
	sf.sub()
	sf.pub()
	sf.unsub()
}

// Запрос на отписку
func (sf *tester) unsub() {
	sf.t.Log("unsub")
	sf.unsubBad1()
	sf.unsubBad2()
	sf.unsubGood1()
	sf.unsubBad3()
	sf.unsubGood2()
}

func (sf *tester) unsubGood2() {
	sf.t.Log("unsubGood2")
	err := bus.Subscribe(sf.handSub)
	if err != nil {
		sf.t.Fatalf("unsubGood1(): err=%v", err)
	}
	req := &msg_unsub.UnsubReq{
		Name_: sf.handSub.Name_,
		Uuid_: "test_uuid",
	}
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))
	fibApp := kserv_http.GetKernelServHttp().Fiber()
	hReq, err := http.NewRequest("POST", "/bus/unsub", body)
	if err != nil {
		sf.t.Fatalf("unsubGood2(): err=%v", err)
	}
	hReq.Header.Add("Content-Type", "application/json")

	_resp, err := fibApp.Test(hReq)
	if err != nil {
		sf.t.Fatalf("unsubGood2(): after request, err=%v", err)
	}
	if _resp.StatusCode != 200 {
		sf.t.Fatalf("unsubGood2(): statusCode(%v)!=200", _resp.StatusCode)
	}
	defer _resp.Body.Close()
	binBody, _ := io.ReadAll(_resp.Body)
	resp := &msg_unsub.UnsubResp{}
	err = json.Unmarshal(binBody, resp)
	if err != nil {
		sf.t.Fatalf("unsubGood2(): err=%v", err)
	}
	if string(resp.Status_) != "ok" {
		sf.t.Fatalf("unsubGood2(): resp(%v)!='ok'", string(resp.Status_))
	}
}

// Кривой запрос
func (sf *tester) unsubBad3() {
	sf.t.Log("unsubBad3")
	req := "tra-la-la"
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))
	fibApp := kserv_http.GetKernelServHttp().Fiber()
	hReq, err := http.NewRequest("POST", "/bus/unsub", body)
	hReq.Header.Add("Content-Type", "application/json")
	if err != nil {
		sf.t.Fatalf("unsubBad3(): err=%v", err)
	}
	_resp, err := fibApp.Test(hReq)
	if err != nil {
		sf.t.Fatalf("unsubBad3(): after request, err=%v", err)
	}
	if _resp.StatusCode == 200 {
		sf.t.Fatalf("unsubBad3(): statusCode(%v)==200", _resp.StatusCode)
	}
}

func (sf *tester) unsubGood1() {
	sf.t.Log("unsubGood1")
	err := bus.Subscribe(sf.handSub)
	if err != nil {
		sf.t.Fatalf("unsubGood1(): err=%v", err)
	}
	count := 0
	for count < 100 {
		SleepMs()
		count++
	}
	req := &msg_unsub.UnsubReq{
		Name_: sf.handSub.Name_,
		Uuid_: "test_uuid",
	}
	resp := bus.processUnsubRequest(req)
	if resp.Status_ != "ok" {
		sf.t.Fatalf("unsubGood1(): status(%v)!='ok'", resp.Status_)
	}
}

// Все поля на месте, нет подписчика
func (sf *tester) unsubBad2() {
	sf.t.Log("unsubBad2")
	bus.Unsubscribe(sf.handSub)
	count := 0
	for count < 100 {
		SleepMs()
		count++
	}
	req := &msg_unsub.UnsubReq{
		Name_: sf.handSub.Name_,
		Uuid_: "test_uuid",
	}
	resp := bus.processUnsubRequest(req)
	if resp.Status_ == "ok" {
		sf.t.Fatalf("unsubBad2(): status(%v)=='ok'", resp.Status_)
	}
}

// Нет полей для процесса отписки
func (sf *tester) unsubBad1() {
	sf.t.Log("unsubBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("unsubBad1(): panic==nil")
		}
	}()
	req := &msg_unsub.UnsubReq{}
	_ = bus.processUnsubRequest(req)
}

// Запрос на публикацию
func (sf *tester) pub() {
	sf.t.Log("pub")
	sf.pubBad1()
	sf.pubGood1()
	sf.pubBad2()
	sf.pubBad3()
	sf.pubGood2()
}

func (sf *tester) pubGood2() {
	sf.t.Log("pubGood2")
	req := &msg_pub.PublishReq{
		Topic_:  "topic_sub",
		Uuid_:   "test_uuid",
		BinMsg_: []byte("http_pub"),
	}
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))
	fibApp := kserv_http.GetKernelServHttp().Fiber()
	hReq, err := http.NewRequest("POST", "/bus/pub", body)
	hReq.Header.Add("Content-Type", "application/json")
	if err != nil {
		sf.t.Fatalf("pubGood2(): err=%v", err)
	}
	_resp, err := fibApp.Test(hReq)
	if err != nil {
		sf.t.Fatalf("pubGood2(): after request, err=%v", err)
	}
	if _resp.StatusCode != 200 {
		sf.t.Fatalf("pubGood2(): statusCode(%v)!=200", _resp.StatusCode)
	}
	defer _resp.Body.Close()
	binBody, _ := io.ReadAll(_resp.Body)
	resp := &msg_pub.PublishResp{}
	err = json.Unmarshal(binBody, resp)
	if err != nil {
		sf.t.Fatalf("pubGood2(): err=%v", err)
	}
	if string(resp.Status_) != "ok" {
		sf.t.Fatalf("pubGood2(): resp(%v)!='ok'", string(resp.Status_))
	}
}

// Кривой запрос
func (sf *tester) pubBad3() {
	sf.t.Log("pubBad3")
	req := "tra-la-la"
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))
	fibApp := kserv_http.GetKernelServHttp().Fiber()
	hReq, err := http.NewRequest("POST", "/bus/pub", body)
	hReq.Header.Add("Content-Type", "application/json")
	if err != nil {
		sf.t.Fatalf("pubBad3(): err=%v", err)
	}
	_resp, err := fibApp.Test(hReq)
	if err != nil {
		sf.t.Fatalf("pubBad3(): after request, err=%v", err)
	}
	if _resp.StatusCode == 200 {
		sf.t.Fatalf("pubBad3(): statusCode(%v)==200", _resp.StatusCode)
	}
}

// Что-то случилось внутри шины
func (sf *tester) pubBad2() {
	sf.t.Log("pubBad2")
	bus.IsWork_.Reset()
	defer bus.IsWork_.Set()

	req := &msg_pub.PublishReq{
		Topic_:  "topic_sub",
		Uuid_:   "test_uuid",
		BinMsg_: []byte("test_pub"),
	}
	resp := bus.processPublish(req)
	if resp.Status_ == "ok" {
		sf.t.Fatalf("pubBad2(): status(%v)=='ok'", resp.Status_)
	}
}

// Все поля на месте
func (sf *tester) pubGood1() {
	sf.t.Log("pubGood1")
	err := bus.Subscribe(sf.handSub)
	if err != nil {
		sf.t.Fatalf("pubGood1(): err=%v", err)
	}
	req := &msg_pub.PublishReq{
		Topic_:  "topic_sub",
		Uuid_:   "test_uuid",
		BinMsg_: []byte("test_pub"),
	}
	_ = bus.processPublish(req)
	for {
		SleepMs()
		msg := string(sf.handSub.Msg())
		if msg != "" {
			break
		}
	}
	msg := string(sf.handSub.Msg())
	if msg != "test_pub" {
		sf.t.Fatalf("pubGood1(): msg(%v)!='test_pub'", msg)
	}
}

// Нет полей для процесса публикации
func (sf *tester) pubBad1() {
	sf.t.Log("pubBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("pubBad1(): panic==nil")
		}
	}()
	req := &msg_pub.PublishReq{}
	_ = bus.processPublish(req)
}

// Запрос подписки на топик
func (sf *tester) sub() {
	sf.t.Log("sub")
	sf.subBad1()
	sf.subBad2()
	sf.subGood1()
	sf.subBad3()
	sf.subGood2()
}

// Полный процесс подписки
func (sf *tester) subGood2() {
	sf.t.Log("subGood2")
	req := &msg_sub.SubscribeReq{
		Topic_:   "topic_serv",
		Uuid_:    "test_uuid",
		WebHook_: "http://localhost:18200/bus/pub/",
	}
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))
	fibApp := kserv_http.GetKernelServHttp().Fiber()
	hReq, err := http.NewRequest("POST", "/bus/sub", body)
	hReq.Header.Add("Content-Type", "application/json")
	if err != nil {
		sf.t.Fatalf("subBad1(): err=%v", err)
	}
	_resp, err := fibApp.Test(hReq)
	if err != nil {
		sf.t.Fatalf("subBad1(): after request, err=%v", err)
	}
	if _resp.StatusCode != 200 {
		sf.t.Fatalf("subBad1(): statusCode(%v)!=200", _resp.StatusCode)
	}
	defer _resp.Body.Close()
	binBody, _ := io.ReadAll(_resp.Body)
	resp := &msg_sub.SubscribeResp{}
	err = json.Unmarshal(binBody, resp)
	if err != nil {
		sf.t.Fatalf("subBad1(): err=%v", err)
	}
	if string(resp.Status_) != "ok" {
		sf.t.Fatalf("subBad1(): resp(%v)!='ok'", string(resp.Status_))
	}
}

// Отключена базовая шина
func (sf *tester) subBad3() {
	sf.t.Log("subBad3")
	req := &msg_sub.SubscribeReq{
		Topic_:   "topic_serv",
		Uuid_:    "test_uuid",
		WebHook_: "http://localhost:18200/bus/pub/",
	}
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("subBad3(): panic!=nil")
		}
	}()
	// _bus := kernel_bus_base.GetKernelBusBase()
	bus.IsWork_.Reset()
	defer bus.IsWork_.Set()
	resp := bus.processSubscribe(req)
	if resp.Status_ == "ok" {
		sf.t.Fatalf("subBad3(): resp==ok")
	}
}

// Проверка полей запроса в процессе подписки
func (sf *tester) subGood1() {
	sf.t.Log("subGood1")
	req := &msg_sub.SubscribeReq{
		Topic_:   "topic_serv",
		Uuid_:    "test_uuid",
		WebHook_: "http://localhost:18200/bus/",
	}
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("subGood1(): panic!=nil")
		}
	}()
	_ = bus.processSubscribe(req)
}

// Проверка кривых полей запроса в процессе подписки
func (sf *tester) subBad2() {
	sf.t.Log("subBad2")
	req := &msg_sub.SubscribeReq{
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
	fibApp := kserv_http.GetKernelServHttp().Fiber()
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

// Входящий запрос
func (sf *tester) req() {
	sf.t.Log("req")
	sf.reqBad1()
	sf.reqBad2()
	sf.reqBad3()
	sf.reqGood1()
	sf.reqBad4()
}

// Что-то с обработчиком
func (sf *tester) reqBad4() {
	sf.t.Log("reqBad4")
	sf.handServ.IsBad_.Set()
	defer sf.handServ.IsBad_.Reset()
	req := &msg_serve.ServeReq{
		Topic_:  sf.handServ.Topic_,
		Uuid_:   "test_uuid",
		BinReq_: []byte("test_msg"),
	}
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))
	fibApp := kserv_http.GetKernelServHttp().Fiber()
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
	resp := &msg_serve.ServeResp{}
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
	req := &msg_serve.ServeReq{
		Topic_:  sf.handServ.Topic_,
		Uuid_:   "test_uuid",
		BinReq_: []byte("test_msg"),
	}
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))
	fibApp := kserv_http.GetKernelServHttp().Fiber()
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
	resp := &msg_serve.ServeResp{}
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
	req := &msg_serve.ServeReq{
		Topic_:  "bad_topic",
		Uuid_:   "test_uuid",
		BinReq_: []byte("test_msg"),
	}
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))
	fibApp := kserv_http.GetKernelServHttp().Fiber()
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
	fibApp := kserv_http.GetKernelServHttp().Fiber()
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
	_ = os.Unsetenv("LOCAL_HTTP_URL")
	os.Setenv("LOCAL_HTTP_URL", "http://localhost:18312/")
	_ = GetKernelBusHttp()
	if bus == nil {
		sf.t.Fatalf("get(): bus==nil")
	}
	_ = GetKernelBusHttp()
	bus.RegisterServe(sf.handServ)
}

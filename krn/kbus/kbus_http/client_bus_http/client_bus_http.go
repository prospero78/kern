// package client_bus_http -- клиент HTTP-шины
package client_bus_http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/local_ctx"
	"github.com/prospero78/kern/kc/safe_bool"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kbus/kbus_http"
	"github.com/prospero78/kern/krn/kbus/kbus_msg/msg_pub"
	"github.com/prospero78/kern/krn/kbus/kbus_msg/msg_serve"
	"github.com/prospero78/kern/krn/kbus/kbus_msg/msg_sub"
	"github.com/prospero78/kern/krn/kbus/kbus_msg/msg_unsub"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
	"github.com/prospero78/kern/mock/mock_hand_sub_http"
)

// ClientBusHttp -- клиент HTTP-шины
type ClientBusHttp struct {
	bus       IKernelBus
	ctx       ILocalCtx
	log       ILogBuf
	isWork    ISafeBool
	urlRemote string // URL дистанционной шины
	urlLocal  string // URL локальной шины
}

// NewClientBusHttp - -возвращает новый клиент HTTP-шины
func NewClientBusHttp(url string) IBusClient {
	Hassert(url != "", "NewClientBusHttp(): url is empty")
	kCtx := kctx.GetKernelCtx()
	urlLocal := os.Getenv("LOCAL_HTTP_URL")
	Hassert(urlLocal != "", "NewClientBusHttp(): env LOCAL_HTTP_URL not set")
	sf := &ClientBusHttp{
		ctx:       local_ctx.NewLocalCtx(kCtx.BaseCtx()),
		bus:       kbus_http.GetKernelBusHttp(),
		isWork:    safe_bool.NewSafeBool(),
		urlRemote: strings.TrimSuffix(url, "/"),
		urlLocal:  strings.TrimSuffix(urlLocal, "/"),
	}
	sf.log = sf.ctx.Log()
	return sf
}

// Unsubscribe -- отписывается от топика в дистанционной шине
func (sf *ClientBusHttp) Unsubscribe(handler IBusHandlerSubscribe) {
	_uuid, err := uuid.NewV6()
	Hassert(err == nil, "ClientBusHttp.Unsubscribe(): in generate UUID v6, err=\n\t%v", err)

	req := &msg_unsub.UnsubReq{
		Name_: handler.Name(),
		Uuid_: _uuid.String(),
	}
	req.SelfCheck()
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))

	hReq, err := http.NewRequest("POST", sf.urlRemote+"/bus/unsub", body)
	Hassert(err == nil, "ClientBusHttp.Unsubscribe(): in new request, err=\n\t%v")

	binBody, err := sf.makePost(hReq)
	if err != nil {
		sf.log.Err("ClientBusHttp.Unsubscribe(): in make request, err=\n\t%v")
		return
	}
	resp := &msg_unsub.UnsubResp{}
	err = json.Unmarshal(binBody, resp)
	Hassert(err == nil, "ClientBusHttp.Unsubscribe(): in unmarshal response,  err=\n\t%v", err)
	if string(resp.Status_) != "ok" {
		sf.log.Err("ClientBusHttp.Unsubscribe(): resp!='ok', err=\n\t%v", resp.Status_)
	}
	Hassert(resp.Uuid_ == req.Uuid_, "ClientBusHttp.Unsubscribe(): resp uuid(%v) bad", resp.Uuid_)
}

// Subscribe -- подписывается на топик в дистанционной шине
func (sf *ClientBusHttp) Subscribe(handler IBusHandlerSubscribe) error {
	_uuid, err := uuid.NewV6()
	Hassert(err == nil, "ClientBusHttp.Subscribe(): in generate UUID v6, err=\n\t%v", err)
	req := &msg_sub.SubscribeReq{
		Topic_:   handler.Topic(),
		Uuid_:    _uuid.String(),
		WebHook_: sf.urlLocal + "/bus/pub",
	}
	req.SelfCheck()
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))

	hReq, err := http.NewRequest("POST", sf.urlRemote+"/bus/sub", body)
	Hassert(err == nil, "ClientBusHttp.Subscribe(): in new request, err=\n\t%v")

	binBody, err := sf.makePost(hReq)
	if err != nil {
		err := fmt.Errorf("ClientBusHttp.Subscribe(): in make request, err=\n\t%w", err)
		sf.log.Err(err.Error())
		return err
	}
	resp := &msg_sub.SubscribeResp{}
	err = json.Unmarshal(binBody, resp)
	Hassert(err == nil, "ClientBusHttp.Subscribe(): in unmarshal response,  err=\n\t%v", err)
	if string(resp.Status_) != "ok" {
		err := fmt.Errorf("ClientBusHttp.Subscribe(): resp!='ok', err=\n\t%v", resp.Status_)
		sf.log.Err(err.Error())
		return err
	}
	Hassert(resp.Uuid_ == req.Uuid_, "ClientBusHttp.Subscribe(): resp uuid(%v) bad", resp.Uuid_)
	// FIXME: вот тут похоже дичь
	_handler := handler.(*mock_hand_sub_http.MockHandSubHttp)
	_handler.SetName(resp.Name_)
	err = sf.bus.Subscribe(_handler)
	return err
}

// SendRequest -- отправляет в дистанционную шину запрос
func (sf *ClientBusHttp) SendRequest(topic ATopic, binReq []byte) ([]byte, error) {
	_uuid, err := uuid.NewV6()
	Hassert(err == nil, "ClientBusHttp.SendRequest(): in generate UUID v6, err=\n\t%v", err)
	req := &msg_serve.ServeReq{
		Topic_:  topic,
		Uuid_:   _uuid.String(),
		BinReq_: binReq,
	}
	req.SelfCheck()
	_binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(_binReq))

	hReq, err := http.NewRequest("POST", sf.urlRemote+"/bus/request", body)
	Hassert(err == nil, "ClientBusHttp.SendRequest(): in new request, err=\n\t%v")

	binBody, err := sf.makePost(hReq)
	if err != nil {
		err := fmt.Errorf("ClientBusHttp.SendRequest(): in make request, err=\n\t%w", err)
		sf.log.Err(err.Error())
		return nil, err
	}
	resp := &msg_serve.ServeResp{}
	err = json.Unmarshal(binBody, resp)
	Hassert(err == nil, "ClientBusHttp.SendRequest(): in unmarshal response,  err=\n\t%v", err)
	if string(resp.Status_) != "ok" {
		err := fmt.Errorf("ClientBusHttp.SendRequest(): resp!='ok', err=\n\t%v", resp.Status_)
		sf.log.Err(err.Error())
		return nil, err
	}
	Hassert(resp.Uuid_ == req.Uuid_, "ClientBusHttp.SendRequest(): resp uuid(%v) bad", resp.Uuid_)
	return resp.BinResp_, nil
}

// RegisterServe -- регистрирует в локальной шине обработчик
func (sf *ClientBusHttp) RegisterServe(handler IBusHandlerServe) {
	Hassert(handler != nil, "ClientBusHttp.RegisterServe(): handler==nil")
	sf.bus.RegisterServe(handler)
}

// Publish -- публикует сообщение в дистанционной шину
func (sf *ClientBusHttp) Publish(topic ATopic, binMsg []byte) error {
	_uuid, err := uuid.NewV6()
	Hassert(err == nil, "ClientBusHttp.Publish(): in generate UUID v6, err=\n\t%v", err)
	req := &msg_pub.PublishReq{
		Topic_:  topic,
		Uuid_:   _uuid.String(),
		BinMsg_: binMsg,
	}
	req.SelfCheck()
	binReq, _ := json.MarshalIndent(req, "", "  ")
	body := strings.NewReader(string(binReq))

	hReq, err := http.NewRequest("POST", sf.urlRemote+"/bus/pub", body)
	Hassert(err == nil, "ClientBusHttp.Publish(): in new request, err=\n\t%v")

	binBody, err := sf.makePost(hReq)
	if err != nil {
		err := fmt.Errorf("ClientBusHttp.Publish(): in make request, err=\n\t%w", err)
		sf.log.Err(err.Error())
		return err
	}
	resp := &msg_pub.PublishResp{}
	err = json.Unmarshal(binBody, resp)
	Hassert(err == nil, "ClientBusHttp.Publish(): in unmarshal response,  err=\n\t%v", err)
	if string(resp.Status_) != "ok" {
		err := fmt.Errorf("ClientBusHttp.Publish(): resp!='ok', err=\n\t%v", resp.Status_)
		sf.log.Err(err.Error())
		return err
	}
	Hassert(resp.Uuid_ == req.Uuid_, "ClientBusHttp.Publish(): resp uuid(%v) bad", resp.Uuid_)
	return nil
}

// Единый обработчик запросов
func (sf *ClientBusHttp) makePost(hReq *http.Request) ([]byte, error) {
	hReq.Header.Add("Content-Type", "application/json")
	_resp, err := http.DefaultClient.Do(hReq)
	if err != nil {
		err := fmt.Errorf("ClientBusHttp.makePost(): after request, err=\n\t%w", err)
		sf.log.Err(err.Error())
		return nil, err
	}
	defer _resp.Body.Close()
	binBody, _ := io.ReadAll(_resp.Body)
	return binBody, nil
}

// Log -- возвращает локальный лог клиента
func (sf *ClientBusHttp) Log() ILogBuf {
	return sf.log
}

// IsWork -- возвращает признак работы
func (sf *ClientBusHttp) IsWork() bool {
	return sf.bus.IsWork()
}

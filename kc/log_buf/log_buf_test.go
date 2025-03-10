package log_buf

import (
	"testing"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

type tester struct {
	t   *testing.T
	log ILogBuf
}

func TestLogBuf(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
}

// Создаёт новый буферный лог
func (sf *tester) new() {
	sf.t.Log("new")
	sf.log = NewLogBuf()
	if sf.log == nil {
		sf.t.Fatalf("new(): log==nil")
	}
	msg := sf.log.Get(-1)
	_ = sf.log.GetErr(-1)
	if msg == nil {
		sf.t.Fatalf("new(): msg==nil")
	}
	sf.log.Debug("test msg: %v", 45)
	sf.log.Info("test msg: %v", 46)
	sf.log.Warn("test msg: %v", 47)
	for i := range 120 {
		sf.log.Err("test err: %v", i)
	}
	_ = sf.log.Get(120)
	_ = sf.log.GetErr(120)

	_ = sf.log.Get(-1)
	_ = sf.log.GetErr(-1)

	_ = sf.log.Get(19)
	_ = sf.log.GetErr(20)
	_ = sf.log.Size()
}

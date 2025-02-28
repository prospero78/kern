package main

import (
	"testing"
	"time"

	"github.com/prospero78/kern/krn/kctx"
	"github.com/prospero78/kern/mock/mock_env"
)

func TestMain(t *testing.T) {
	_ = mock_env.MakeEnv()
	go main()
	time.Sleep(time.Millisecond * 250)
	kCtx := kctx.GetKernelCtx()
	kCtx.Cancel()
	kCtx.Wg().Wait()
}

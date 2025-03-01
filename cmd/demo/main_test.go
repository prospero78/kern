package main

import (
	"os"
	"testing"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/krn/kctx"
	"github.com/prospero78/kern/mock/mock_env"
)

func TestMain(t *testing.T) {
	_ = mock_env.MakeEnv()
	_ = os.Unsetenv("LOCAL_HTTP_URL")
	os.Setenv("LOCAL_HTTP_URL", "http://localhost:18330/")
	go main()
	kCtx := kctx.GetKernelCtx()
	for {
		SleepMs()
		if kCtx.Get("monolitName") != nil {
			break
		}
	}

	kCtx.Cancel()
	kCtx.Wg().Wait()
}

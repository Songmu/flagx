package flagx

import (
	"flag"
	"runtime"
	"testing"
	"time"
	"weak"
)

func TestRegistrationLifetime(t *testing.T) {
	key := func() weak.Pointer[flag.FlagSet] {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		var value string
		StringVar(fs, &value, "output|o", "", "an option")
		return weak.Make(fs)
	}()
	deadline := time.Now().Add(10 * time.Second)
	for {
		runtime.GC()
		nameOrders.Lock()
		_, retained := nameOrders.sets[key]
		nameOrders.Unlock()
		if !retained {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("registration metadata retained a discarded FlagSet")
		}
		time.Sleep(time.Millisecond)
	}
}

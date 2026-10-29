package flagx

import (
	"flag"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"weak"
)

// Keep only names, not Values or FlagSets, so registrations cannot retain a set
// through a custom Value that refers back to it.
var nameOrders = struct {
	sync.Mutex
	sets map[weak.Pointer[flag.FlagSet]][][]string
}{sets: make(map[weak.Pointer[flag.FlagSet]][][]string)}

func registerNames(fs *flag.FlagSet, names string, register func(string)) {
	list := strings.Split(names, "|")
	for _, name := range list {
		if name == "" {
			panic("flagx: empty flag name in " + strconv.Quote(names))
		}
	}
	for _, name := range list {
		register(name)
	}
	rememberNames(fs, list)
}

func rememberNames(fs *flag.FlagSet, names []string) {
	defer runtime.KeepAlive(fs)
	nameOrders.Lock()
	defer nameOrders.Unlock()
	key := weak.Make(fs)
	if _, ok := nameOrders.sets[key]; !ok {
		runtime.AddCleanup(fs, forgetNames, key)
	}
	nameOrders.sets[key] = append(nameOrders.sets[key], names)
}

func forgetNames(key weak.Pointer[flag.FlagSet]) {
	nameOrders.Lock()
	defer nameOrders.Unlock()
	delete(nameOrders.sets, key)
}

type nameOrder struct {
	group int
	rank  int
}

func registeredNames(fs *flag.FlagSet) map[string]nameOrder {
	defer runtime.KeepAlive(fs)
	nameOrders.Lock()
	defer nameOrders.Unlock()
	result := make(map[string]nameOrder)
	for group, names := range nameOrders.sets[weak.Make(fs)] {
		for _, name := range names {
			result[name] = nameOrder{group: group, rank: len(result)}
		}
	}
	return result
}

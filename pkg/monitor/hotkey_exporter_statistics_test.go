package monitor

import (
	"fmt"
	"math/rand"
	"testing"
)

func mockHotkeyPair(n int) []HotKeyPair {
	pairs := make([]HotKeyPair, n)
	for i := 0; i < n; i++ {
		pairs[i].key = fmt.Sprintf("mock-key-%d", rand.Intn(n))
		pairs[i].count = uint64(rand.Intn(n) + rand.Intn(n)*2)
	}
	return pairs
}

func checkFilter(paris []HotKeyPair) bool {
	aux := make(map[string]bool)
	for _, hotkey := range paris {

		if _, ok := aux[hotkey.key]; ok {
			return false
		}

		aux[hotkey.key] = true
	}

	return true
}

func TestHotkeyStatistics_Filter(t *testing.T) {
	testRound := 1 << 4
	dataCount := 1 << 16
	for i := 0; i < testRound; i++ {
		t.Run(fmt.Sprintf("parallel round %d", i+1), func(t *testing.T) {
			statistics := new(HotKeyStatistics)
			statistics.hotkeys = mockHotkeyPair(dataCount)
			statistics.Filter()
			if pass := checkFilter(statistics.hotkeys); !pass {
				t.Error("Filter test failed")
			}
		})
	}
}

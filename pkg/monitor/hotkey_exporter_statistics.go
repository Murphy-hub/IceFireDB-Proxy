package monitor

import (
	"sort"
)

type HotKeyPair struct {
	key               string
	count             uint64
	countAvgPerSecond float64
}

type HotKeyStatistics struct {
	count   int
	total   int
	hotkeys []HotKeyPair
}

func (h *HotKeyStatistics) Filter() {
	auxiliary := make(map[string]HotKeyPair)
	for _, hotkeyPair := range h.hotkeys {
		if v, ok := auxiliary[hotkeyPair.key]; !ok {
			auxiliary[hotkeyPair.key] = hotkeyPair
		} else {
			if hotkeyPair.count > v.count {
				auxiliary[hotkeyPair.key] = hotkeyPair
			}
		}
	}

	newHotkeys := make([]HotKeyPair, 0, len(auxiliary))
	for _, hotkey := range auxiliary {
		newHotkeys = append(newHotkeys, hotkey)
	}

	h.hotkeys = nil
	h.hotkeys = newHotkeys
}

func (h *HotKeyStatistics) Init(data *HotKeyMonitorDataS, threshold int) {
	h.count = len(data.HotKeyData)
	h.hotkeys = make([]HotKeyPair, 0, h.count)

	rangeSec := data.TimeRange.End.Sub(data.TimeRange.Start).Seconds()
	if rangeSec == 0 {
		rangeSec = 1.0
	}

	for hotkey, count := range data.HotKeyData {

		h.hotkeys = append(h.hotkeys, HotKeyPair{
			key:               hotkey,
			count:             count,
			countAvgPerSecond: float64(count) / rangeSec,
		})

		h.total += int(count)
	}

	h.Filter()

	if threshold > 0 && len(h.hotkeys) > threshold {
		sort.Sort(h)
		h.hotkeys = h.hotkeys[:threshold]
	}
}

func (h *HotKeyStatistics) Len() int {
	return len(h.hotkeys)
}

func (h *HotKeyStatistics) Swap(i, j int) {
	h.hotkeys[i], h.hotkeys[j] = h.hotkeys[j], h.hotkeys[i]
}

func (h *HotKeyStatistics) Less(i, j int) bool {
	return h.hotkeys[i].count > h.hotkeys[j].count
}

func (h *HotKeyStatistics) GetHotKeyPairArray() []HotKeyPair {
	return h.hotkeys
}

func (h *HotKeyStatistics) GetHotKeyTotal() int {
	return h.count
}

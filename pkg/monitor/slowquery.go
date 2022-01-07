package monitor

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/IceFireDB/IceFireDB-Proxy/utils"
)

type SlowQueryConfS struct {
	Enable                 bool
	SlowQueryTimeThreshold int
	MaxListSize            int
	sync.RWMutex
}

type SlowQueryDataS struct {
	Resp      []interface{}
	StartTime time.Time
	EndTime   time.Time
}

type SlowQueryMonitorDataS struct {
	SlowQueryDataList []*SlowQueryDataS
	Index             int32
	sync.RWMutex
}

func (m *Monitor) IsSlowQuery(args []interface{}, startTime time.Time, endTime time.Time) {
	subMilliseconds := endTime.Sub(startTime).Milliseconds()
	if subMilliseconds < int64(m.SlowQueryConf.SlowQueryTimeThreshold) {
		return
	}

	m.SlowQueryMonitorData.RLock()
	currentIndex := m.SlowQueryMonitorData.Index
	m.SlowQueryMonitorData.RUnlock()

	if currentIndex >= int32(m.SlowQueryConf.MaxListSize-1) {
		return
	}

	m.SlowQueryMonitorData.Lock()

	index := m.SlowQueryMonitorData.Index + 1

	if index >= int32(m.SlowQueryConf.MaxListSize) {
		m.SlowQueryMonitorData.Unlock()
		return
	}

	m.SlowQueryMonitorData.SlowQueryDataList[index].Resp = args
	m.SlowQueryMonitorData.SlowQueryDataList[index].StartTime = startTime
	m.SlowQueryMonitorData.SlowQueryDataList[index].EndTime = endTime
	m.SlowQueryMonitorData.Index = index
	m.SlowQueryMonitorData.Unlock()

	respStr := ""
	for _, v := range args {
		respStr += utils.GetInterfaceString(v) + " "
	}

	/*commandLen := len(args)
	for i := 0; i < commandLen; i++ {
		respStr += string(resp.Array[i].Value) + " "
	}*/
	logrus.Warnf("Found slowquery: %s, cost : %d ms.", respStr, endTime.Sub(startTime).Milliseconds())
}

func (m *Monitor) GetSlowQueryData() (data []*SlowQueryDataS, count int) {
	m.SlowQueryMonitorData.Lock()
	defer m.SlowQueryMonitorData.Unlock()

	count = int(m.SlowQueryMonitorData.Index) + 1

	data = make([]*SlowQueryDataS, count)

	for i := 0; i < count; i++ {
		data[i] = m.SlowQueryMonitorData.SlowQueryDataList[i]

		m.SlowQueryMonitorData.SlowQueryDataList[i] = &SlowQueryDataS{}
	}

	m.SlowQueryMonitorData.Index = INIT_INDEX

	return
}

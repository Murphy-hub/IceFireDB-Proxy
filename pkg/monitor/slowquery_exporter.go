package monitor

import (
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	slowQueryExecTimeLabels = append(BasicLabels, "_01_command", "_02_exec_time", "_03_start_time")

	slowQueryTotalDescKey          = "slowquery.total"
	slowQueryExecTimeRecordDescKey = "slowquery.exec_time_record"

	slowQueryGaugeVec = map[string]*prometheus.Desc{
		slowQueryTotalDescKey:          NewDesc("slowquery_total", "Count of slow query", BasicLabels),
		slowQueryExecTimeRecordDescKey: NewDesc("slowquery_exec_time_record", "Record of slow query command execute time", slowQueryExecTimeLabels),
	}
)

type SlowQueryExporter struct {
	mon              *Monitor
	basicConf        *ExporterConf
	slowQueryMetrics map[string]*prometheus.Desc
}

func NewSlowQueryExporter(mon *Monitor, c *ExporterConf) *SlowQueryExporter {
	metrics := make(map[string]*prometheus.Desc)
	metrics[slowQueryTotalDescKey] = slowQueryGaugeVec[slowQueryTotalDescKey]
	metrics[slowQueryExecTimeRecordDescKey] = slowQueryGaugeVec[slowQueryExecTimeRecordDescKey]
	return &SlowQueryExporter{
		mon:              mon,
		basicConf:        c,
		slowQueryMetrics: metrics,
	}
}

func (s *SlowQueryExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, slowQueryMetric := range s.slowQueryMetrics {
		ch <- slowQueryMetric
	}
}

func (s *SlowQueryExporter) Collect(ch chan<- prometheus.Metric) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("slowQuery promethues collect panic:", r)
		}
	}()
	data, count := s.mon.GetSlowQueryData()
	if len(data) == 0 {

		ch <- prometheus.MustNewConstMetric(
			s.slowQueryMetrics[slowQueryTotalDescKey],
			prometheus.GaugeValue,
			0,
			s.basicConf.Host,
		)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		s.slowQueryMetrics[slowQueryTotalDescKey],
		prometheus.GaugeValue,
		float64(count),
		s.basicConf.Host,
	)

	statistics := new(SlowQueryStatistics)
	statistics.Init(data, s.basicConf.SlowQueryConf.RecordLimit)

	slowQueryPairArray := statistics.GetQueryPairArray()
	for _, q := range slowQueryPairArray {
		ch <- prometheus.MustNewConstMetric(
			s.slowQueryMetrics[slowQueryExecTimeRecordDescKey],
			prometheus.GaugeValue,

			0,
			s.basicConf.Host,
			q.formatRespCommand,
			strconv.FormatInt(q.execTime, 10),
			q.startTime.Format("2006-01-02 15:04:05"),
		)
	}
}

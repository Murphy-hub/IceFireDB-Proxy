package monitor

import (
	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	tcpStateDescKey = "tcp.state"
	networkGaugeVec = map[string]*prometheus.Desc{
		tcpStateDescKey: NewDesc("tcp_state_count", "Count of tcp state", append(BasicLabels, "tcp_state")),
	}
)

func NewNetExporter(c *ExporterConf) *NetworkExporter {
	cf := prometheus.NewCounterFunc(prometheus.CounterOpts{
		Namespace:   Namespace,
		Name:        "net_input_bytes_total",
		Help:        "net_input_bytes_total metric",
		ConstLabels: BasicLabelsMap,
	}, CurrentNetworkStatInputByte)

	cfw := prometheus.NewCounterFunc(prometheus.CounterOpts{
		Namespace:   Namespace,
		Name:        "net_output_bytes_total",
		Help:        "net_output_bytes_total metric",
		ConstLabels: BasicLabelsMap,
	}, CurrentNetworkStatOutputByte)
	prometheus.MustRegister(cf, cfw)

	metrics := make(map[string]*prometheus.Desc)
	metrics[tcpStateDescKey] = networkGaugeVec[tcpStateDescKey]
	return &NetworkExporter{
		networkMetrics: metrics,
		basicConf:      c,
	}
}

type NetworkExporter struct {
	basicConf      *ExporterConf
	networkMetrics map[string]*prometheus.Desc
}

func (s *NetworkExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, nm := range s.networkMetrics {
		ch <- nm
	}
}

func (s *NetworkExporter) Collect(ch chan<- prometheus.Metric) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("network promethues collect panic", r)
		}
	}()
	tcp, err := GetNetstat()
	if err != nil {
		return
	}
	for k, v := range tcp {
		ch <- prometheus.MustNewConstMetric(
			s.networkMetrics[tcpStateDescKey],
			prometheus.GaugeValue,
			float64(v),
			s.basicConf.Host,
			k,
		)
	}
}

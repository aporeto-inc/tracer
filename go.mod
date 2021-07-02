module github.com/aporeto-inc/tracer

go 1.15

require (
	github.com/google/go-querystring v1.0.0
	github.com/gorilla/websocket v1.4.2
	github.com/grafana/loki v1.6.1
	github.com/json-iterator/go v1.1.10
	github.com/olekukonko/tablewriter v0.0.4
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/common v0.10.0
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/spf13/pflag v1.0.5
	go.aporeto.io/addedeffect v1.77.1-0.20210630212802-34063134b871
	go.aporeto.io/elemental v1.100.1-0.20210630205640-d4330d5ba464
	go.aporeto.io/gaia v1.94.1-0.20210701231023-b77f8658f7cf
	go.aporeto.io/tg v1.34.1-0.20210528201128-159c302ba155
	go.aporeto.io/underwater v1.170.0
	go.uber.org/zap v1.16.0
)

replace github.com/hpcloud/tail => github.com/grafana/tail v0.0.0-20191024143944-0b54ddf21fe7

replace k8s.io/client-go => k8s.io/client-go v0.18.3

// >v1.2.0 has some conflict with prometheus/alertmanager. Hence prevent the upgrade till it's fixed.
replace github.com/satori/go.uuid => github.com/satori/go.uuid v1.2.0

// Use fork of gocql that has gokit logs and Prometheus metrics.
replace github.com/gocql/gocql => github.com/grafana/gocql v0.0.0-20200605141915-ba5dc39ece85

// Same as Cortex, we can't upgrade to grpc 1.30.0 until go.etcd.io/etcd will support it.
replace google.golang.org/grpc => google.golang.org/grpc v1.29.1

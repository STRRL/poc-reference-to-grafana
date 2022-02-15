package grafana

import (
	"context"
	"fmt"
	"time"
)

func Example_datasourceProxyDatasourceQuery() {
	proxy := NewDatasourceProxy(
		"http://localhost:13000",
		ViewerAPIKey,
	)
	series, err := proxy.DatasourceQuery(
		context.TODO(),
		2,
		"sum by (instance)(rate(node_cpu_seconds_total{mode=\"system\",instance=\"node-exporter:9100\",job=\"node-exporter\"}[60s])) * 100",
		time.Now().Add(-time.Hour),
		time.Now(),
		"30",
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", series)
	// Output:
}

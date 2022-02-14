package main

import (
	"github.com/STRRL/poc-reference-to-grafana/pkg/grafana"
	"github.com/STRRL/poc-reference-to-grafana/pkg/repository"
)

func main() {
	grafana.FetchDataFromGrafana()
	_, _ = repository.ConnectToMysql()
}

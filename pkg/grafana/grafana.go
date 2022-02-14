package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/STRRL/poc-reference-to-grafana/pkg/entity"
	"github.com/grafana/grafana-api-golang-client"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"log"
)

const ViewerAPIKey = "eyJrIjoibzBXVzFsenN1S1d1Nnh5amdnUVNTT1hSNTlydFpmaFciLCJuIjoiQUQiLCJpZCI6MX0="

const ExamplePanelRawURL = "http://localhost:13000/d/rYdddlPWk/node-exporter-full?orgId=1&refresh=1m&viewPanel=77"

func FetchDataFromGrafana() {
	client, err := gapi.New("http://localhost:13000", gapi.Config{
		APIKey: ViewerAPIKey,
	})
	if err != nil {
		log.Fatalln(err)
	}
	dashboards, err := client.Dashboards()
	if err != nil {
		log.Fatalln(err)
	}
	dashboardsJSON, err := json.Marshal(dashboards)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(dashboardsJSON))
	dashboard, err := client.DashboardByUID(dashboards[0].UID)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(dashboard.Message)
}

type GrafanaPanelTargetModel struct {
	Expr   string `json:"expr"`
	Format string `json:"format"`
	Hide   bool   `json:"hide"`
}

type GrafanaPanelModel struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
	// TODO change another type
	Datasource  string                    `json:"datasource"`
	Description string                    `json:"description"`
	Targets     []GrafanaPanelTargetModel `json:"targets"`
}

type GrafanaService struct {
	grafanaClient *gapi.Client
}

func NewGrafanaService(grafanaClient *gapi.Client) *GrafanaService {
	return &GrafanaService{grafanaClient: grafanaClient}
}

type GrafanaDashboardTemplatingModel struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Query       string `json:"query"`
	Label       string `json:"label"`
	Hide        int    `json:"hide"`
	Refresh     int    `json:"refresh"`
	Regex       string `json:"regex"`
	Multi       bool   `json:"multi"`
	IncludeAdd  bool   `json:"includeAll"`
	SkipURLSync bool   `json:"skipUrlSync"`
}

type DashboardModel struct {
	Panels     []GrafanaPanelModel                          `json:"panels"`
	Templating map[string][]GrafanaDashboardTemplatingModel `json:"templating"`
}

func (it *GrafanaService) FetchPanelWithRawURL(ctx context.Context, rawURL string) (*GrafanaPanelIdentifier, *GrafanaPanelModel, error) {
	grafanaPanelIdentifier, err := FromRawURL(rawURL)
	if err != nil {
		return nil, nil, err
	}
	dashboard, err := it.grafanaClient.DashboardByUID(grafanaPanelIdentifier.DashboardUID)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "fetch grafana dashbaord %s", grafanaPanelIdentifier.DashboardUID)
	}
	var dashboardModel DashboardModel
	err = mapstructure.Decode(dashboard.Model, &dashboardModel)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "parse grafana dashboard model: %s", grafanaPanelIdentifier.DashboardUID)
	}
	panic("implement me")
}

func (it GrafanaService) FetchAvailableVariableOfPanel(ctx context.Context, panel entity.GrafanaPanel) ([]entity.GrafanaPanelVariable, error) {
	panic("implement me")

}

func (it GrafanaService) FetchMetricsFromPanel(ctx context.Context, panel entity.GrafanaPanel) (interface{}, error) {
	panic("implement me")
}

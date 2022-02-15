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
	"strings"
	"time"
)

const ViewerAPIKey = "eyJrIjoiNnNyc0tKQmJrc3B1UUlpTHhCVlFra2NBNFR2cHV1U1giLCJuIjoidiIsImlkIjoxfQ=="

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
	grafanaClient   *gapi.Client
	datasourceProxy DatasourceProxy
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

func (it *GrafanaService) FetchAvailableVariableOfPanel(ctx context.Context, dashboard DashboardModel) ([]entity.GrafanaPanelVariable, error) {
	var result []entity.GrafanaPanelVariable
	templatingModels, ok := dashboard.Templating["list"]
	if !ok {
		return nil, nil
	}
	for _, item := range templatingModels {
		result = append(result,
			entity.GrafanaPanelVariable{
				GrafanaPanelID: 0,
				VariableName:   item.Name,
				VariableType:   item.Type,
			},
		)
	}
	return result, nil
}

type VariableValue struct {
	Variable string `json:"variable"`
	Value    string `json:"value"`
}

func (it *GrafanaService) FetchMetricsFromPanel(ctx context.Context, panel entity.GrafanaPanel, datasourceID int64, variableValues []VariableValue) ([]DatasourceProxyResponse, error) {
	panelIdentifier, err := FromRawURL(panel.RawURL)
	if err != nil {
		return nil, err

	}

	dashboard, err := it.DashboardByUID(panelIdentifier.DashboardUID)
	if err != nil {
		return nil, err
	}

	var result []DatasourceProxyResponse
	for _, target := range dashboard.Panels[panelIdentifier.PanelID].Targets {
		renderedExpr := it.renderExpr(target.Expr, variableValues)
		timeSeries, err := it.datasourceProxy.DatasourceQuery(ctx, datasourceID, renderedExpr, time.Now().Add(-time.Hour), time.Now(), "2s")
		if err != nil {
			return nil, err
		}
		result = append(result, *timeSeries)
	}
	return result, nil
}

func (it *GrafanaService) DashboardByUID(dashboardUID string) (*DashboardModel, error) {
	dashboard, err := it.grafanaClient.DashboardByUID(dashboardUID)
	if err != nil {
		return nil, errors.Wrapf(err, "fetch grafana dashbaord %s", dashboardUID)
	}
	var dashboardModel DashboardModel
	err = mapstructure.Decode(dashboard.Model, &dashboardModel)
	if err != nil {
		return nil, errors.Wrapf(err, "parse grafana dashboard model: %s", dashboardUID)
	}
	return &dashboardModel, nil
}

func (it *GrafanaService) renderExpr(expr string, variableValues []VariableValue) string {
	var result = expr
	for _, item := range variableValues {
		result = strings.Replace(result, item.Variable, item.Value, -1)
	}
	return result
}

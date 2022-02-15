package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type DatasourceProxyResponse struct {
	Status string       `json:"status"`
	Data   ResponseData `json:"data"`
}

type Point struct {
	Time  time.Time
	Value float64
}

func (it *Point) UnmarshalJSON(bytes []byte) error {
	var temp []interface{}
	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return err
	}
	it.Time = time.Unix((temp[0]).(int64), 0)
	result, err := strconv.ParseFloat(temp[1].(string), 64)
	if err != nil {
		return err
	}
	it.Value = result

	return nil
}

func (it *Point) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d, %f]", it.Time.Unix(), it.Value)), nil
}

type Metric struct {
	Metric map[string]string `json:"metric"`
	Values []Point           `json:"values"`
}

type ResponseData struct {
	ResultType string          `json:"resultType"`
	Result     []Metric        `json:"result"`
	Data       []VariableValue `json:"data"`
}

type DatasourceProxy struct {
	// host should contain schema like http, https
	host  string
	token string
}

func NewDatasourceProxy(host string, token string) *DatasourceProxy {
	return &DatasourceProxy{host: host, token: token}
}

// DatasourceQuery might be useful with older grafana version < 8.0.0
func (it *DatasourceProxy) DatasourceQuery(ctx context.Context, datasourceID int64, query string, start, end time.Time, step string) (*DatasourceProxyResponse, error) {
	param :=
		fmt.Sprintf("query=%s&start=%d&end=%d&step=%s",
			url.QueryEscape(query),
			start.Unix(),
			end.Unix(),
			step,
		)
	request, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/datasources/proxy/%d/api/v1/query_range?%s",
			it.host,
			datasourceID,
			param,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", it.token))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var result DatasourceProxyResponse
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// TSDBQuery, used in 8.3.4
func (it DatasourceProxy) TSDBQuery() {
	panic("implement me")
}

package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

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

package grafana

import (
	"context"
	gapi "github.com/grafana/grafana-api-golang-client"
)

func Example_grafanaServiceFetchPanelWithRawURL() {
	client, err := gapi.New("http://localhost:13000", gapi.Config{
		APIKey: ViewerAPIKey,
	})
	if err != nil {
		panic(err)
	}
	service := NewGrafanaService(client)
	_, _, err = service.FetchPanelWithRawURL(context.TODO(), ExamplePanelRawURL)
	if err != nil {
		panic(err)
	}

	// Output:
}

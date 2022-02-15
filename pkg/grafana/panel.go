package grafana

import (
	"github.com/pkg/errors"
	"net/url"
	"path"
	"strconv"
)

const keyOrgID = "orgId"
const keyViewPanel = "viewPanel"

type GrafanaPanelIdentifier struct {
	OrgID        uint
	DashboardUID string
	PanelID      uint
}

func FromRawURL(rawURL string) (*GrafanaPanelIdentifier, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.Wrapf(err, "parse raw url to grafana panel identifier")
	}
	dashboardID := path.Base(path.Dir(parsedURL.Path))

	orgID := parsedURL.Query().Get(keyOrgID)
	orgIDInt, err := strconv.Atoi(orgID)
	if err != nil {
		return nil, errors.Wrapf(err, "parse raw url, invalid org id: %s", orgID)
	}

	viewPanel := parsedURL.Query().Get(keyViewPanel)
	viewPanelInt, err := strconv.Atoi(viewPanel)
	if err != nil {
		return nil, errors.Wrapf(err, "parse raw url, invalid view panel: %s", viewPanel)
	}
	return &GrafanaPanelIdentifier{
		DashboardUID: dashboardID,
		PanelID:      uint(viewPanelInt),
		OrgID:        uint(orgIDInt),
	}, nil
}

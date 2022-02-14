package grafana

import (
	"reflect"
	"testing"
)

func TestFromRawURL(t *testing.T) {
	type args struct {
		rawURL string
	}
	tests := []struct {
		name    string
		args    args
		want    *GrafanaPanelIdentifier
		wantErr bool
	}{
		{
			name: "example parse url",
			args: args{
				rawURL: ExamplePanelRawURL,
			},
			want: &GrafanaPanelIdentifier{
				DashboardUID: "rYdddlPWk",
				PanelID:      77,
				OrgID:        1,
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromRawURL(tt.args.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromRawURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromRawURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}


.PHONY: bin/reference-to-grafana
bin/reference-to-grafana:
	go build -o bin/reference-to-grafana ./cmd/reference-to-grafana

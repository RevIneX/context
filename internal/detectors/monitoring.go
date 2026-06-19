package detectors

import (
	"strings"

	"github.com/RevIneX/context/internal/data"
)

const (
	monPrometheus    = "Prometheus"
	monGrafana       = "Grafana"
	monZabbix        = "Zabbix"
	monELK           = "ELK"
	monElasticsearch = "Elasticsearch"
	monLogstash      = "Logstash"
	monKibana        = "Kibana"
	monDatadog       = "Datadog"
	monNewRelic      = "New Relic"
	monSplunk        = "Splunk"
	monJaeger        = "Jaeger"
	monOpenTelemetry = "OpenTelemetry"
)

type MonitoringDetector struct{}

func (d *MonitoringDetector) Name() string {
	return "monitoring"
}

func (d *MonitoringDetector) Detect(filePath string) []data.Finding {
	return d.DetectWithBuf(filePath, nil)
}

func (d *MonitoringDetector) DetectWithBuf(filePath string, buf []data.Finding) []data.Finding {
	if filePath == "" {
		return buf
	}

	lastSlash := strings.LastIndexByte(filePath, '/')
	lastBackslash := strings.LastIndexByte(filePath, '\\')
	if lastBackslash > lastSlash {
		lastSlash = lastBackslash
	}
	baseName := filePath[lastSlash+1:]

	if baseName == "" {
		return buf
	}

	var tech string
	switch baseName {
	case "prometheus.yml", "prometheus.yaml":
		tech = monPrometheus
	case "grafana.yaml", "grafana.yml", "grafana.json":
		tech = monGrafana
	case "zabbix.conf", "zabbix.yaml", "zabbix.yml":
		tech = monZabbix
	case "elasticsearch.yml", "elasticsearch.yaml":
		tech = monElasticsearch
	case "logstash.conf", "logstash.yml", "logstash.yaml":
		tech = monLogstash
	case "kibana.yml", "kibana.yaml":
		tech = monKibana
	case "elk.yml", "elk.yaml":
		tech = monELK
	case "datadog.yaml", "datadog.yml":
		tech = monDatadog
	case "newrelic.yml", "newrelic.yaml":
		tech = monNewRelic
	case "splunk.yml", "splunk.yaml":
		tech = monSplunk
	case "jaeger.yml", "jaeger.yaml":
		tech = monJaeger
	case "opentelemetry.yml", "opentelemetry.yaml":
		tech = monOpenTelemetry
	default:
		lowName := strings.ToLower(baseName)
		switch lowName {
		case "prometheus.yml", "prometheus.yaml":
			tech = monPrometheus
		case "grafana.yaml", "grafana.yml", "grafana.json":
			tech = monGrafana
		case "zabbix.conf", "zabbix.yaml", "zabbix.yml":
			tech = monZabbix
		case "elasticsearch.yml", "elasticsearch.yaml":
			tech = monElasticsearch
		case "logstash.conf", "logstash.yml", "logstash.yaml":
			tech = monLogstash
		case "kibana.yml", "kibana.yaml":
			tech = monKibana
		case "elk.yml", "elk.yaml":
			tech = monELK
		case "datadog.yaml", "datadog.yml":
			tech = monDatadog
		case "newrelic.yml", "newrelic.yaml":
			tech = monNewRelic
		case "splunk.yml", "splunk.yaml":
			tech = monSplunk
		case "jaeger.yml", "jaeger.yaml":
			tech = monJaeger
		case "opentelemetry.yml", "opentelemetry.yaml":
			tech = monOpenTelemetry
		}
	}

	if tech != "" {
		buf = append(buf, data.NewFinding(data.CatMonitoring, tech, "", filePath, 1, nil))
	}
	return buf
}
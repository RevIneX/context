package detectors

import (
	"strings"

	"github.com/RevIneX/context/internal/data"
)

const (
	orchKubernetes = "Kubernetes"
	orchSwarm      = "Docker Swarm"
	orchNomad      = "Nomad"
	orchOpenShift  = "OpenShift"
	orchRancher    = "Rancher"
	orchMesos      = "Mesos"
	orchEKS        = "EKS"
	orchGKE        = "GKE"
	orchAKS        = "AKS"
)

type OrchestrationDetector struct{}

func (d *OrchestrationDetector) Name() string {
	return "orchestration"
}

func (d *OrchestrationDetector) Detect(filePath string) []data.Finding {
	return d.DetectWithBuf(filePath, nil)
}

func (d *OrchestrationDetector) DetectWithBuf(filePath string, buf []data.Finding) []data.Finding {
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
	case "deployment.yaml", "deployment.yml":
		tech = orchKubernetes
	case "service.yaml", "service.yml":
		tech = orchKubernetes
	case "k8s.yaml", "k8s.yml":
		tech = orchKubernetes
	case "kubernetes.yaml", "kubernetes.yml":
		tech = orchKubernetes
	case "swarm.yml", "swarm.yaml":
		tech = orchSwarm
	case "docker-swarm.yml", "docker-swarm.yaml":
		tech = orchSwarm
	case "nomad.hcl", "nomad.yaml", "nomad.yml":
		tech = orchNomad
	case "openshift.yaml", "openshift.yml":
		tech = orchOpenShift
	case "rancher.yaml", "rancher.yml":
		tech = orchRancher
	case "mesos.yaml", "mesos.yml":
		tech = orchMesos
	default:
		lowName := strings.ToLower(baseName)
		switch lowName {
		case "deployment.yaml", "deployment.yml":
			tech = orchKubernetes
		case "service.yaml", "service.yml":
			tech = orchKubernetes
		case "k8s.yaml", "k8s.yml":
			tech = orchKubernetes
		case "kubernetes.yaml", "kubernetes.yml":
			tech = orchKubernetes
		case "swarm.yml", "swarm.yaml":
			tech = orchSwarm
		case "docker-swarm.yml", "docker-swarm.yaml":
			tech = orchSwarm
		case "nomad.hcl", "nomad.yaml", "nomad.yml":
			tech = orchNomad
		case "openshift.yaml", "openshift.yml":
			tech = orchOpenShift
		case "rancher.yaml", "rancher.yml":
			tech = orchRancher
		case "mesos.yaml", "mesos.yml":
			tech = orchMesos
		}
	}

	if tech != "" {
		buf = append(buf, data.NewFinding(data.CatOrchestration, tech, "", filePath, 1, nil))
	}
	return buf
}
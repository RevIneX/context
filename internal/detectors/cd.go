package detectors

import (
	"strings"

	"github.com/RevIneX/context/internal/data"
)

const (
	cdArgoCD    = "ArgoCD"
	cdSpinnaker = "Spinnaker"
	cdFluxCD    = "FluxCD"
	cdCodeFresh = "CodeFresh"
	cdJenkinsX  = "Jenkins X"
	cdHarness   = "Harness"
	cdOctopus   = "Octopus Deploy"
)

type CDDetector struct{}

func (d *CDDetector) Name() string {
	return "cd"
}

func (d *CDDetector) Detect(filePath string) []data.Finding {
	return d.DetectWithBuf(filePath, nil)
}

func (d *CDDetector) DetectWithBuf(filePath string, buf []data.Finding) []data.Finding {
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
	case "argo-app.yaml", "argo-app.yml", "argocd.yaml", "argocd.yml":
		tech = cdArgoCD
	case "spinnaker.yml", "spinnaker.yaml":
		tech = cdSpinnaker
	case "fluxcd.yaml", "fluxcd.yml":
		tech = cdFluxCD
	case "codefresh.yml", "codefresh.yaml":
		tech = cdCodeFresh
	case "jenkins-x.yml", "jenkins-x.yaml":
		tech = cdJenkinsX
	case "harness.yaml", "harness.yml":
		tech = cdHarness
	case "octopus.yaml", "octopus.yml":
		tech = cdOctopus
	default:
		lowName := strings.ToLower(baseName)
		switch lowName {
		case "argo-app.yaml", "argo-app.yml", "argocd.yaml", "argocd.yml":
			tech = cdArgoCD
		case "spinnaker.yml", "spinnaker.yaml":
			tech = cdSpinnaker
		case "fluxcd.yaml", "fluxcd.yml":
			tech = cdFluxCD
		case "codefresh.yml", "codefresh.yaml":
			tech = cdCodeFresh
		case "jenkins-x.yml", "jenkins-x.yaml":
			tech = cdJenkinsX
		case "harness.yaml", "harness.yml":
			tech = cdHarness
		case "octopus.yaml", "octopus.yml":
			tech = cdOctopus
		}
	}

	if tech != "" {
		buf = append(buf, data.NewFinding(data.CatCD, tech, "", filePath, 1, nil))
	}
	return buf
}
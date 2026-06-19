package detectors

import (
	"strings"

	"github.com/RevIneX/context/internal/data"
)

const (
	ciJenkins            = "Jenkins"
	ciGitHubActions      = "GitHub Actions"
	ciGitLabCI           = "GitLab CI"
	ciCircleCI           = "CircleCI"
	ciTravis             = "Travis CI"
	ciDrone              = "Drone"
	ciBuildkite          = "Buildkite"
	ciTeamCity           = "TeamCity"
	ciAzurePipelines     = "Azure Pipelines"
	ciBitbucketPipelines = "Bitbucket Pipelines"
)

var ciFiles = map[string]string{
	"jenkinsfile":             ciJenkins,
	".gitlab-ci.yml":          ciGitLabCI,
	".travis.yml":             ciTravis,
	".drone.yml":              ciDrone,
	"buildkite.yml":           ciBuildkite,
	"buildkite.yaml":          ciBuildkite,
	"bitbucket-pipelines.yml": ciBitbucketPipelines,
	"azure-pipelines.yml":     ciAzurePipelines,
	"azure-pipelines.yaml":    ciAzurePipelines,
}

type CIDetector struct{}

func (d *CIDetector) Name() string {
	return "ci"
}

func (d *CIDetector) Detect(filePath string) []data.Finding {
	return d.DetectWithBuf(filePath, nil)
}

func (d *CIDetector) DetectWithBuf(filePath string, buf []data.Finding) []data.Finding {
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
	var ok bool

	tech, ok = ciFiles[baseName]
	if !ok {
		lowName := strings.ToLower(baseName)
		tech, ok = ciFiles[lowName]
		if !ok {
			if lastSlash > 0 {
				parentStart := lastSlash - 1
				for parentStart >= 0 && filePath[parentStart] != '/' && filePath[parentStart] != '\\' {
					parentStart--
				}
				parentDir := filePath[parentStart+1 : lastSlash]
				if parentDir == "workflows" || parentDir == "WORKFLOWS" {
					if parentStart > 0 {
						grandStart := parentStart - 1
						for grandStart >= 0 && filePath[grandStart] != '/' && filePath[grandStart] != '\\' {
							grandStart--
						}
						grandDir := filePath[grandStart+1 : parentStart]
						if grandDir == ".github" {
							if strings.HasSuffix(lowName, ".yml") || strings.HasSuffix(lowName, ".yaml") {
								tech, ok = ciGitHubActions, true
							}
						}
					}
				}

				if !ok && (parentDir == ".circleci" || parentDir == ".CIRCLECI") {
					if lowName == "config.yml" || lowName == "config.yaml" {
						tech, ok = ciCircleCI, true
					}
				}
			}
		}
	}

	if ok {
		buf = append(buf, data.NewFinding(data.CatCI, tech, "", filePath, 1, nil))
	}
	return buf
}
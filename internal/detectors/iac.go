package detectors

import (
	"strings"

	"github.com/RevIneX/context/internal/data"
)

const (
	iacTerraform      = "Terraform"
	iacPulumi         = "Pulumi"
	iacCloudFormation = "CloudFormation"
	iacAnsible        = "Ansible"
	iacHelm           = "Helm"
	iacVagrant        = "Vagrant"
	iacPacker         = "Packer"
	iacCrossplane     = "Crossplane"
	iacCDK            = "CDK"
	iacCDKTF          = "CDKTF"
)

var iacFiles = map[string]string{
	"main.tf":             iacTerraform,
	"terraform.tfvars":    iacTerraform,
	"pulumi.yaml":         iacPulumi,
	"cloudformation.json": iacCloudFormation,
	"cloudformation.yaml": iacCloudFormation,
	"cloudformation.yml":  iacCloudFormation,
	"site.yml":            iacAnsible,
	"playbook.yml":        iacAnsible,
	"playbook.yaml":       iacAnsible,
	"chart.yaml":          iacHelm,
	"chart.yml":           iacHelm,
	"vagrantfile":         iacVagrant,
	"packer.json":         iacPacker,
	"crossplane.yaml":     iacCrossplane,
	"crossplane.yml":      iacCrossplane,
}

var iacExtensions = map[string]string{
	".tf":     iacTerraform,
	".hcl":    iacTerraform,
	".tfvars": iacTerraform,
	".pulumi": iacPulumi,
	".cdk":    iacCDK,
	".cdktf":  iacCDKTF,
}

type IaCDetector struct{}

func (d *IaCDetector) Name() string {
	return "iac"
}

func (d *IaCDetector) DetectWithBuf(filePath string, buf []data.Finding) []data.Finding {
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

	tech, ok = iacFiles[baseName]
	if !ok {
		lowName := strings.ToLower(baseName)
		tech, ok = iacFiles[lowName]
		if !ok {
			if dot := strings.LastIndexByte(lowName, '.'); dot != -1 {
				ext := lowName[dot:]
				tech, ok = iacExtensions[ext]
			}
		}
	}

	if ok {
		buf = append(buf, data.NewFinding(data.CatIaC, tech, "", filePath, 1, nil))
	}
	return buf
}

func (d *IaCDetector) Detect(filePath string) []data.Finding {
	return d.DetectWithBuf(filePath, nil)
}
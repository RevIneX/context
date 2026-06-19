package main

import (
	"fmt"
	"os"
	"time"

	"github.com/RevIneX/context/internal/analyzer"
	"github.com/RevIneX/context/internal/detectors"
	"github.com/RevIneX/context/internal/formatting"
	"github.com/RevIneX/context/internal/stages"
)

func main() {
	start := time.Now()

	var root string
	if len(os.Args) > 1 {
		root = os.Args[1]
	} else {
		var err error
		root, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	files, err := analyzer.WalkProject(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	deps := []stages.Detector{
		&detectors.LanguageDetector{},
		&detectors.ContainerDetector{},
		&detectors.ArchitectureDetector{},
		&detectors.DatabaseDetector{},
		&detectors.FrameworkDetector{},
		&detectors.WebserverDetector{},
		&detectors.QueueDetector{},
		&detectors.CacheDetector{},
		&detectors.APIDetector{},
		&detectors.IaCDetector{},
		&detectors.CIDetector{},
		&detectors.CDDetector{},
		&detectors.MonitoringDetector{},
		&detectors.OrchestrationDetector{},
	}

	runner := stages.NewRunner(deps, 0)
	findings := runner.Run(files)

	output := formatting.Format(findings)
	if output != "" {
		fmt.Print(output)
	}

	elapsed := time.Since(start)
	fmt.Fprintf(os.Stderr, "\nscan complete: %d findings in %v\n", len(findings), elapsed)
}

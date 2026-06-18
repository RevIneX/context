package main

import (
	"github.com/RevIneX/context/internal/analyzer"
	"github.com/RevIneX/context/internal/detectors"
	"github.com/RevIneX/context/internal/formatting"
	"github.com/RevIneX/context/internal/stages"
	"fmt"
	"os"
	"time"
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

	var deps []stages.Detector
	deps = append(deps,
		&detectors.LanguageDetector{},
		&detectors.ContainerDetector{},
		&detectors.ArchitectureDetector{},
	)

	runner := stages.NewRunner(deps, 0)
	findings := runner.Run(files)

	output := formatting.Format(findings)

	if output != "" {
		fmt.Print(output)
	}

	elapsed := time.Since(start)
	fmt.Fprintf(os.Stderr, "\nscan complete: %d findings in %v\n", len(findings), elapsed)
}
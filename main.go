package main

import (
	"deprec/agent"
	"deprec/configuration"
	"deprec/model"
	"log"
)

func main() {
	log.Printf("DepRec run started...")

	config := configuration.Load()

	sbom := []*model.Dependency{
		{
			Name:     "log4j",
			Version:  "1.2.6",
			MetaData: map[string]string{"vcs": "https://github.com/apache/logging-log4j1"},
		},
		{
			Name:     "log4j",
			Version:  "2.19",
			MetaData: map[string]string{"vcs": "https://github.com/apache/logging-log4j2"},
		}}

	var result []float64
	for _, dep := range sbom {
		log.Printf("Running Agent for %s:%s", dep.Name, dep.Version)

		a := agent.NewAgent(dep, config)
		r := a.Start().Result
		result = append(result, r)
	}

	log.Printf("...DepRec run done")
	log.Printf("Results: %f", result)
}

package agent

import (
	"deprec/extraction"
	"deprec/model"
	"log"
)

type Agent struct {
	Dependency *model.Dependency
	DataModel  *model.DataModel
}

func NewAgent(dependency *model.Dependency) *Agent {
	agent := Agent{Dependency: dependency, DataModel: &model.DataModel{}}
	return &agent
}

func (agent *Agent) Start() model.AgentResult {
	log.Printf("Starting Extraction...")
	agent.Extraction()
	log.Printf("...Extraction complete")

	log.Printf("Starting Combination & Conclusion...")
	agent.CombinationAndConclusion()
	log.Printf("...Combination & Conclusion complete")

	return model.AgentResult{
		Dependency: agent.Dependency,
		Result:     0.0,
	}
}

func (agent *Agent) Extraction() {

	extractor := extraction.NewRepositoryExtractor(agent.Dependency)
	extractor.Extract(agent.DataModel)
}

func (agent *Agent) CombinationAndConclusion() {

}

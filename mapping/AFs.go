package mapping

import (
	"deprec/model"
	"strings"
)

func Network(model *model.DataModel) float64 {

	var result float64
	if model.Repository == nil {
		return 0
	}

	for _, contributor := range model.Repository.Contributors {
		result += float64(contributor.Repositories)
		result += float64(contributor.Organizations)
	}

	result += float64(len(model.Repository.Contributors))

	return result
}

func Interconnectedness(model *model.DataModel) float64 {
	return Network(model)*0.5 + 1
}

func DeityGiven(model *model.DataModel) float64 {

	archived := model.Repository.Archivation
	if archived {
		return 1
	}

	readme := model.Repository.ReadMe

	if strings.Contains(readme, "deprecated") || strings.Contains(strings.ToLower(readme), "end-of-life") {
		return 1
	}

	return 0

}

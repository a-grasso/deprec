package mapping

import "deprec/model"

func Network(model *model.DataModel) float64 {
	var result float64

	for _, contributor := range model.Repository.Contributors {
		result += float64(contributor.Repositories)
	}

	result += float64(len(model.Repository.Contributors))

	return result
}

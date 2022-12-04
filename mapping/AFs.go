package mapping

import (
	"deprec/model"
	"log"
)

func Network(model *model.DataModel) float64 {

	var result float64
	if model.Repository == nil {
		return 0
	}

	log.Printf("calcing network of repo %s", model.Repository.Owner)

	for _, contributor := range model.Repository.Contributors {
		result += float64(contributor.Repositories)
		result += float64(contributor.Organizations)
	}

	result += float64(len(model.Repository.Contributors))

	return result
}

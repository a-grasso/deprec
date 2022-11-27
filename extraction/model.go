package extraction

import "deprec/model"

type Extractor interface {
	// Extract information into a DataModel
	Extract(model model.DataModel)
}

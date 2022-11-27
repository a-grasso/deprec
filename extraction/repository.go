package extraction

import "deprec/model"

type RepositoryExtractor struct {
	Repository string
}

func NewRepositoryExtractor(dependency *model.Dependency) *RepositoryExtractor {
	return &RepositoryExtractor{Repository: dependency.MetaData["vcs"]}
}

func (re *RepositoryExtractor) Extract(dataModel *model.DataModel) {

	repository := &model.Repository{Contributors: &model.Contributors{}}

	dataModel.Repository = repository
}

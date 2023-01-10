package extraction

import (
	"deprec/cache"
	"deprec/configuration"
	"deprec/logging"
	"deprec/model"
	"deprec/ossindexapi"
)

type OSSIndexExtractor struct {
	PackageURL string
	Config     *configuration.Configuration
	Client     *ossindexapi.ClientWrapper
}

func NewOSSIndexExtractor(dependency *model.Dependency, config *configuration.Configuration) *OSSIndexExtractor {

	cache := cache.NewCache(config.MongoDB)
	client := ossindexapi.NewClient(config.OSSIndex)

	wrapper := ossindexapi.NewClientWrapper(client, cache)

	packageURL := dependency.PackageURL

	return &OSSIndexExtractor{
		PackageURL: packageURL,
		Config:     config,
		Client:     wrapper,
	}
}

func (ossie *OSSIndexExtractor) Extract(dataModel *model.DataModel) {
	logging.SugaredLogger.Infof("extracting ossindex '%s'", ossie.PackageURL)

	index := &model.VulnerabilityIndex{}

	reports, _ := ossie.Client.GetComponentReport(ossie.PackageURL)

	vulnerabilities := 0

	for _, report := range reports {
		vulnerabilities += len(report.Vulnerabilities)
	}

	index.TotalVulnerabilitiesCount = vulnerabilities

	dataModel.VulnerabilityIndex = index
}

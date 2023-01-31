package extraction

import (
	"github.com/a-grasso/deprec/cache"
	"github.com/a-grasso/deprec/configuration"
	"github.com/a-grasso/deprec/logging"
	"github.com/a-grasso/deprec/model"
	"github.com/a-grasso/deprec/ossindexapi"
	"net/http"
	"strings"
)

type OSSIndexExtractor struct {
	PackageURL string
	Config     configuration.OSSIndex
	Client     *ossindexapi.ClientWrapper
}

func NewOSSIndexExtractor(dependency model.Dependency, config configuration.OSSIndex, cache *cache.Cache) *OSSIndexExtractor {

	client := ossindexapi.NewClient(config)

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

	purl := strings.Split(ossie.PackageURL, "?type")[0]

	reports, err := ossie.Client.GetComponentReport(purl)
	if err != nil {
		return
	}
	if len(reports) != 1 {
		return
	}

	componentReport := reports[0]

	resp, err := http.Get(componentReport.Reference)
	if err != nil {
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		return
	}

	index.TotalVulnerabilitiesCount = len(componentReport.Vulnerabilities)

	dataModel.VulnerabilityIndex = index
}

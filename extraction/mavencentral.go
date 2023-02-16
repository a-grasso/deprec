package extraction

import (
	"fmt"
	"github.com/a-grasso/deprec/cache"
	"github.com/a-grasso/deprec/logging"
	"github.com/a-grasso/deprec/mavencentralapi"
	"github.com/a-grasso/deprec/model"
	"github.com/thoas/go-funk"
	"github.com/vifraa/gopom"
	"time"
)

type MavenCentralExtractor struct {
	DependencyName string
	SHA1           string
	Client         *mavencentralapi.ClientWrapper
}

func NewMavenCentralExtractor(dependency model.Dependency, cache *cache.Cache) *MavenCentralExtractor {

	client := mavencentralapi.NewClient()

	wrapper := mavencentralapi.NewClientWrapper(client, cache)

	sha1 := dependency.Hashes[model.SHA1]

	return &MavenCentralExtractor{
		DependencyName: dependency.Name,
		SHA1:           sha1,
		Client:         wrapper,
	}
}

func (mce *MavenCentralExtractor) Extract(dataModel *model.DataModel) {
	logging.SugaredLogger.Infof("extracting maven central '%s' with SHA-1 '%s'", mce.DependencyName, mce.SHA1)

	search, err := mce.Client.SearchMavenCentralSHA1(mce.SHA1)

	if err != nil {
		logging.SugaredLogger.Debugf("could not search maven central '%s' with SHA-1 '%s'", mce.DependencyName, mce.SHA1)
		return
	}

	if len(search.Response.Docs) == 0 {
		return
	}

	response := search.Response.Docs[0]

	version := response.V
	artifactId := response.A
	groupId := response.G
	var msToNs int64 = 1000000
	timestamp := time.Unix(0, response.Timestamp*msToNs)

	library := mce.extractLibrary(groupId, artifactId)

	artifact := mce.extractArtifact(groupId, artifactId, version, timestamp)

	dataModel.Distribution = &model.Distribution{
		Library:  library,
		Artifact: artifact,
	}
}

func (mce *MavenCentralExtractor) extractArtifact(groupId string, artifactId string, version string, date time.Time) *model.Artifact {
	pom, err := mce.Client.GetArtifactPom(groupId, artifactId, version)
	if err != nil {
		logging.SugaredLogger.Debugf("could not get artifact pom for '%s' with SHA-1 '%s'", mce.DependencyName, mce.SHA1)
		return nil
	}

	repos := funk.Map(pom.Repositories, func(r gopom.Repository) string { return r.Name }).([]string)

	dependencies := collectDependencies(pom)

	contributors := funk.Map(pom.Contributors, func(c gopom.Contributor) string { return c.Email }).([]string)

	developers := funk.Map(pom.Developers, func(c gopom.Developer) string { return c.Email }).([]string)

	licenses := funk.Map(pom.Licenses, func(l gopom.License) string { return l.Name }).([]string)

	mailingLists := funk.Map(pom.MailingLists, func(ml gopom.MailingList) string { return ml.Name }).([]string)

	return &model.Artifact{
		Version:              version,
		ArtifactRepositories: repos,
		Date:                 date,
		Vulnerabilities:      nil,
		Dependents:           nil,
		Dependencies:         dependencies,
		DeprecationWarning:   false,
		Contributors:         contributors,
		Developers:           developers,
		Organization:         pom.Organization.Name,
		Licenses:             licenses,
		MailingLists:         mailingLists,
		Description:          pom.Description,
	}
}

func collectDependencies(pom *gopom.Project) []string {
	var dependencies []string

	var pomDependencies []gopom.Dependency
	var pomPlugins []gopom.Plugin

	pomDependencies = append(pomDependencies, pom.Dependencies...)
	pomDependencies = append(pomDependencies, pom.DependencyManagement.Dependencies...)

	pomPlugins = append(pomPlugins, pom.Build.Plugins...)
	pomPlugins = append(pomPlugins, pom.Build.PluginManagement.Plugins...)

	for _, dependency := range pomDependencies {
		dep := fmt.Sprintf("%s-%s-%s", dependency.GroupID, dependency.ArtifactID, dependency.Version)

		if funk.Contains(dependencies, dep) {
			continue
		}

		dependencies = append(dependencies, dep)
	}

	for _, plugin := range pomPlugins {
		dep := fmt.Sprintf("%s-%s-%s", plugin.GroupID, plugin.ArtifactID, plugin.Version)

		if funk.Contains(dependencies, dep) {
			continue
		}

		dependencies = append(dependencies, dep)
	}

	return dependencies
}

func (mce *MavenCentralExtractor) extractLibrary(groupId string, artifactId string) *model.Library {
	metadata, err := mce.Client.GetLibraryMetadata(groupId, artifactId)
	if err != nil {
		logging.SugaredLogger.Debugf("could not get library metadata for '%s' with SHA-1 '%s'", mce.DependencyName, mce.SHA1)
		return nil
	}

	lastUpdated, err := time.Parse("20060102150405", metadata.Versioning.LastUpdated)
	if err != nil {
		lastUpdated = time.Time{}
	}

	return &model.Library{
		Ranking:       nil,
		Licenses:      nil,
		UsedBy:        nil,
		Versions:      metadata.Versioning.Versions.Version,
		LastUpdated:   lastUpdated,
		LatestVersion: metadata.Versioning.Latest,
		LatestRelease: metadata.Versioning.Release,
	}
}

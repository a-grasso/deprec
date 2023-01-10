package mavencentralapi

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/vifraa/gopom"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	BaseURLSHASearch      string
	BaseURLBrowseArtifact string
	BaseURLBrowseLibrary  string
	BasePOMName           string
	MetadataName          string
}

func NewClient() *Client {

	return &Client{
		BaseURLSHASearch:      "https://search.maven.org/solrsearch/select?q=1:%s&rows=20&wt=json",
		BaseURLBrowseArtifact: "https://repo1.maven.org/maven2/%s/%s/%s/%s.%s",
		BaseURLBrowseLibrary:  "https://repo1.maven.org/maven2/%s/%s/%s.%s",
		BasePOMName:           "%s-%s",
		MetadataName:          "maven-metadata",
	}
}

type MavenCentralSearch struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
		Params struct {
			Q       string `json:"q"`
			Core    string `json:"core"`
			Indent  string `json:"indent"`
			Fl      string `json:"fl"`
			Start   string `json:"start"`
			Sort    string `json:"sort"`
			Rows    string `json:"rows"`
			Wt      string `json:"wt"`
			Version string `json:"version"`
		} `json:"params"`
	} `json:"responseHeader"`
	Response struct {
		NumFound int `json:"numFound"`
		Start    int `json:"start"`
		Docs     []struct {
			ID        string   `json:"id"`
			G         string   `json:"g"`
			A         string   `json:"a"`
			V         string   `json:"v"`
			P         string   `json:"p"`
			Timestamp int64    `json:"timestamp"`
			Ec        []string `json:"ec"`
			Tags      []string `json:"tags"`
		} `json:"docs"`
	} `json:"response"`
}

func (c *Client) SearchMavenCentralSHA1(sha1 string) (*MavenCentralSearch, error) {
	url := fmt.Sprintf(c.BaseURLSHASearch, sha1)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var j MavenCentralSearch
	err = json.NewDecoder(resp.Body).Decode(&j)
	if err != nil {
		return nil, err
	}

	return &j, err
}

func (c *Client) GetArtifactPom(groupId, artifactId, version string) (*gopom.Project, error) {

	extension := "pom"

	groupId = strings.ReplaceAll(groupId, ".", "/")
	artifactId = strings.ReplaceAll(artifactId, ".", "/")

	url := fmt.Sprintf(c.BaseURLBrowseArtifact, groupId, artifactId, version, fmt.Sprintf(c.BasePOMName, artifactId, version), extension)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	pomString := string(body)

	var pom gopom.Project
	err = xml.Unmarshal([]byte(pomString), &pom)
	if err != nil {
		return nil, err
	}

	return &pom, nil
}

type Metadata struct {
	XMLName      xml.Name `xml:"metadata"`
	Text         string   `xml:",chardata"`
	ModelVersion string   `xml:"modelVersion,attr"`
	GroupId      string   `xml:"groupId"`
	ArtifactId   string   `xml:"artifactId"`
	Versioning   struct {
		Text     string `xml:",chardata"`
		Latest   string `xml:"latest"`
		Release  string `xml:"release"`
		Versions struct {
			Text    string   `xml:",chardata"`
			Version []string `xml:"version"`
		} `xml:"versions"`
		LastUpdated string `xml:"lastUpdated"`
	} `xml:"versioning"`
}

func (c *Client) GetLibraryMetadata(groupId, artifactId string) (*Metadata, error) {

	extension := "xml"

	groupId = strings.ReplaceAll(groupId, ".", "/")
	artifactId = strings.ReplaceAll(artifactId, ".", "/")

	url := fmt.Sprintf(c.BaseURLBrowseLibrary, groupId, artifactId, c.MetadataName, extension)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	pomString := string(body)

	var pom Metadata
	err = xml.Unmarshal([]byte(pomString), &pom)
	if err != nil {
		return nil, err
	}

	return &pom, nil
}

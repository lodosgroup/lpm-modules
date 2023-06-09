package template

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	common "lpm_builder/pkg/common"
	"net/url"
	"regexp"
)

type Template struct {
	Name                  string         `json:"name"`
	Description           string         `json:"description"`
	Maintainer            string         `json:"maintainer"`
	SourceRepository      string         `json:"source_repository"`
	Homepage              string         `json:"homepage"`
	Arch                  string         `json:"arch"`
	Kind                  string         `json:"kind"`
	FileChecksumAlgo      string         `json:"file_checksum_algo"`
	Tags                  []string       `json:"tags"`
	Version               common.Version `json:"version"`
	License               string         `json:"license"`
	MandatoryDependencies Dependencies   `json:"mandatory_dependencies"`
	SuggestedDependencies Dependencies   `json:"suggested_dependencies"`
}

type Dependencies struct {
	Build   []common.Dependency `json:"build"`
	Runtime []common.Dependency `json:"runtime"`
}

func (template *Template) validate() error {
	common.Logger.Println("Validating template fields")

	// field lengths
	{
		if len(template.Description) > 200 {
			return errors.New("description exceeds maximum length of 155 characters")
		}

		if len(template.Tags) > 5 {
			return errors.New("tags exceeds maximum length of 5 items")
		}
	}

	// Package name
	{
		pkgNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

		if !pkgNameRegex.MatchString(template.Name) {
			return errors.New("Package name can only contain English alphabets, numbers, '-' and '_' characters")
		}
	}

	// Architecture
	{
		var supportedArchitectures []string = []string{
			"amd64",
			"aarch64",
			"i686",
			"arm",
			"armv7",
			"mips",
			"mips64",
			"mips64el",
			"mipsel",
			"powerpc",
			"powerpc64",
			"powerpc64le",
			"riscv64gc",
			"s390x",
		}

		if !common.Contains(supportedArchitectures, template.Arch) {
			return errors.New("Unsupported architecture.")
		}

	}

	// File checksum algorithm
	{
		var supportedAlgorithms []string = []string{
			"md5",
			"sha256",
			"sha512",
		}

		if !common.Contains(supportedAlgorithms, template.FileChecksumAlgo) {
			return errors.New("Unsupported checksum algorithm.")
		}
	}

	// URLs
	{
		if len(template.SourceRepository) > 0 {
			_, err := url.ParseRequestURI(template.SourceRepository)
			common.FailOnError(err, "Invalid URL in source_repository.")
		}

		_, err := url.ParseRequestURI(template.Homepage)
		common.FailOnError(err, "Invalid URL in homepage.")
	}

	return nil
}

func DeserializeTemplate(templateDirPath string) Template {
	common.Logger.Println("Deserializing template fields")

	const templateLeafPath = "/template"

	var template = Template{
		MandatoryDependencies: Dependencies{
			Build:   []common.Dependency{},
			Runtime: []common.Dependency{},
		},
		SuggestedDependencies: Dependencies{
			Build:   []common.Dependency{},
			Runtime: []common.Dependency{},
		},
	}

	template_json_content, err := ioutil.ReadFile(templateDirPath + templateLeafPath)
	common.FailOnError(err, "Failed reading template json file")

	err = json.Unmarshal(template_json_content, &template)
	common.FailOnError(err, "Failed reading template json file")

	err = template.validate()
	common.FailOnError(err)

	for i := range template.MandatoryDependencies.Runtime {
		if len(template.MandatoryDependencies.Runtime[i].Version.Condition) == 0 {
			template.MandatoryDependencies.Runtime[i].Version.Condition = ">="
		}
	}

	for i := range template.MandatoryDependencies.Build {
		if len(template.MandatoryDependencies.Build[i].Version.Condition) == 0 {
			template.MandatoryDependencies.Build[i].Version.Condition = ">="
		}
	}

	for i := range template.SuggestedDependencies.Runtime {
		if len(template.SuggestedDependencies.Runtime[i].Version.Condition) == 0 {
			template.SuggestedDependencies.Runtime[i].Version.Condition = ">="
		}
	}

	for i := range template.SuggestedDependencies.Build {
		if len(template.SuggestedDependencies.Build[i].Version.Condition) == 0 {
			template.SuggestedDependencies.Build[i].Version.Condition = ">="
		}
	}

	return template
}

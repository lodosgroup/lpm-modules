package template

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	common "lpm_builder/pkg/common"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func readAndFillString(prompt string, reader *bufio.Reader, to_be_filled *string) {
	fmt.Printf(prompt)
	value, err := reader.ReadString('\n')
	common.FailOnError(err, "Failed on reading value from stdin")
	*to_be_filled = strings.TrimSpace(value)

}

func CreateInteractively(c *cli.Context, template_name string, out_path string) {
	template := &common.Template{
		Name:                template_name,
		Description:         "",
		Maintainer:          "",
		SourcePkg:           "",
		Repository:          "",
		Homepage:            "",
		Arch:                "",
		Kind:                "",
		Tags:                []string{},
		Version:             common.Version{},
		License:             "",
		RuntimeDependencies: []common.Dependency{},
		RuntimeSuggestions:  []common.Dependency{},
		BuildDependencies:   []common.Dependency{},
	}

	reader := bufio.NewReader(os.Stdin)

	readAndFillString(fmt.Sprintf("Description of '%s': ", template_name), reader, &template.Description)
	readAndFillString(fmt.Sprintf("Maintainer of '%s': ", template_name), reader, &template.Maintainer)
	readAndFillString(fmt.Sprintf("Source package of '%s': ", template_name), reader, &template.SourcePkg)
	readAndFillString(fmt.Sprintf("Repository of '%s': ", template_name), reader, &template.Repository)
	readAndFillString(fmt.Sprintf("Homepage of '%s': ", template_name), reader, &template.Homepage)
	readAndFillString(fmt.Sprintf("Architecture of '%s': ", template_name), reader, &template.Arch)
	readAndFillString(fmt.Sprintf("Kind of '%s': ", template_name), reader, &template.Kind)
	readAndFillString(fmt.Sprintf("License of '%s': ", template_name), reader, &template.License)

	fmt.Printf("Tags of '%s': ", template_name)
	value, err := reader.ReadString('\n')
	common.FailOnError(err, "Failed on reading value from stdin")
	template.Tags = strings.Split(strings.TrimSpace(value), " ")

	common.SetReadableVersion(&template.Version)

	template_json, err := json.MarshalIndent(template, "", " ")
	common.FailOnError(err, "Failed on serializing template to json string")

	err = ioutil.WriteFile(out_path+"template.json", template_json, 0644)
	common.FailOnError(err, "Failed on writing to file")
}

func CreateDefault(c *cli.Context, template_name string, out_path string) {
	template := &common.Template{
		Name:                template_name,
		Description:         "",
		Maintainer:          "",
		SourcePkg:           "",
		Repository:          "",
		Homepage:            "",
		Arch:                "",
		Kind:                "",
		Tags:                []string{},
		Version:             common.Version{},
		License:             "",
		RuntimeDependencies: []common.Dependency{},
		RuntimeSuggestions:  []common.Dependency{},
		BuildDependencies:   []common.Dependency{},
	}

	common.SetReadableVersion(&template.Version)

	template_json, err := json.MarshalIndent(template, "", " ")
	common.FailOnError(err, "Failed on serializing template to json string")

	err = ioutil.WriteFile(out_path+"template.json", template_json, 0644)
	common.FailOnError(err, "Failed on writing to file")
}

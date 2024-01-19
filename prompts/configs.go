package prompts

import (
	"fmt"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/pkg/errors"
	"github.com/speakeasy-api/openapi-generation/v2/pkg/generate"
	config "github.com/speakeasy-api/sdk-gen-config"
	"github.com/speakeasy-api/sdk-gen-config/workflow"
	"github.com/speakeasy-api/speakeasy/charm"
)

func PromptForTargetConfig(targetName string, target *workflow.Target) (*config.Configuration, error) {
	output, err := config.GetDefaultConfig(true, generate.GetLanguageConfigDefaults, map[string]bool{target.Target: true})
	if err != nil {
		return nil, errors.Wrapf(err, "error generating config for target %s of type %s", targetName, target.Target)
	}

	var sdkClassName string
	configFields := []huh.Field{
		huh.NewInput().
			Title("Name your SDK object:").
			Placeholder("your users will access SDK methods with <sdk_name>.doThing()").
			Inline(true).
			Prompt(" ").
			Value(&sdkClassName),
	}
	languageForms, err := languageSpecificForms(target.Target)
	if err != nil {
		return nil, err
	}

	configFields = append(configFields, languageForms...)
	form := huh.NewForm(
		huh.NewGroup(
			configFields...,
		))
	if _, err := tea.NewProgram(charm.NewForm(form,
		fmt.Sprintf("Let's configure your %s target (%s)", target.Target, targetName),
		"This will create a gen.yaml config file that defines parameters for how your SDK is generated. \n"+
			"We will go through a few basic configurations here, but you can always modify this file directly in the future.")).
		Run(); err != nil {
		return nil, err
	}

	output.Generation.SDKClassName = sdkClassName

	saveLanguageConfigValues(target.Target, form, output)

	return output, nil
}

func configBaseForm(quickstart *Quickstart) (*QuickstartState, error) {
	for key, target := range quickstart.WorkflowFile.Targets {
		output, err := PromptForTargetConfig(key, &target)
		if err != nil {
			return nil, err
		}

		quickstart.LanguageConfigs[key] = output
	}

	var nextState QuickstartState = GithubWorkflowBase
	return &nextState, nil
}

type configPrompt struct {
	Key    string
	Prompt string
}

var languageSpecificPrompts = map[string][]configPrompt{
	"go": {
		{
			Key:    "packageName",
			Prompt: "Choose a go module package name:",
		},
	},
	"typescript": {
		{
			Key:    "packageName",
			Prompt: "Choose a npm package name:",
		},
		{
			Key:    "author",
			Prompt: "Choose an author of the published package:",
		},
	},
	"python": {
		{
			Key:    "packageName",
			Prompt: "Choose a PyPI package name:",
		},
		{
			Key:    "author",
			Prompt: "Choose an author of the published package:",
		},
	},
	"java": {
		{
			Key:    "projectName",
			Prompt: "Choose a Gradle rootProject.name, which gives a name to the Gradle build:",
		},
		{
			Key:    "groupID",
			Prompt: "Choose a groupID to use for namespacing the package. This is usually the reversed domain name of your organization:",
		},
		{
			Key:    "artifactID",
			Prompt: "Choose a artifactID to use for namespacing the package. This is usually the name of your project:",
		},
	},
	"terraform": {
		{
			Key:    "packageName",
			Prompt: "Choose a terraform provider package name:",
		},
		{
			Key:    "author",
			Prompt: "Choose an author of the published provider:",
		},
	},
	"docs": {
		{
			Key:    "defaultLanguage",
			Prompt: "Choose a default language for your doc site:",
		},
	},
}

func languageSpecificForms(language string) ([]huh.Field, error) {
	t, err := generate.GetTargetFromTargetString(language)
	if err != nil {
		return nil, err
	}

	configFields, err := generate.GetLanguageConfigFields(t, true)
	if err != nil {
		return nil, err
	}

	fields := []huh.Field{}
	if prompts, ok := languageSpecificPrompts[language]; ok {
		for _, prompt := range prompts {
			if exists, defaultValue, validateRegex, validateMessage := getValuesForFieldName(configFields, prompt.Key); exists {
				fields = append(fields, addPromptForField(prompt.Key, prompt.Prompt, defaultValue, validateRegex, validateMessage)...)
			}
		}
	}

	return fields, nil
}

func getValuesForFieldName(configFields []config.SDKGenConfigField, fieldName string) (bool, string, string, string) {
	var packageNameConfig *config.SDKGenConfigField
	for _, field := range configFields {
		if field.Name == fieldName {
			packageNameConfig = &field
			break
		}
	}
	if packageNameConfig == nil {
		return false, "", "", ""
	}

	defaultValue := ""
	if packageNameConfig.DefaultValue != nil {
		defaultValue, _ = (*packageNameConfig.DefaultValue).(string)
	}

	validationRegex := ""
	if packageNameConfig.ValidationRegex != nil {
		validationRegex = *packageNameConfig.ValidationRegex
		validationRegex = strings.Replace(validationRegex, `\u002f`, `/`, -1)
	}

	validationMessage := ""
	if packageNameConfig.ValidationRegex != nil {
		validationMessage = *packageNameConfig.ValidationMessage
	}

	return true, defaultValue, validationRegex, validationMessage
}

func addPromptForField(key, question, defaultValue, validateRegex, validateMessage string) []huh.Field {
	return []huh.Field{
		huh.NewInput().
			Key(key).
			Title(question).
			Placeholder(defaultValue).
			Inline(true).
			Validate(func(s string) error {
				if validateRegex != "" {
					r, err := regexp.Compile(validateRegex)
					if err != nil {
						return err
					}
					if !r.MatchString(s) {
						return errors.New(validateMessage)
					}
				}
				return nil
			}).
			Prompt(" "),
	}
}

func saveLanguageConfigValues(language string, form *huh.Form, configuration *config.Configuration) {
	if prompts, ok := languageSpecificPrompts[language]; ok {
		for _, prompt := range prompts {
			configuration.Languages[language].Cfg[prompt.Key] = form.GetString(prompt.Key)
		}
	}
}

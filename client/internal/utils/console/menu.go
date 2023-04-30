package console

import (
	"github.com/AlecAivazis/survey/v2"
)

func AskPrompt(question string, options []string) (string, error) {
	prompt := &survey.Select{
		Message: question,
		Options: options,
	}
	var answer string
	err := survey.AskOne(prompt, &answer, survey.WithValidator(survey.Required))
	return answer, err
}

func Ask(question string) (string, error) {
	prompt := &survey.Input{
		Message: question,
	}
	var answer string
	err := survey.AskOne(prompt, &answer, survey.WithValidator(survey.Required))
	return answer, err
}

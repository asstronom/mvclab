package consoleview

import "github.com/manifoldco/promptui"

var (
	defaultSelect = &promptui.SelectTemplates{
		Active:   "{{ . | cyan }} \U0001F336",
		Inactive: "{{ . | white }}",
	}
	handlerSelect = &promptui.SelectTemplates{
		Active:   "{{ .Name | cyan }} \U0001F336",
		Inactive: "{{ .Name | white }}",
	}
)

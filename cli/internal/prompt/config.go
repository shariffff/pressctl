package prompt

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

// PromptEditorSelection prompts the user to select their preferred editor
func PromptEditorSelection() (string, error) {
	editors := []string{"nano", "vim", "code", "Other"}
	opts := make([]huh.Option[string], len(editors))
	for i, e := range editors {
		opts[i] = huh.NewOption(e, e)
	}

	var selected string
	if err := huh.NewSelect[string]().
		Title("Preferred editor").
		Description("This will be saved for future use").
		Options(opts...).
		Value(&selected).
		Run(); err != nil {
		return "", normalizeErr(err)
	}

	if selected == "Other" {
		var custom string
		if err := huh.NewInput().
			Title("Editor command").
			Description("e.g., subl, emacs, micro").
			Value(&custom).
			Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf("editor command is required")
				}
				return nil
			}).
			Run(); err != nil {
			return "", normalizeErr(err)
		}
		return custom, nil
	}

	return selected, nil
}

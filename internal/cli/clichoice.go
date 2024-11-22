package cli

type CliChoice struct {
	Name        string
	Description string
	Default     string
	Choices     []string
}

func (c CliChoice) IsValidChoice(choice string) bool {
	for _, c := range c.Choices {
		if c == choice {
			return true
		}
	}
	return false
}

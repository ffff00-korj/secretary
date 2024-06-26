package utils

import (
	"fmt"

	"github.com/ffff00-korj/secretary/internal/config"
)

func HelpMessage() string {
	return fmt.Sprintf(`Here's what I can do:
/%s to start application,
/%s to see help message,
/%s <name> <sum> <payment day> to add,
/%s to see how many dollars you spent on your next sallary.`, config.CmdStart, config.CmdHelp, config.CmdAdd, config.CmdExpenseReport)
}

func TextToMarkdown(text string) string {
	return fmt.Sprintf("```%s```", text)
}

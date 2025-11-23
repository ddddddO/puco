package model

import (
	"errors"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddddddO/puco/internal"
	"github.com/ddddddO/puco/internal/command"
	"github.com/ddddddO/puco/internal/phpunitxml"
)

type yesnoView struct {
	cursor   int
	choices  []string
	selected string

	cmdPHPUnit *command.CmdPHPUnit
}

func newYesNoView(commandToSpecifyBeforePHPCommand string) (*yesnoView, error) {
	return &yesnoView{
		choices: []string{"Yes", "No"},
		cmdPHPUnit: &command.CmdPHPUnit{
			CommandToSpecifyBeforePHPCommand: commandToSpecifyBeforePHPCommand,
		},
	}, nil
}

func (v *yesnoView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyOfQuiet1, KeyOfQuiet2:
			return m, tea.Quit

		case "enter":
			v.selected = v.choices[v.cursor]

			if v.selected == "Yes" {
				if err := v.prepare(m); err != nil {
					m.currentView = ViewOfError
					m.err = err
					return m, nil
				}

				return m, v.cmdPHPUnit.Command()
			}

			m.quitting = true
			return m, tea.Quit

		case "down", "j":
			v.cursor++
			if v.cursor >= len(v.choices) {
				v.cursor = 0
			}

		case "up", "k":
			v.cursor--
			if v.cursor < 0 {
				v.cursor = len(v.choices) - 1
			}
		}
	}

	return m, nil
}

func (v *yesnoView) prepare(m model) error {
	targetTests := []string{}
	for s := range m.selectTestFilesView.selected {
		targetTests = append(targetTests, s)
	}
	targetCoverages := []string{}
	for s := range m.selectCoverageFilesView.selected {
		targetCoverages = append(targetCoverages, s)
	}
	cfg := internal.Config{
		CommandToSpecifyBeforePHPCommand: v.cmdPHPUnit.CommandToSpecifyBeforePHPCommand,
		LatestExecutedData: struct {
			SelectedTestFilePaths       []string
			SelectedCoverageTargetPaths []string
		}{
			SelectedTestFilePaths:       targetTests,
			SelectedCoverageTargetPaths: targetCoverages,
		},
	}

	if err := internal.StoreConfig(cfg); err != nil {
		return err
	}

	if err := os.RemoveAll(command.OutputCoverageDir); err != nil {
		return err
	}
	if _, err := os.Stat(phpunitxml.OutputPHPUnitXML); err == nil {
		if e := os.Remove(phpunitxml.OutputPHPUnitXML); e != nil {
			return err
		}
	} else {
		if !os.IsNotExist(err) {
			return err
		}
	}

	// 微妙
	return phpunitxml.Generate(v.cmdPHPUnit.CommandToSpecifyBeforePHPCommand, targetTests, m.selectCoverageFilesView.longestMatchDirPath())
}

func (v *yesnoView) view(width int, cfv *selectCoverageFilesView) string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("%s\n\n", internal.ColorLightPinkStyle.Render("Execute PHPUnit?")))

	v.cmdPHPUnit.Build(cfv.longestMatchDirPath(), "PUCO", phpunitxml.OutputPHPUnitXML)

	rawCmd := v.cmdPHPUnit.RawCmd()
	if len(rawCmd) <= width {
		s.WriteString(fmt.Sprintf("%s\n\n", internal.ColorBrightGreenStyle.Render(rawCmd)))
	} else {
		// ターミナルの横幅より長いコマンドを改行して表示するため
		splited, err := splitStringByN(rawCmd, width)
		if err != nil {
			return temporaryErrorView(err, width)
		}
		for i := range splited {
			s.WriteString(fmt.Sprintf("%s\n", internal.ColorBrightGreenStyle.Render(splited[i])))
		}
		s.WriteString("\n")
	}

	for i := 0; i < len(v.choices); i++ {
		if v.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(v.choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(↑↓: move, Enter: next, Esc: quit)\n")

	return s.String()
}

func splitStringByN(s string, n int) ([]string, error) {
	if n <= 0 {
		return nil, errors.New("require n is upper 0")
	}

	runes := []rune(s)
	var chunks []string
	for i := 0; i < len(runes); i += n {
		start := i
		end := min(i+n, len(runes))

		chunk := string(runes[start:end])
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

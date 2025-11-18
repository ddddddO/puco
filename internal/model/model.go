package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddddddO/puco/internal"
	"github.com/ddddddO/puco/internal/command"
)

const (
	ViewOfSelectTestFiles = iota
	ViewOfSelectCoverageFiles
	ViewOfYesNo
	ViewOfCoverageList
)

var (
	KeyOfQuiet1 = tea.KeyCtrlC.String()
	KeyOfQuiet2 = tea.KeyEsc.String()
)

type model struct {
	cfg *internal.Config

	height int
	width  int

	currentView int
	quitting    bool

	selectTestFilesView     *selectTestFilesView
	selectCoverageFilesView *selectCoverageFilesView
	yesnoView               *yesnoView
	coverageListView        *coveragedListView

	err error
}

func New(cfg internal.Config, shouldRestoreLatestExecutedData bool) (model, error) {
	tfv, err := newSelectTestFilesView(cfg, shouldRestoreLatestExecutedData)
	if err != nil {
		return model{}, err
	}

	cfv, err := newSelectCoverageFilesView(cfg, shouldRestoreLatestExecutedData)
	if err != nil {
		return model{}, err
	}

	ynv, err := newYesNoView(cfg.CommandToSpecifyBeforePHPCommand)
	if err != nil {
		return model{}, err
	}

	clv := newCoverageListView()

	return model{
		currentView: ViewOfSelectTestFiles,

		selectTestFilesView:     tfv,
		selectCoverageFilesView: cfv,
		yesnoView:               ynv,
		coverageListView:        clv,
	}, nil
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case KeyOfQuiet1, KeyOfQuiet2:
			m.quitting = true
			return m, tea.Quit
		}
	case command.PHPUitFinishedMsg:
		if msg.Err() != nil {
			m.err = msg.Err()
			m.quitting = true
			return m, tea.Quit
		}

		m.currentView = ViewOfCoverageList
		return m, nil

	}

	switch m.currentView {
	case ViewOfSelectTestFiles:
		return m.selectTestFilesView.update(msg, m)
	case ViewOfSelectCoverageFiles:
		return m.selectCoverageFilesView.update(msg, m)
	case ViewOfYesNo:
		return m.yesnoView.update(msg, m)
	case ViewOfCoverageList:
		return m.coverageListView.update(msg, m)
	default:
		return m, nil
	}
}

func (m model) View() string {
	// 最終結果出力
	if m.quitting {
		// var sb strings.Builder
		// sb.WriteString("Result:\n\n")

		// sb.WriteString("Selected test files:\n")
		// if len(m.selectTestFilesView.selected) == 0 {
		// 	sb.WriteString("  (no selected))\n")
		// } else {
		// 	for choice := range m.selectTestFilesView.selected {
		// 		sb.WriteString(fmt.Sprintf("  - %s\n", choice))
		// 	}
		// }

		// sb.WriteString("\nSelected coverage target:\n")
		// if len(m.selectCoverageFilesView.selected) == 0 {
		// 	sb.WriteString("  (no selected)\n")
		// } else {
		// 	for choice := range m.selectCoverageFilesView.selected {
		// 		sb.WriteString(fmt.Sprintf("  - %s\n", choice))
		// 	}
		// }

		// return sb.String()
		return "end"
	}

	if m.err != nil {
		return fmt.Sprintf("failed...: \n%v\n", m.err)
	}

	switch m.currentView {
	case ViewOfSelectTestFiles:
		return m.selectTestFilesView.view(m.height)
	case ViewOfSelectCoverageFiles:
		return m.selectCoverageFilesView.view(m.height)
	case ViewOfYesNo:
		return m.yesnoView.view(m.width, m.selectCoverageFilesView)
	case ViewOfCoverageList:
		return m.coverageListView.view(m.height)
	default:
		return "unknown view"
	}
}

package model

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ddddddO/gtree"
	"github.com/ddddddO/puco/internal"
	"github.com/ddddddO/puco/internal/command"
)

type coveragedListView struct {
	ready    bool
	content  string
	viewport viewport.Model
}

func newCoverageListView() *coveragedListView {
	return &coveragedListView{}
}

func (c *coveragedListView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyOfQuiet1, KeyOfQuiet2:
			return m, tea.Quit

		case "enter":
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(c.headerView())
		footerHeight := lipgloss.Height(c.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !c.ready {
			c.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			c.viewport.YPosition = headerHeight

			c.viewport.SetContent(c.content)
			c.ready = true
		} else {
			c.viewport.Width = msg.Width
			c.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	c.viewport, cmd = c.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (c *coveragedListView) view(height, width int, phpCommandResult string) string {
	if !c.ready {
		cn, err := c.getContent()
		if err != nil {
			return temporaryErrorView(err, width)
		}
		escMsg := "(Enter: quit)"
		resultTitle := fmt.Sprintf("\n\n%s\n\n", internal.ColorLightPinkStyle.Render("===== Result of PHPUnit ====="))
		c.content = escMsg + resultTitle + fmt.Sprintf("%s%s\n", phpCommandResult, cn)

		headerHeight := lipgloss.Height(c.headerView())
		footerHeight := lipgloss.Height(c.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		c.viewport = viewport.New(width, height-verticalMarginHeight)
		c.viewport.YPosition = headerHeight

		c.viewport.SetContent(c.content)
		c.ready = true
	}
	return fmt.Sprintf("%s\n%s\n%s", c.headerView(), c.viewport.View(), c.footerView())
}

func (c *coveragedListView) getContent() (string, error) {
	s := &strings.Builder{}
	lvl := 2
	s.WriteString(fmt.Sprintf("\n\n%s\n\n", internal.ColorLightPinkStyle.Render(fmt.Sprintf("===== Cveraged file list (Max depth: %d) =====", lvl+1))))

	coverages, err := internal.GetCoveragedFilePaths(lvl)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			e := fmt.Errorf("'%s' directory did not exist, probably because pcov was not installed", command.OutputCoverageDir)
			return "", e
		}

		return "", err
	}

	root := gtree.NewRoot(command.OutputCoverageDir)
	var node *gtree.Node
	for i := range coverages {
		for j, name := range strings.Split(coverages[i], string(filepath.Separator)) {
			if j == 0 {
				node = root
				continue
			}

			node = node.Add(name)
		}
	}
	for iter, err := range gtree.WalkIterFromRoot(root) {
		if err != nil {
			return "", err
		}

		if iter.Level() == 1 {
			s.WriteString(fmt.Sprintf("%s\n", iter.Row()))
			continue
		}
		s.WriteString(fmt.Sprintf("  %s\n", iter.Row()))
	}
	s.WriteString("\n(Enter: quit)\n")

	return s.String(), nil
}

func (c *coveragedListView) headerView() string {
	title := ""
	line := strings.Repeat("─", max(0, c.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (c *coveragedListView) footerView() string {
	info := fmt.Sprintf("%3.f%%", c.viewport.ScrollPercent()*100)
	line := strings.Repeat("─", max(0, c.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

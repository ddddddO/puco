package model

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddddddO/gtree"
	"github.com/ddddddO/puco/internal"
	"github.com/ddddddO/puco/internal/command"
)

type coveragedListView struct {
}

func newCoverageListView() *coveragedListView {
	return &coveragedListView{}
}

func (c *coveragedListView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyOfQuiet1, KeyOfQuiet2:
			return m, tea.Quit

		case "enter":
			m.quitting = true
			return m, tea.Quit

			// case "down", "j":
			// 	v.cursor++
			// 	if v.cursor >= len(v.choices) {
			// 		v.cursor = 0
			// 	}

			// case "up", "k":
			// 	v.cursor--
			// 	if v.cursor < 0 {
			// 		v.cursor = len(v.choices) - 1
			// 	}
		}
	}

	return m, tea.Quit
}

func (c *coveragedListView) view(viewHeight int) string {
	s := &strings.Builder{}
	lvl := 2
	s.WriteString(fmt.Sprintf("\n\n%s\n\n", internal.ColorLightPinkStyle.Render(fmt.Sprintf("===== Cveraged file list (Max depth: %d) =====", lvl+1))))

	coverages, err := internal.GetCoveragedFilePaths(lvl)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			panic(fmt.Errorf("%s directory did not exist, probably because pcov was not installed.", command.OutputCoverageDir))
		}

		panic(err)
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
	rowCnt := 0
	height := max(0, viewHeight-15) // 数は一旦決め打ち。phpunit実行時の出力が多いと表示崩れちゃいそうだけど...
	for iter, err := range gtree.WalkIterFromRoot(root) {
		if err != nil {
			panic(err)
		}

		if rowCnt > height {
			break
		}

		rowCnt++
		if iter.Level() == 1 {
			s.WriteString(fmt.Sprintf("%s\n", iter.Row()))
			continue
		}
		s.WriteString(fmt.Sprintf("  %s\n", iter.Row()))
	}
	if rowCnt > viewHeight {
		// TODO: ↑の決め打ちの数次第か表示されてない
		s.WriteString(fmt.Sprintln("  ... more"))
	}
	s.WriteString("\n")

	// TODO: なんとかしたい
	s.WriteString("If it is not finished, press any key to finish it...")

	return s.String()
}

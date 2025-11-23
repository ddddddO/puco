package model

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddddddO/puco/internal"
)

type errorView struct{}

func newErrorView() *errorView {
	return &errorView{}
}

// 各modelのviewメソッド内で呼ぶための
func temporaryErrorView(err error, width int) string {
	return newErrorView().view(err, width)
}

func (e *errorView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyOfQuiet1, KeyOfQuiet2:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (e *errorView) view(err error, width int) string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("%s\n\n", internal.ColorBrightRedBoldStyle.Render("Ooops... Error has occurred.")))

	errMsg := err.Error()
	if len(errMsg) <= width {
		s.WriteString(fmt.Sprintf("%s\n\n", errMsg))
	} else {
		// ターミナルの横幅より長いコマンドを改行して表示するため
		splited, err := splitStringByN(errMsg, width)
		if err != nil {
			// panic(err)
			s.WriteString(fmt.Sprintf("%s\n\n", errMsg))
			s.WriteString("(Esc: quit)\n")
			return s.String()
		}
		for i := range splited {
			s.WriteString(fmt.Sprintf("%s\n", splited[i]))
		}
		s.WriteString("\n")
	}

	s.WriteString("(Esc: quit)\n")
	return s.String()
}

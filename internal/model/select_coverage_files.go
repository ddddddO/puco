package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddddddO/puco/internal"
	"github.com/sahilm/fuzzy"
)

type selectCoverageFilesView struct {
	height int

	choices         []string      // 選択肢のリスト
	filteredChoices []fuzzy.Match // 絞り込まれた選択肢のリスト

	cursor      int
	searchInput textinput.Model

	selected map[string]struct{}
}

func newSelectCoverageFilesView(cfg internal.Config, shouldRestoreLatestExecutedData bool) (*selectCoverageFilesView, error) {
	paths, err := internal.GetPHPCodeFilePaths()
	if err != nil {
		return nil, err
	}

	ti := textinput.New()
	ti.Placeholder = "Filter coverage target files..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	selected := make(map[string]struct{})
	if shouldRestoreLatestExecutedData {
		for i := range paths {
			if cfg.IsMatchedCoverageTargetFile(paths[i]) {
				selected[paths[i]] = struct{}{}
			}
		}
	}

	return &selectCoverageFilesView{
		choices:         paths,
		filteredChoices: fuzzy.Find("", paths),
		selected:        selected,
		searchInput:     ti,
	}, nil
}

func (t *selectCoverageFilesView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if t.cursor > 0 {
				t.cursor--
			}

			t.searchInput, cmd = t.searchInput.Update(msg)
			return m, cmd

		case "down", "j":
			// 検索が空っぽの時
			if len(strings.TrimSpace(t.searchInput.Value())) == 0 {
				if t.cursor < len(t.choices)-1 {
					t.cursor++
				}
			}

			if t.cursor < len(t.filteredChoices)-1 {
				t.cursor++
			}

			t.searchInput, cmd = t.searchInput.Update(msg)
			return m, cmd

		// spaceで選択・選択解除
		case " ":
			// 検索が空っぽの時
			if len(strings.TrimSpace(t.searchInput.Value())) == 0 {
				choice := t.choices[t.cursor]
				_, ok := t.selected[choice]
				if ok {
					delete(t.selected, choice)
				} else {
					t.selected[choice] = struct{}{}
				}
			}

			if len(t.filteredChoices) > 0 {
				choice := t.filteredChoices[t.cursor].Str
				if _, ok := t.selected[choice]; ok {
					delete(t.selected, choice)
				} else {
					t.selected[choice] = struct{}{}
				}
			}

			// スペースで選択するため、検索窓への入力でスペースは許容しない
			// t.searchInput, cmd = t.searchInput.Update(msg)
			return m, nil

		case "enter":
			m.currentView = ViewOfYesNo
			return m, nil

		default:
			t.searchInput, cmd = t.searchInput.Update(msg)
			t.filteredChoices = fuzzy.Find(t.searchInput.Value(), t.choices)
			t.cursor = 0
			return m, cmd
		}
	}

	t.searchInput, cmd = t.searchInput.Update(msg)

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

func (t *selectCoverageFilesView) view(viewHeight int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\n\n", internal.ColorLightPinkStyle.Render("Select target coverage files (press Space)")))
	sb.WriteString(t.searchInput.View())
	sb.WriteString("\n\n")

	// t.filteredChoices は、fuzzy.Findの第一引数が空文字だとマッチしない。検索文字列が空なら全権表示するようにするため以下の処理
	matchs := []string{}
	if len(strings.TrimSpace(t.searchInput.Value())) == 0 {
		matchs = t.choices
	} else {
		for _, m := range t.filteredChoices {
			matchs = append(matchs, m.Str)
		}
	}

	// マイナスしてるのは、パス一覧を除いた高さを一旦決め打ちした数
	height := min(len(matchs), viewHeight-7)
	height = max(0, height) // 起動時、heightがマイナス値になることあってパニックになるから

	for i, match := range matchs[:height] {
		coloredMatch := match
		cursor := " " // no cursor
		if t.cursor == i {
			cursor = ">" // cursor!
			coloredMatch = internal.ColorBrightGreenStyle.Render(match)

		}

		checked := " " // not selected
		if _, ok := t.selected[match]; ok {
			checked = "x" // selected!
			coloredMatch = internal.ColorBrightBlueStyle.Render(match)

		}

		// Render the row
		sb.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, coloredMatch))
	}
	if len(matchs) > height {
		sb.WriteString(fmt.Sprintf("  ... %d more", len(matchs)-height))
	}

	// The footer
	sb.WriteString("\n(↑↓: move, Space: select, Enter: next, Esc: quit)\n")

	// Send the UI for rendering
	return sb.String()
}

// 選択されたカバレッジ取りたいPHPファイルのパスの最大公約数的なディレクトリのパスを返す
func (t *selectCoverageFilesView) longestMatchDirPath() string {
	if len(t.selected) == 0 {
		return "./"
	}

	segments := make([][]string, len(t.selected))
	i := 0
	for p := range t.selected {
		dirPath := filepath.Dir(p)
		cleanedPath := filepath.Clean(dirPath)
		segments[i] = strings.Split(cleanedPath, string(filepath.Separator))
		i++
	}
	minLen := len(segments[0])
	for i := 1; i < len(segments); i++ {
		if len(segments[i]) < minLen {
			minLen = len(segments[i])
		}
	}

	commonSegments := []string{}
	firstSegments := segments[0]
	for i := 0; i < minLen; i++ {
		currentSegment := firstSegments[i]
		isCommon := true

		for j := 1; j < len(segments); j++ {
			if segments[j][i] != currentSegment {
				isCommon = false
				break
			}
		}

		if isCommon {
			commonSegments = append(commonSegments, currentSegment)
		} else {
			// 一致しない要素が現れたら、そこで終了
			break
		}
	}

	if len(commonSegments) == 0 {
		return "./"
	}

	return filepath.Join(commonSegments...)
}

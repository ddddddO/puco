package internal

import "github.com/charmbracelet/lipgloss"

// 各Viewのヘッダ文言用
var ColorLightPinkStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFB6C1")).
	Bold(false)

// list移動時のカーソルがさしてる時の色
var ColorBrightGreenStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("10")).
	Bold(true)

// listにチェック済みの色
var ColorBrightBlueStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("12")).
	Bold(true)

// error字のタイトル色
var ColorBrightRedBoldStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("9")).
	Bold(true)

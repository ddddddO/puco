package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type PHPUitFinishedMsg struct {
	Stdout string
	Stderr string
	err    error
}

func (p PHPUitFinishedMsg) Err() error {
	return p.err
}

type CmdPHPUnit struct {
	CommandToSpecifyBeforePHPCommand string
	cmd                              *exec.Cmd
}

const OutputCoverageDir = "coverage-puco"

// TODO: test/coverage以外の、-d memory_limit指定は、設定ファイルから受け取るようにする
func (c *CmdPHPUnit) Build(targetCoverageDir string, testSuiteName string, configFile string) {
	if len(c.CommandToSpecifyBeforePHPCommand) == 0 {
		c.cmd = exec.Command("php", []string{
			"-d",
			fmt.Sprintf("pcov.directory=%s", targetCoverageDir), // ここカンマ区切りで複数指定可能のようだけど、どうもうまくいってない。なので、最大公約数的なパスを一旦指定しておく
			"-d",
			fmt.Sprintf("memory_limit=%s", "512M"),
			"vendor/bin/phpunit",
			"--testsuite",
			testSuiteName,
			"--configuration",
			configFile,
			"--coverage-html",
			OutputCoverageDir,
		}...)

		return
	}

	parsedCmd := strings.Split(c.CommandToSpecifyBeforePHPCommand, " ")
	if len(parsedCmd) == 1 {
		c.cmd = exec.Command(parsedCmd[0], []string{
			"php",
			"-d",
			fmt.Sprintf("pcov.directory=%s", targetCoverageDir),
			"-d",
			fmt.Sprintf("memory_limit=%s", "512M"),
			"vendor/bin/phpunit",
			"--testsuite",
			testSuiteName,
			"--configuration",
			configFile,
			"--coverage-html",
			OutputCoverageDir,
		}...)

		return
	}

	args := append(parsedCmd[1:], []string{
		"php",
		"-d",
		fmt.Sprintf("pcov.directory=%s", targetCoverageDir),
		"-d",
		fmt.Sprintf("memory_limit=%s", "512M"),
		"vendor/bin/phpunit",
		"--testsuite",
		testSuiteName,
		"--configuration",
		configFile,
		"--coverage-html",
		OutputCoverageDir,
	}...)
	c.cmd = exec.Command(parsedCmd[0], args...)
}

func (c *CmdPHPUnit) RawCmd() string {
	return c.cmd.String()
}

func (c *CmdPHPUnit) Command() tea.Cmd {
	var stdout, stderr bytes.Buffer
	// リアルタイムでも出力したいのでこうする
	c.cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	c.cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	return tea.ExecProcess(c.cmd, func(err error) tea.Msg {
		return PHPUitFinishedMsg{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
			err:    err,
		}
	})
}

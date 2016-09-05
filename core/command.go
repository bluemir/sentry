package core

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type shellCommander struct {
	shell string
	proc  *os.Process
}

func newShellCommander(shell string) *shellCommander {
	return &shellCommander{shell, nil}
}

func (c *shellCommander) exec(command string) {
	cmd := exec.Command(c.shell, "-c", command)
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", out.String())

	c.proc = cmd.Process
}
func (c *shellCommander) stop() {
	if c.proc == nil {
		return
	}

	err := c.proc.Kill()
	if err != nil {
		log.Info(err)
	}
	c = nil
}

package core

import (
	"bufio"
	"io"
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

	var reader io.Reader
	reader, cmd.Stdout = io.Pipe()
	//cmd.Stdout = stdout
	go func() {
		r := bufio.NewReader(reader)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				log.Debugln("[end of stdout]")
				return
			}
			log.Info(string(line))
		}
		log.Debugln("process read line terminated")
	}()

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	c.proc = cmd.Process
}
func (c *shellCommander) stop() {
	if c.proc == nil {
		return
	}
	log.Infof("kill process")
	err := c.proc.Kill()
	if err != nil {
		log.Info(err)
	}
	c = nil
}

package core

import (
	"bufio"
	"io"
	"os"
	"os/exec"

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
	//cmd.Stdin = strings.NewReader("some input")

	link(log.Info, cmd.StdoutPipe)
	link(log.Warn, cmd.StderrPipe)

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
		log.Warnln(err)
	}

	c = nil
}

func link(out func(...interface{}), in func() (io.ReadCloser, error)) {

	reader, err := in()
	if err != nil {
		log.Warn(err)
	}
	go func() {

		r := bufio.NewReader(reader)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				log.Debugln("[end of stdout]")
				return
			}
			out(string(line))
		}
		log.Debugln("process read line terminated")
	}()
}

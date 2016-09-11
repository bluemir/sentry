package core

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"syscall"

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

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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
	log.Infof("try to kill process...")
	pgid, err := syscall.Getpgid(c.proc.Pid)
	if err != nil {
		log.Warnln(err)
		return
	}
	if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
		log.Warnln(err)
		return
	}

	log.Warn("kill process : ", c.proc.Pid)
	c.proc = nil
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

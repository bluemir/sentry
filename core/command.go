package core

import (
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

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
	log.Infof("Try to kill process...")
	pgid, err := syscall.Getpgid(c.proc.Pid)
	if err != nil {
		log.Warnln(err)
		return
	}
	if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
		log.Warnln(err)
		return
	}

	log.Warnf("Process killed. PID: %d", c.proc.Pid)
	c.proc = nil
}

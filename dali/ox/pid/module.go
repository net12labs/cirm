package pid

import (
	"fmt"
	"os"
	"syscall"
)

type PidModule struct {
	// Pid fields here
}

func (p *PidModule) NewPid() *Pid {
	return &Pid{}
}

type Pid struct {
	Pid         int
	PidFilePath string
}

func (p *Pid) IsProcessRunning() bool {
	process, err := os.FindProcess(p.Pid)
	if err != nil {
		return false
	}
	// Sending signal 0 to a process is a way to check for its existence
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func (p *Pid) PidFileExists() bool {
	if err := p.ensure_filePath(); err != nil {
		panic(err)
	}
	_, err := os.Stat(p.PidFilePath)
	return err == nil
}

func (p *Pid) KillProcess() error {
	process, err := os.FindProcess(p.Pid)
	if err != nil {
		return err
	}
	return process.Kill()
}

func (p *Pid) SetPidFromFile() error {

	data, err := os.ReadFile(p.PidFilePath)
	if err != nil {
		return err
	}

	var pid int
	_, err = fmt.Sscanf(string(data), "%d", &pid)
	if err != nil {
		return err
	}

	p.Pid = pid
	return nil
}

func (p *Pid) RemovePidFile() error {
	if err := p.ensure_filePath(); err != nil {
		return err
	}
	return os.Remove(p.PidFilePath)
}

func (p *Pid) SetOwnPid() {
	p.Pid = os.Getpid()
}

func (p *Pid) ensure_pid() error {
	if p.Pid == 0 {
		return fmt.Errorf("PID is not set")
	}
	return nil
}
func (p *Pid) ensure_filePath() error {
	if p.PidFilePath == "" {
		return fmt.Errorf("PID file path is not set")
	}
	return nil
}

func (p *Pid) WritePidFile() error {
	if err := p.ensure_pid(); err != nil {
		return err
	}
	if err := p.ensure_filePath(); err != nil {
		return err
	}

	err := os.WriteFile(p.PidFilePath, []byte(fmt.Sprintf("%d", p.Pid)), 0644)
	if err != nil {
		return err
	}

	return nil
}

var Module = PidModule{}

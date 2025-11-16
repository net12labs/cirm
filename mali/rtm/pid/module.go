package pid

import (
	"fmt"

	pid "github.com/net12labs/cirm/mali/rtm/pid/def"
)

var Pid = pid.Module

type PidHandler struct {
	Pid *pid.Pid
}

func (h *PidHandler) Init() *PidHandler {
	h.Pid = Pid.NewPid()
	h.Pid.SetOwnPid()
	return h
}

func (h *PidHandler) Handle_ExitOnDuplicate() error {

	if h.Pid.PidFileExists() {
		existingPid := Pid.NewPid()
		existingPid.PidFilePath = h.Pid.PidFilePath
		existingPid.SetPidFromFile()
		if existingPid.IsProcessRunning() {
			return fmt.Errorf("service already running")
		} else {
			// Stale PID file, remove it
			existingPid.RemovePidFile()
		}
	}
	h.Pid.WritePidFile()
	return nil
}

func (h *PidHandler) Handle_CleanupOnExit() error {
	return h.Pid.RemovePidFile()
}

func NewPidHandler() *PidHandler {
	handler := &PidHandler{}
	return handler.Init()
}

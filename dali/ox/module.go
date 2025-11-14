package os

import (
	etc "dali/ox/etc"
	pid "dali/ox/pid"
	"fmt"
)

var Pid = pid.Module
var Etc = etc.Module

type PidHandler struct {
	Pid *pid.Pid
}

func (h *PidHandler) Init() *PidHandler {
	h.Pid = Pid.NewPid()
	h.Pid.SetOwnPid()
	h.Pid.PidFilePath = Etc.Get("pid_dir").String() + "/" + Etc.Get("pid_file_name").String()
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

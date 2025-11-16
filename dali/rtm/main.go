package rtm

import (
	"github.com/net12labs/cirm/mali/rtm/args"
	"github.com/net12labs/cirm/mali/rtm/do"
	"github.com/net12labs/cirm/mali/rtm/etc"
	"github.com/net12labs/cirm/mali/rtm/pid"
	"github.com/net12labs/cirm/mali/rtm/unit"
)

var Runtime = unit.NewRuntimeUnit()
var Etc = etc.Module
var Pid = pid.NewPidHandler()
var Args = args.NewArgs()

var Do = do.NewDo()

func init() {
	Pid.Init()
}

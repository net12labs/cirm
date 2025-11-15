package rtm

import (
	"github.com/net12labs/cirm/dali/runtime/args"
	"github.com/net12labs/cirm/dali/runtime/do"
	"github.com/net12labs/cirm/dali/runtime/etc"
	"github.com/net12labs/cirm/dali/runtime/pid"
	"github.com/net12labs/cirm/dali/runtime/unit"
)

var Runtime = unit.NewRuntimeUnit()
var Etc = etc.Module
var Pid = pid.NewPidHandler()
var Args = args.NewArgs()

var Do = do.NewDo()

func init() {
	Pid.Init()
}

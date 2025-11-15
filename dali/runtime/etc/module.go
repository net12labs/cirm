package etc

import (
	mod "github.com/net12labs/cirm/dali/x-module"
)

var Module = mod.NewEtcStore()

func init() {
	Module.SetKV("pid_dir", "/var/run")
}

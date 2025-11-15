package etc

import (
	mod "github.com/net12labs/cirm/bin/etc-store"
)

var Module = mod.NewEtcStore()

func init() {
	Module.SetKV("pid_dir", "/var/run")
}

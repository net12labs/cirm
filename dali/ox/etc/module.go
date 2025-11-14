package etc

import (
	mod "github.com/net12labs/cirm/dali/x-module"
)

var Module = mod.NewEtcStoreCb(func(es *mod.EtcStore) {
	es.SetKV("pid_dir", "/var/run")
})

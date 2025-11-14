package etc

import (
	mod "dali/x-module"
)

var Module = mod.NewEtcStoreCb(func(es *mod.EtcStore) {
	es.SetKV("pid_dir", "/var/run")
})

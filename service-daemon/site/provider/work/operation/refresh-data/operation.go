package refreshdata

import (
	wop "github.com/net12labs/cirm/dali/work/operation"
	"github.com/net12labs/cirm/dali/work/task"
	fetchdata "github.com/net12labs/cirm/service-daemon/site/provider/work/task/fetch-asn-prefixes"
)

type RefreshData struct {
	Op wop.Operation
}

func (r *RefreshData) Init() {
	r.Op = *wop.NewOperation()
	r.Op.Name = "Refresh IP Data"
	r.Op.Execute = func() error {
		retrieveTask := fetchdata.FetchIpData{
			Task: task.Task{},
			OnStart: func() {
				// Task start logic here
			},
			OnDone: func() {
				// Task done logic here
			},
		}
		retrieveTask.Run()
		return nil
	}
}

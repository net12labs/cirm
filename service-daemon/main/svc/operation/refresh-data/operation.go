package refreshdata

import (
	wop "cirm/lib/work/operation"
	"cirm/lib/work/task"
	fetchdata "cirm/svc/task/fetch-asn-prefixes"
)

type RefreshData struct {
	Op wop.Operation
}

func (r *RefreshData) Init() {
	r.Op = wop.Operation{
		Name: "Refresh IP Data",
		Execute: func() error {
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
		},
	}
}

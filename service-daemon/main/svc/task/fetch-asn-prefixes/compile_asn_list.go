package fetchdata

import (
	"bufio"
	"cirm/lib/work/task"
	"os"
	"strconv"
	"strings"
)

type CompileAsnList struct {
	Task    task.Task // this actually should be a subtask
	OnStart func()
}

// extract ASN data
// but first we need to instantiate a task, build a plan and then start execution that ca be resumable or progress trackacble
func (r *CompileAsnList) loadASNsFromFile(filename string) ([]int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var asns []int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		asn, err := strconv.Atoi(line)
		if err != nil {
			continue
		}
		asns = append(asns, asn)
	}

	return asns, scanner.Err()
}

package fetchdata

import (
	"cirm/lib/work/task"
	"fmt"
	"os"
	"strconv"
)

var defaultASNList = []int{
	// China Telecom (CT)
	4134, 4809, 23724, 58466, 58519, 58536, 58540, 58541, 58542, 58543,
	// China Unicom (CU)
	4837, 9929, 17621, 17622, 17623,
	// China Mobile (CMCC)
	9808, 24400, 56040, 56041, 56042, 56044, 56046, 56047, 56048,
	// CERNET (China Education and Research Network)
	4538, 23910,
	// CSTNET (Chinese Science and Technology Network)
	7497,
	// Dr.Peng Telecom & Media Group
	17816, 17897,
	// Tencent
	45090, 132203,
	// Alibaba
	37963, 45102,
	// Baidu
	38365, 55967,
	// Other major Chinese ISPs and networks
	4229, 4812, 4847, 9394, 17816, 18118,
	21859, 23764, 24407, 24429, 24444, 24482,
	38019, 38283, 38587, 38854,
	45062, 45071, 45075,
	55720, 58453, 58593, 58772, 58834, 59019,
	62610, 63525, 63593, 63629, 63655,
	132525, 133111, 133219, 133365, 133492, 134419, 134543, 134756, 134770,
	135377, 136190, 136958, 137266, 137539, 137687, 137710, 137718, 137726, 137735,
}

type FetchIpData struct {
	Task    task.Task
	Fetcher *RoutesFetcher
	OnStart func()
	OnDone  func()
	asnList []int
	AsnFile string
}

func (r *FetchIpData) Run() {
	success := 0
	failed := 0

	for _, asn := range r.asnList {
		count, err := r.Fetcher.fetchASN(asn)
		if err != nil {
			failed++
			fmt.Fprintf(os.Stderr, "✗ AS%d: Error: %v\n", asn, err)
		} else if count > 0 {
			success++
			fmt.Fprintf(os.Stderr, "✓ AS%d: %d routes\n", asn, count)
		} else {
			failed++
			fmt.Fprintf(os.Stderr, "✗ AS%d: No routes found\n", asn)
		}
	}
}

func (r *FetchIpData) compileAsnList() {
	// Determine ASN list
	var asns []int
	if len(r.asnList) > 0 {
		asns = r.asnList
	} else if *asnFile != "" {
		var err error
		asns, err = loadASNsFromFile(*asnFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading ASN file: %v\n", err)
			os.Exit(1)
		}
	} else {
		asns = defaultASNList
	}

}

type intSlice []int

func (i *intSlice) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *intSlice) Set(value string) error {
	val, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*i = append(*i, val)
	return nil
}

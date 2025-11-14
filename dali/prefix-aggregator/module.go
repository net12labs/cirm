package routeaggregator

import (
	"math"
	"math/big"
	"net/netip"
	"sort"
)

type PrefixAggregator struct {
	// PrefixAggregator fields here
	InputPrefixList        []string
	AggregateNonContigious bool
	OutputPrefixList       []string

	// Aggregation parameters
	MinPrefixLength    int     // Minimum prefix length to aggregate to (default: 24)
	MaxPrefixReduction int     // Maximum bits to reduce from original prefix (default: 16)
	WasteRatio2        float64 // Max waste ratio for 2-network groups (default: 100)
	WasteRatio3        float64 // Max waste ratio for 3-network groups (default: 500)
	WasteRatio5        float64 // Max waste ratio for 4-5 network groups (default: 2000)
	WasteRatio10       float64 // Max waste ratio for 6-10 network groups (default: 10000)
	WasteRatioMany     float64 // Max waste ratio for 10+ network groups (default: 100000)
}

func NewPrefixAggregator() *PrefixAggregator {
	return &PrefixAggregator{
		AggregateNonContigious: true,
		MinPrefixLength:        24,
		MaxPrefixReduction:     16,
		WasteRatio2:            100,
		WasteRatio3:            500,
		WasteRatio5:            2000,
		WasteRatio10:           10000,
		WasteRatioMany:         100000,
	}
}

// Aggregate performs prefix aggregation on InputPrefixList and stores result in OutputPrefixList
func (pa *PrefixAggregator) Aggregate() error {
	if len(pa.InputPrefixList) == 0 {
		pa.OutputPrefixList = []string{}
		return nil
	}

	// Parse input prefixes
	networks := make([]netip.Prefix, 0, len(pa.InputPrefixList))
	for _, cidr := range pa.InputPrefixList {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			continue // Skip invalid prefixes
		}
		networks = append(networks, prefix)
	}

	// Sort networks for consistent processing
	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Addr().Compare(networks[j].Addr()) < 0
	})

	// First pass: collapse overlapping/contiguous networks
	collapsed := collapseNetworks(networks)

	// If aggressive aggregation is enabled, find optimal supernets
	if pa.AggregateNonContigious {
		collapsed = pa.findOptimalSupernets(collapsed)
		// Collapse again to remove any overlaps created
		collapsed = collapseNetworks(collapsed)
	}

	// Convert back to strings
	pa.OutputPrefixList = make([]string, len(collapsed))
	for i, net := range collapsed {
		pa.OutputPrefixList[i] = net.String()
	}

	return nil
}

// collapseNetworks removes overlapping networks
func collapseNetworks(networks []netip.Prefix) []netip.Prefix {
	if len(networks) == 0 {
		return networks
	}

	sort.Slice(networks, func(i, j int) bool {
		if networks[i].Addr().Compare(networks[j].Addr()) == 0 {
			return networks[i].Bits() < networks[j].Bits()
		}
		return networks[i].Addr().Compare(networks[j].Addr()) < 0
	})

	result := []netip.Prefix{networks[0]}

	for i := 1; i < len(networks); i++ {
		current := networks[i]
		lastIdx := len(result) - 1
		last := result[lastIdx]

		// Check if current is a subnet of last
		if isSubnetOf(current, last) {
			continue // Skip, already covered
		}

		// Check if last is a subnet of current
		if isSubnetOf(last, current) {
			result[lastIdx] = current // Replace with broader network
			continue
		}

		result = append(result, current)
	}

	return result
}

// findOptimalSupernets finds optimal supernets covering groups of networks
func (pa *PrefixAggregator) findOptimalSupernets(networks []netip.Prefix) []netip.Prefix {
	if len(networks) == 0 {
		return networks
	}

	result := []netip.Prefix{}
	used := make(map[int]bool)

	for i := 0; i < len(networks); i++ {
		if used[i] {
			continue
		}

		current := networks[i]
		bestGroup := []int{i}
		bestSupernet := current
		bestScore := 0.0

		// Determine minimum acceptable prefix
		minPrefix := max(pa.MinPrefixLength, current.Bits()-pa.MaxPrefixReduction)

		// Try different prefix lengths
		for targetPrefix := current.Bits() - 1; targetPrefix >= minPrefix; targetPrefix-- {
			supernet := supernetAt(current, targetPrefix)
			if !supernet.IsValid() {
				break
			}

			// Find all networks covered by this supernet
			group := []int{i}
			for j := i + 1; j < len(networks); j++ {
				if !used[j] && isSubnetOf(networks[j], supernet) {
					group = append(group, j)
				}
			}

			// Calculate if this aggregation is worthwhile
			if len(group) >= 2 {
				totalCoverage := big.NewInt(0)
				for _, idx := range group {
					size := networkSize(networks[idx])
					totalCoverage.Add(totalCoverage, size)
				}

				supernetSize := networkSize(supernet)
				wasteRatio := float64(0)
				if totalCoverage.Sign() > 0 {
					wasteRatioFloat, _ := new(big.Float).Quo(
						new(big.Float).SetInt(supernetSize),
						new(big.Float).SetInt(totalCoverage),
					).Float64()
					wasteRatio = wasteRatioFloat
				} else {
					wasteRatio = math.Inf(1)
				}

				// Adjusted max waste based on group size
				adjustedMax := pa.getAdjustedMaxWaste(len(group))

				// Calculate efficiency score
				networksSaved := float64(len(group) - 1)
				bonus := math.Pow(float64(len(group)), 1.5)
				score := networksSaved * bonus / math.Max(1, math.Pow(wasteRatio, 0.3))

				if wasteRatio <= adjustedMax && score > bestScore {
					bestGroup = group
					bestSupernet = supernet
					bestScore = score
				}
			}
		}

		// Use the best aggregation found
		if len(bestGroup) > 1 {
			result = append(result, bestSupernet)
			for _, idx := range bestGroup {
				used[idx] = true
			}
		} else {
			result = append(result, current)
			used[i] = true
		}
	}

	return result
}

// isSubnetOf checks if prefix is a subnet of parent
func isSubnetOf(prefix, parent netip.Prefix) bool {
	if prefix.Bits() < parent.Bits() {
		return false
	}
	parentMasked := parent.Masked()
	prefixMasked := netip.PrefixFrom(prefix.Addr(), parent.Bits()).Masked()
	return parentMasked.Addr().Compare(prefixMasked.Addr()) == 0
}

// supernetAt creates a supernet at the specified prefix length
func supernetAt(prefix netip.Prefix, newBits int) netip.Prefix {
	if newBits >= prefix.Bits() {
		return netip.Prefix{}
	}
	return netip.PrefixFrom(prefix.Addr(), newBits).Masked()
}

// networkSize calculates the number of addresses in a network
func networkSize(prefix netip.Prefix) *big.Int {
	bits := prefix.Bits()
	var maxBits int
	if prefix.Addr().Is6() {
		maxBits = 128
	} else {
		maxBits = 32
	}
	hostBits := maxBits - bits
	size := new(big.Int).Lsh(big.NewInt(1), uint(hostBits))
	return size
}

// getAdjustedMaxWaste returns the maximum acceptable waste ratio based on group size
func (pa *PrefixAggregator) getAdjustedMaxWaste(groupSize int) float64 {
	switch {
	case groupSize == 2:
		return pa.WasteRatio2
	case groupSize == 3:
		return pa.WasteRatio3
	case groupSize <= 5:
		return pa.WasteRatio5
	case groupSize <= 10:
		return pa.WasteRatio10
	default:
		return pa.WasteRatioMany
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

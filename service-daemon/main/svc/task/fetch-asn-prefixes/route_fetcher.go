package fetchdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type RoutesFetcher struct {
	verbose bool
	routes  map[int][]string
	client  *http.Client
}

type RIPEResponse struct {
	Data struct {
		Prefixes []struct {
			Prefix string `json:"prefix"`
		} `json:"prefixes"`
	} `json:"data"`
}

type BGPViewResponse struct {
	Data struct {
		IPv6Prefixes []struct {
			Prefix string `json:"prefix"`
		} `json:"ipv6_prefixes"`
	} `json:"data"`
}

func NewRoutesFetcher(verbose bool) *RoutesFetcher {
	return &RoutesFetcher{
		verbose: verbose,
		routes:  make(map[int][]string),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (rf *RoutesFetcher) log(msg string) {
	if rf.verbose {
		fmt.Fprintf(os.Stderr, "[INFO] %s\n", msg)
	}
}

func (rf *RoutesFetcher) fetchFromRIPE(asn int) ([]string, error) {
	url := fmt.Sprintf("https://stat.ripe.net/data/announced-prefixes/data.json?resource=AS%d", asn)

	resp, err := rf.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ripeResp RIPEResponse
	if err := json.Unmarshal(body, &ripeResp); err != nil {
		return nil, err
	}

	var ipv6Prefixes []string
	for _, p := range ripeResp.Data.Prefixes {
		if strings.Contains(p.Prefix, ":") {
			ipv6Prefixes = append(ipv6Prefixes, p.Prefix)
		}
	}

	rf.log(fmt.Sprintf("AS%d: Found %d prefixes from RIPE", asn, len(ipv6Prefixes)))
	return ipv6Prefixes, nil
}

func (rf *RoutesFetcher) fetchFromBGPView(asn int) ([]string, error) {
	url := fmt.Sprintf("https://api.bgpview.io/asn/%d/prefixes", asn)

	resp, err := rf.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var bgpResp BGPViewResponse
	if err := json.Unmarshal(body, &bgpResp); err != nil {
		return nil, err
	}

	var ipv6Prefixes []string
	for _, p := range bgpResp.Data.IPv6Prefixes {
		ipv6Prefixes = append(ipv6Prefixes, p.Prefix)
	}

	rf.log(fmt.Sprintf("AS%d: Found %d prefixes from BGPView", asn, len(ipv6Prefixes)))
	return ipv6Prefixes, nil
}

func (rf *RoutesFetcher) fetchASN(asn int) (int, error) {
	allRoutes := make(map[string]bool)

	// Try RIPE first (most reliable)
	routes, err := rf.fetchFromRIPE(asn)
	if err != nil {
		rf.log(fmt.Sprintf("AS%d: RIPE fetch failed: %v", asn, err))
	} else {
		for _, route := range routes {
			allRoutes[route] = true
		}
	}

	// Try BGPView as backup
	if len(allRoutes) == 0 {
		routes, err := rf.fetchFromBGPView(asn)
		if err != nil {
			rf.log(fmt.Sprintf("AS%d: BGPView fetch failed: %v", asn, err))
		} else {
			for _, route := range routes {
				allRoutes[route] = true
			}
		}
	}

	// Convert map to sorted slice
	var routeSlice []string
	for route := range allRoutes {
		routeSlice = append(routeSlice, route)
	}
	sort.Strings(routeSlice)

	rf.routes[asn] = routeSlice
	return len(routeSlice), nil
}

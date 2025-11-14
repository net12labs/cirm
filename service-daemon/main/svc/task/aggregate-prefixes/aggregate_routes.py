#!/usr/bin/env python3

import sys
import re
from ipaddress import IPv6Network, collapse_addresses
from collections import defaultdict

def extract_routes_by_asn(filename):
    """Extract routes grouped by ASN from the config file."""
    routes_by_asn = defaultdict(list)
    current_asn = None
    
    with open(filename, 'r') as f:
        for line in f:
            # Check if this is an ASN comment line
            asn_match = re.search(r'# Configuration for ASN AS(\d+)', line)
            if asn_match:
                current_asn = asn_match.group(1)
                continue
            
            # Extract route CIDR blocks
            if 'route' in line and current_asn:
                route_match = re.search(r'route\s+([0-9a-f:]+/\d+)\s+via', line)
                if route_match:
                    cidr = route_match.group(1)
                    routes_by_asn[current_asn].append(cidr)
    
    return routes_by_asn

def aggregate_routes(cidrs, prefix_len=None, aggressive=False):
    """Aggregate CIDR blocks by merging overlapping/contiguous ranges.
    
    Args:
        cidrs: List of CIDR strings to aggregate
        prefix_len: Optional target prefix length for aggressive aggregation
        aggressive: If True, find common supernets even with gaps between routes
    
    Returns:
        List of aggregated CIDR strings
    """
    try:
        # Convert strings to IPv6Network objects with strict=False to normalize addresses
        networks = [IPv6Network(cidr, strict=False) for cidr in cidrs]
        
        # First, remove any overlapping networks using collapse_addresses
        aggregated = list(collapse_addresses(networks))
        
        if aggressive:
            # Aggressive mode: find the smallest supernets that contain all networks
            # This will combine routes even if they have gaps between them
            # Run multiple passes until no more aggregation is possible
            prev_count = len(aggregated)
            max_iterations = 5
            for iteration in range(max_iterations):
                aggregated = find_optimal_supernets(aggregated)
                # After aggressive aggregation, we may have created overlaps
                # Collapse again to remove any redundant routes
                aggregated = list(collapse_addresses(aggregated))
                if len(aggregated) >= prev_count:
                    # No more aggregation possible
                    break
                prev_count = len(aggregated)
        
        # If prefix_len is specified, try to aggregate to that length
        if prefix_len is not None:
            aggregated_by_prefix = defaultdict(list)
            for net in aggregated:
                if net.prefixlen > prefix_len:
                    # Network is more specific than target prefix, need to aggregate
                    # Get the supernet at the target prefix length
                    supernet = net.supernet(new_prefix=prefix_len)
                    aggregated_by_prefix[supernet].append(net)
                else:
                    # Network is less specific or equal, keep as is
                    aggregated_by_prefix[net].append(net)
            
            # Return supernets for overly specific networks
            result = list(aggregated_by_prefix.keys())
            # Final collapse to remove any overlaps from prefix aggregation
            result = list(collapse_addresses(result))
        else:
            result = aggregated
        
        # Sort for consistent output
        result.sort()
        return [str(net) for net in result]
    except Exception as e:
        print(f"Error aggregating routes: {e}", file=sys.stderr)
        return cidrs

def find_optimal_supernets(networks, max_waste_ratio=4.0):
    """Find optimal supernets that cover groups of networks intelligently.
    
    Balance between routing table size and specificity:
    - Accept higher waste for more networks being aggregated
    - But limit minimum prefix length to avoid overly broad routes
    - Prevent aggregating unrelated /32 blocks together
    
    Args:
        networks: List of IPv6Network objects
        max_waste_ratio: Base maximum ratio (supernet_size / actual_coverage)
    """
    if not networks:
        return networks
    
    networks = sorted(networks)
    result = []
    used = set()
    
    for i in range(len(networks)):
        if i in used:
            continue
        
        current = networks[i]
        best_group = [current]
        best_supernet = current
        best_indices = {i}
        best_score = 0
        
        # Determine minimum acceptable prefix based on current network
        # Don't aggregate beyond reasonable boundaries
        min_prefix = max(24, current.prefixlen - 16)  # At most 16 bits broader
        
        # Try different prefix lengths to find the best aggregation
        for target_prefix in range(current.prefixlen - 1, min_prefix - 1, -1):
            try:
                # Create a supernet at this prefix length
                supernet = current.supernet(new_prefix=target_prefix)
                
                # Find all networks that are subnets of this supernet
                group = [current]
                indices = {i}
                
                for j in range(i + 1, len(networks)):
                    if j not in used and networks[j].subnet_of(supernet):
                        group.append(networks[j])
                        indices.add(j)
                
                # Calculate if this aggregation is worthwhile
                if len(group) >= 2:
                    # Calculate total coverage
                    total_coverage = sum(2 ** (128 - net.prefixlen) for net in group)
                    supernet_size = 2 ** (128 - supernet.prefixlen)
                    waste_ratio = supernet_size / total_coverage if total_coverage > 0 else float('inf')
                    
                    # Balanced thresholds:
                    # More networks = more acceptable waste, but keep it reasonable
                    if len(group) == 2:
                        adjusted_max = 100  # 2 routes → 1 route, accept 100x waste
                    elif len(group) == 3:
                        adjusted_max = 500
                    elif len(group) <= 5:
                        adjusted_max = 2000
                    elif len(group) <= 10:
                        adjusted_max = 10000
                    else:
                        # 10+ networks: accept high waste
                        adjusted_max = 100000
                    
                    # Calculate efficiency score
                    # Balance: save routes but avoid excessive waste
                    networks_saved = len(group) - 1
                    # Quadratic bonus for more networks (2^n is too aggressive)
                    bonus = len(group) ** 1.5
                    score = networks_saved * bonus / max(1, waste_ratio ** 0.3)
                    
                    # Accept if waste is within threshold and score improved
                    if waste_ratio <= adjusted_max and score > best_score:
                        best_group = group
                        best_supernet = supernet
                        best_indices = indices
                        best_score = score
                        
            except ValueError:
                # Can't create supernet at this prefix
                break
        
        # Use the best aggregation found
        if len(best_group) > 1:
            result.append(best_supernet)
            used.update(best_indices)
        else:
            result.append(current)
            used.add(i)
    
    return sorted(result)

def print_aggregation_analysis(routes_by_asn, aggressive=False):
    """Print analysis of potential aggregations."""
    print("\n# Route Aggregation Analysis", file=sys.stderr)
    if aggressive:
        print("# Mode: AGGRESSIVE (allowing gaps)", file=sys.stderr)
    else:
        print("# Mode: STANDARD (contiguous only)", file=sys.stderr)
    print("# ASN, Original Routes, Aggregated Routes, Reduction %", file=sys.stderr)
    
    total_before = 0
    total_after = 0
    
    for asn in sorted(routes_by_asn.keys(), key=int):
        cidrs = routes_by_asn[asn]
        before = len(cidrs)
        after = len(aggregate_routes(cidrs, aggressive=aggressive))
        total_before += before
        total_after += after
        reduction = 100 * (before - after) / before if before > 0 else 0
        print(f"# AS{asn}: {before:3d} -> {after:3d} ({reduction:5.1f}%)", file=sys.stderr)
    
    print(f"#", file=sys.stderr)
    overall_reduction = 100 * (total_before - total_after) / total_before if total_before > 0 else 0
    print(f"# Total: {total_before} -> {total_after} ({overall_reduction:.1f}% reduction)", file=sys.stderr)

def main():
    if len(sys.argv) < 2:
        print("Usage: aggregate_routes.py <config_file> [options]")
        print("Options:")
        print("  -a, --asn <asn>        Process only specific ASN")
        print("  -p, --prefix <len>     Aggregate to specific prefix length")
        print("  -g, --aggressive       Aggressive aggregation (with gaps)")
        print("  -s, --stats            Show aggregation statistics for all ASNs")
        print("  -c, --compact          Output as compact BIRD config block")
        sys.exit(1)
    
    config_file = sys.argv[1]
    filter_asn = None
    prefix_len = None
    show_stats = False
    compact = False
    aggressive = False
    
    i = 2
    while i < len(sys.argv):
        arg = sys.argv[i]
        if arg in ['-a', '--asn']:
            filter_asn = sys.argv[i + 1]
            i += 2
        elif arg in ['-p', '--prefix']:
            prefix_len = int(sys.argv[i + 1])
            i += 2
        elif arg in ['-g', '--aggressive']:
            aggressive = True
            i += 1
        elif arg in ['-s', '--stats']:
            show_stats = True
            i += 1
        elif arg in ['-c', '--compact']:
            compact = True
            i += 1
        else:
            print(f"Unknown option: {arg}", file=sys.stderr)
            sys.exit(1)
    
    routes_by_asn = extract_routes_by_asn(config_file)
    
    if show_stats:
        print_aggregation_analysis(routes_by_asn, aggressive)
        return
    
    if filter_asn:
        if filter_asn not in routes_by_asn:
            print(f"ASN {filter_asn} not found in config file", file=sys.stderr)
            sys.exit(1)
        asns_to_process = [filter_asn]
    else:
        # Process all ASNs sorted numerically
        asns_to_process = sorted(routes_by_asn.keys(), key=int)
    
    total_before = 0
    total_after = 0
    
    if compact:
        print("protocol static {")
        print("    ipv6;")
    
    for asn in asns_to_process:
        cidrs = routes_by_asn[asn]
        total_before += len(cidrs)
        
        aggregated = aggregate_routes(cidrs, prefix_len, aggressive)
        total_after += len(aggregated)
        
        if not compact:
            print(f"# Configuration for ASN AS{asn}")
        else:
            print(f"    # Configuration for ASN AS{asn}")
        for cidr in aggregated:
            print(f"    route {cidr} via 2::1:9;")
        if not compact:
            print()
    
    if compact:
        print("}")
    
    # Print statistics
    print(f"# Aggregation Summary:", file=sys.stderr)
    print(f"# Total routes before: {total_before}", file=sys.stderr)
    print(f"# Total routes after:  {total_after}", file=sys.stderr)
    if total_before > 0:
        print(f"# Reduction: {total_before - total_after} routes ({100 * (total_before - total_after) / total_before:.1f}%)", file=sys.stderr)

if __name__ == '__main__':
    main()

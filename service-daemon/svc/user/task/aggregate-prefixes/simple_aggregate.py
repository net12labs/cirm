#!/usr/bin/env python3
"""
Simple Multi-Pass Route Aggregator

Strategy:
1. Start with all routes
2. For each pass, try to find pairs or groups that can be aggregated
3. Keep the best aggregations (most routes saved with reasonable coverage)
4. Repeat until no more aggregations are beneficial
"""

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
            asn_match = re.search(r'# Configuration for ASN AS(\d+)', line)
            if asn_match:
                current_asn = asn_match.group(1)
                continue
            
            if 'route' in line and current_asn:
                route_match = re.search(r'route\s+([0-9a-f:]+/\d+)\s+via', line)
                if route_match:
                    cidr = route_match.group(1)
                    routes_by_asn[current_asn].append(cidr)
    
    return routes_by_asn

def find_aggregation_candidates(networks):
    """
    Find all possible aggregations and score them.
    Returns list of (score, supernet, covered_networks) sorted by score.
    """
    candidates = []
    import math
    
    for i, net in enumerate(networks):
        # Try different supernet sizes - go deeper to find more aggregations
        for prefix_len in range(net.prefixlen - 1, max(20, net.prefixlen - 20), -1):
            try:
                supernet = net.supernet(new_prefix=prefix_len)
                
                # Find all networks covered by this supernet
                covered = []
                for other_net in networks:
                    if other_net.subnet_of(supernet):
                        covered.append(other_net)
                
                if len(covered) >= 2:
                    # Calculate score: routes_saved / log(waste_ratio)
                    # But boost score for large aggregations
                    routes_saved = len(covered) - 1
                    total_coverage = sum(2 ** (128 - n.prefixlen) for n in covered)
                    supernet_size = 2 ** (128 - supernet.prefixlen)
                    waste_ratio = supernet_size / total_coverage if total_coverage > 0 else float('inf')
                    
                    # Enhanced scoring:
                    # - Base: routes_saved / log(waste)
                    # - Bonus for many routes: multiply by sqrt(routes_saved)
                    # - This makes aggregating 10+ routes more attractive
                    waste_penalty = math.log10(max(1, waste_ratio))
                    base_score = routes_saved / (1 + waste_penalty)
                    
                    # Exponential bonus for more routes
                    if len(covered) >= 10:
                        bonus = 3.0  # Strong bonus for 10+ routes
                    elif len(covered) >= 5:
                        bonus = 2.0  # Good bonus for 5+ routes
                    elif len(covered) >= 3:
                        bonus = 1.5  # Small bonus for 3+ routes
                    else:
                        bonus = 1.0  # No bonus for 2 routes
                    
                    score = base_score * bonus
                    
                    candidates.append((score, supernet, frozenset(covered)))
                    
            except ValueError:
                continue
    
    # Remove duplicates (same supernet covering same networks)
    seen = set()
    unique_candidates = []
    for score, supernet, covered in candidates:
        key = (supernet, covered)
        if key not in seen:
            seen.add(key)
            unique_candidates.append((score, supernet, covered))
    
    # Sort by score (higher is better)
    unique_candidates.sort(reverse=True)
    return unique_candidates

def aggregate_one_pass(networks):
    """
    Single pass: find the best non-overlapping aggregations.
    Returns new list of networks after aggregation.
    """
    if len(networks) <= 1:
        return networks
    
    # Find all possible aggregations
    candidates = find_aggregation_candidates(networks)
    
    if not candidates:
        return networks
    
    # Greedily select best non-overlapping aggregations
    used_networks = set()
    selected_supernets = []
    
    for score, supernet, covered in candidates:
        # Check if any of the covered networks are already used
        if not any(net in used_networks for net in covered):
            # This aggregation is valid
            selected_supernets.append(supernet)
            used_networks.update(covered)
    
    # Add networks that weren't aggregated
    result = list(selected_supernets)
    for net in networks:
        if net not in used_networks:
            result.append(net)
    
    return sorted(result)

def aggregate_routes_simple(cidrs, max_passes=10):
    """
    Simple multi-pass aggregation.
    Keep aggregating until no more improvement or max passes reached.
    """
    # Convert to networks and do initial collapse
    networks = [IPv6Network(cidr, strict=False) for cidr in cidrs]
    networks = list(collapse_addresses(networks))
    
    for pass_num in range(max_passes):
        prev_count = len(networks)
        
        # Try one pass of aggregation
        networks = aggregate_one_pass(networks)
        
        # Collapse to merge any adjacent networks created
        networks = list(collapse_addresses(networks))
        
        # Check if we made progress
        if len(networks) >= prev_count:
            break
    
    networks.sort()
    return [str(net) for net in networks]

def print_stats(routes_by_asn):
    """Print aggregation statistics."""
    print("\n# Route Aggregation Analysis (Simple Multi-Pass)", file=sys.stderr)
    print("# ASN, Original Routes, Aggregated Routes, Reduction %", file=sys.stderr)
    
    total_before = 0
    total_after = 0
    
    for asn in sorted(routes_by_asn.keys(), key=int):
        cidrs = routes_by_asn[asn]
        before = len(cidrs)
        after = len(aggregate_routes_simple(cidrs))
        total_before += before
        total_after += after
        reduction = 100 * (before - after) / before if before > 0 else 0
        print(f"# AS{asn}: {before:3d} -> {after:3d} ({reduction:5.1f}%)", file=sys.stderr)
    
    print(f"#", file=sys.stderr)
    overall_reduction = 100 * (total_before - total_after) / total_before if total_before > 0 else 0
    print(f"# Total: {total_before} -> {total_after} ({overall_reduction:.1f}% reduction)", file=sys.stderr)

def main():
    if len(sys.argv) < 2:
        print("Usage: simple_aggregate.py <config_file> [--stats|--compact]")
        sys.exit(1)
    
    config_file = sys.argv[1]
    show_stats = '--stats' in sys.argv
    compact = '--compact' in sys.argv
    
    routes_by_asn = extract_routes_by_asn(config_file)
    
    if show_stats:
        print_stats(routes_by_asn)
        return
    
    asns_to_process = sorted(routes_by_asn.keys(), key=int)
    total_before = 0
    total_after = 0
    
    if compact:
        print("protocol static {")
        print("    ipv6;")
    
    for asn in asns_to_process:
        cidrs = routes_by_asn[asn]
        total_before += len(cidrs)
        
        aggregated = aggregate_routes_simple(cidrs)
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
    
    # Print statistics to stderr
    print(f"# Aggregation Summary:", file=sys.stderr)
    print(f"# Total routes before: {total_before}", file=sys.stderr)
    print(f"# Total routes after:  {total_after}", file=sys.stderr)
    if total_before > 0:
        print(f"# Reduction: {total_before - total_after} routes ({100 * (total_before - total_after) / total_before:.1f}%)", file=sys.stderr)

if __name__ == '__main__':
    main()

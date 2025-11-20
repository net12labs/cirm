// filepath: /home/lxk/Desktop/cirm/bins/vfsql/search.go
package vfsql

import (
	"fmt"
	"strings"
)

// Search performs advanced search with multiple criteria
func (vfs *VFS) Search(query *SearchQuery) (*SearchResults, error) {
	if query == nil {
		query = &SearchQuery{}
	}

	// Build SQL query
	sqlQuery, args := vfs.buildSearchQuery(query, false)

	// Get total count if limit is set
	var totalCount int
	if query.Limit > 0 {
		countQuery, countArgs := vfs.buildSearchQuery(query, true)
		err := vfs.db.db.QueryRow(countQuery, countArgs...).Scan(&totalCount)
		if err != nil {
			return nil, err
		}
	}

	// Execute search query
	rows, err := vfs.db.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// If no limit was set, total count is the result count
	if query.Limit == 0 {
		totalCount = len(paths)
	}

	hasMore := query.Limit > 0 && totalCount > query.Offset+len(paths)

	return &SearchResults{
		Paths:      paths,
		TotalCount: totalCount,
		HasMore:    hasMore,
	}, nil
}

// buildSearchQuery constructs the SQL query for search
func (vfs *VFS) buildSearchQuery(query *SearchQuery, countOnly bool) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	// Always filter by filesystem
	conditions = append(conditions, "fs_id = ?")
	args = append(args, vfs.id)

	// Type filter
	if query.Type == FileTypeFile {
		conditions = append(conditions, "type = 'file'")
	} else if query.Type == FileTypeDir {
		conditions = append(conditions, "type = 'dir'")
	}

	// Name pattern
	if query.NamePattern != "" {
		likePattern := globToLike(query.NamePattern)
		conditions = append(conditions, "name LIKE ?")
		args = append(args, likePattern)
	}

	// Description filter
	if query.Description != "" {
		conditions = append(conditions, "description LIKE ?")
		args = append(args, "%"+query.Description+"%")
	}

	// Tags filter
	if len(query.Tags) > 0 {
		if query.TagMatchAll {
			// All tags must be present
			for _, tag := range query.Tags {
				conditions = append(conditions, "tags LIKE ?")
				args = append(args, "%"+tag+"%")
			}
		} else {
			// At least one tag must be present
			tagConditions := make([]string, len(query.Tags))
			for i, tag := range query.Tags {
				tagConditions[i] = "tags LIKE ?"
				args = append(args, "%"+tag+"%")
			}
			conditions = append(conditions, "("+strings.Join(tagConditions, " OR ")+")")
		}
	}

	// Size filters
	if query.MinSize > 0 {
		conditions = append(conditions, "size >= ?")
		args = append(args, query.MinSize)
	}
	if query.MaxSize > 0 {
		conditions = append(conditions, "size <= ?")
		args = append(args, query.MaxSize)
	}

	// Time filters
	if query.CreatedAfter > 0 {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, query.CreatedAfter)
	}
	if query.CreatedBefore > 0 {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, query.CreatedBefore)
	}
	if query.ModifiedAfter > 0 {
		conditions = append(conditions, "modified_at >= ?")
		args = append(args, query.ModifiedAfter)
	}
	if query.ModifiedBefore > 0 {
		conditions = append(conditions, "modified_at <= ?")
		args = append(args, query.ModifiedBefore)
	}
	if query.AccessedAfter > 0 {
		conditions = append(conditions, "accessed_at >= ?")
		args = append(args, query.AccessedAfter)
	}
	if query.AccessedBefore > 0 {
		conditions = append(conditions, "accessed_at <= ?")
		args = append(args, query.AccessedBefore)
	}

	// Permission filters
	if query.Mode != 0 {
		conditions = append(conditions, "mode = ?")
		args = append(args, query.Mode)
	}
	if query.UID >= 0 {
		conditions = append(conditions, "uid = ?")
		args = append(args, query.UID)
	}
	if query.GID >= 0 {
		conditions = append(conditions, "gid = ?")
		args = append(args, query.GID)
	}

	// Build WHERE clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	if countOnly {
		return "SELECT COUNT(*) FROM inodes" + whereClause, args
	}

	// Build CTE for path reconstruction
	cte := `
	WITH RECURSIVE path_builder(id, path, level) AS (
		SELECT id, name, 0 FROM inodes WHERE fs_id = ? AND parent_id IS NULL
		UNION ALL
		SELECT i.id, pb.path || '/' || i.name, pb.level + 1
		FROM inodes i
		JOIN path_builder pb ON i.parent_id = pb.id
		WHERE i.fs_id = ?
	)
	`

	// Add fs_id to beginning of args for CTE
	cteArgs := []interface{}{vfs.id, vfs.id}
	cteArgs = append(cteArgs, args...)

	// Base path filter (needs to be done after path reconstruction)
	basePathFilter := ""
	if query.BasePath != "" {
		basePath := normalizePath(query.BasePath)
		if query.Recursive {
			basePathFilter = " AND (pb.path = ? OR pb.path LIKE ?)"
			cteArgs = append(cteArgs, basePath, basePath+"/%")
		} else {
			// Non-recursive: direct children only
			basePathFilter = " AND pb.path = ?"
			cteArgs = append(cteArgs, basePath)
		}
	}

	// Max depth filter
	depthFilter := ""
	if query.MaxDepth > 0 && query.BasePath != "" {
		// Calculate base depth
		baseDepth := len(strings.Split(strings.Trim(query.BasePath, "/"), "/"))
		if query.BasePath == "/" {
			baseDepth = 0
		}
		maxLevel := baseDepth + query.MaxDepth
		depthFilter = fmt.Sprintf(" AND pb.level <= %d", maxLevel)
	}

	// Sort clause
	sortClause := ""
	if !countOnly {
		sortField := "pb.path"
		switch query.SortBy {
		case SortByName:
			sortField = "i.name"
		case SortBySize:
			sortField = "i.size"
		case SortByCreated:
			sortField = "i.created_at"
		case SortByModified:
			sortField = "i.modified_at"
		case SortByAccessed:
			sortField = "i.accessed_at"
		}

		sortDir := "ASC"
		if query.SortOrder == Descending {
			sortDir = "DESC"
		}

		sortClause = fmt.Sprintf(" ORDER BY %s %s", sortField, sortDir)
	}

	// Limit and offset
	limitClause := ""
	if query.Limit > 0 {
		limitClause = fmt.Sprintf(" LIMIT %d OFFSET %d", query.Limit, query.Offset)
	}

	// Full query
	fullQuery := cte + `
		SELECT pb.path FROM path_builder pb
		JOIN inodes i ON pb.id = i.id
	` + whereClause + basePathFilter + depthFilter + sortClause + limitClause

	return fullQuery, cteArgs
}

// FindByName performs a simple name-based search
func (vfs *VFS) FindByName(pattern string, opts *FindOptions) ([]string, error) {
	if opts == nil {
		opts = &FindOptions{}
	}

	query := &SearchQuery{
		NamePattern: pattern,
		BasePath:    opts.BasePath,
		Recursive:   opts.Recursive,
		Type:        opts.Type,
		Limit:       opts.Limit,
		SortBy:      SortByPath,
		SortOrder:   Ascending,
	}

	results, err := vfs.Search(query)
	if err != nil {
		return nil, err
	}

	return results.Paths, nil
}

// FindByTag performs a tag-based search
func (vfs *VFS) FindByTag(tags []string, matchAll bool) ([]string, error) {
	query := &SearchQuery{
		Tags:        tags,
		TagMatchAll: matchAll,
		BasePath:    "/",
		Recursive:   true,
		SortBy:      SortByPath,
		SortOrder:   Ascending,
	}

	results, err := vfs.Search(query)
	if err != nil {
		return nil, err
	}

	return results.Paths, nil
}

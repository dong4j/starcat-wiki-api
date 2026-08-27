package store

import (
	"database/sql"
	"time"
)

// WikiOperationalStats describes probe coverage, status distribution, retries, and freshness.
type WikiOperationalStats struct {
	Repositories  int            `json:"repositories"`
	TotalProbes   int            `json:"total_probes"`
	ExpiredProbes int            `json:"expired_probes"`
	Retryable     int            `json:"retryable"`
	BySource      map[string]int `json:"by_source"`
	ByStatus      map[string]int `json:"by_status"`
	HTTPStatuses  map[string]int `json:"http_statuses"`
	LastCheckedAt string         `json:"last_checked_at,omitempty"`
}

// ProbeError is a bounded diagnostic row. It excludes matched signals and response content.
type ProbeError struct {
	Owner         string `json:"owner"`
	Repo          string `json:"repo"`
	Source        string `json:"source"`
	HTTPStatus    *int   `json:"http_status,omitempty"`
	Attempt       int    `json:"attempt"`
	LastError     string `json:"last_error"`
	NextRetryAt   string `json:"next_retry_at,omitempty"`
	LastCheckedAt string `json:"last_checked_at"`
}

// GetWikiOperationalStats aggregates operational state without triggering network probes.
func (s *SQLiteStore) GetWikiOperationalStats(now time.Time, maxAttempts int) (WikiOperationalStats, error) {
	var result WikiOperationalStats
	var lastChecked sql.NullString
	err := s.db.QueryRow(`
SELECT COUNT(DISTINCT owner || '/' || repo), COUNT(*),
       COALESCE(SUM(CASE WHEN expires_at < ? THEN 1 ELSE 0 END), 0),
       COALESCE(SUM(CASE WHEN status='error' AND attempt < ? AND
         (next_retry_at IS NULL OR next_retry_at <= ?) THEN 1 ELSE 0 END), 0),
       MAX(checked_at)
FROM doc_probes`, now.UTC().Format(time.RFC3339), maxAttempts, now.UTC().Format(time.RFC3339)).
		Scan(&result.Repositories, &result.TotalProbes, &result.ExpiredProbes, &result.Retryable, &lastChecked)
	if err != nil {
		return WikiOperationalStats{}, err
	}
	result.LastCheckedAt = lastChecked.String
	result.BySource, err = s.groupCounts("source")
	if err != nil {
		return WikiOperationalStats{}, err
	}
	result.ByStatus, err = s.groupCounts("status")
	if err != nil {
		return WikiOperationalStats{}, err
	}
	result.HTTPStatuses, err = s.groupCounts("CAST(http_status AS TEXT)")
	if err != nil {
		return WikiOperationalStats{}, err
	}
	return result, nil
}

func (s *SQLiteStore) groupCounts(column string) (map[string]int, error) {
	// column is selected from code, not request input.
	rows, err := s.db.Query("SELECT " + column + ", COUNT(*) FROM doc_probes WHERE " + column + " IS NOT NULL GROUP BY " + column)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]int)
	for rows.Next() {
		var key string
		var count int
		if err := rows.Scan(&key, &count); err != nil {
			return nil, err
		}
		result[key] = count
	}
	return result, rows.Err()
}

// ListRecentProbeErrors returns the most recently checked errors with a hard caller-provided limit.
func (s *SQLiteStore) ListRecentProbeErrors(limit int) ([]ProbeError, error) {
	rows, err := s.db.Query(`
SELECT owner, repo, source, http_status, attempt, COALESCE(last_error, ''),
       COALESCE(next_retry_at, ''), checked_at
FROM doc_probes WHERE status='error'
ORDER BY checked_at DESC, owner, repo, source LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]ProbeError, 0, limit)
	for rows.Next() {
		var item ProbeError
		var httpStatus sql.NullInt64
		if err := rows.Scan(&item.Owner, &item.Repo, &item.Source, &httpStatus, &item.Attempt,
			&item.LastError, &item.NextRetryAt, &item.LastCheckedAt); err != nil {
			return nil, err
		}
		if httpStatus.Valid {
			value := int(httpStatus.Int64)
			item.HTTPStatus = &value
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

// Package handler exposes read-only Wiki probe operations data.
package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/starcat-app/starcat-wiki-api/internal/store"
)

type wikiOperationsStore interface {
	GetWikiOperationalStats(now time.Time, maxAttempts int) (store.WikiOperationalStats, error)
	ListRecentProbeErrors(limit int) ([]store.ProbeError, error)
}

// HandleOperationalStats returns probe coverage, statuses, retries, and freshness.
func HandleOperationalStats(s wikiOperationsStore) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		stats, err := s.GetWikiOperationalStats(time.Now().UTC(), 5)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to aggregate wiki stats", nil)
			return
		}
		writeJSON(w, stats)
	}
}

// HandleRecentProbeErrors returns a bounded error table for diagnostics.
func HandleRecentProbeErrors(s wikiOperationsStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 25
		if raw := r.URL.Query().Get("limit"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil || parsed < 1 || parsed > 100 {
				writeError(w, http.StatusBadRequest, "INVALID_LIMIT", "limit must be between 1 and 100", nil)
				return
			}
			limit = parsed
		}
		items, err := s.ListRecentProbeErrors(limit)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list probe errors", nil)
			return
		}
		writeJSON(w, items)
	}
}

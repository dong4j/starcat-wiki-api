package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/starcat-app/starcat-wiki-api/internal/probe"
)

func TestWikiOperationalStatsAndRecentErrors(t *testing.T) {
	s, err := NewSQLiteStore(filepath.Join(t.TempDir(), "wiki.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	ctx := context.Background()
	if err := s.InsertPendingProbes(ctx, "starcat", "app"); err != nil {
		t.Fatal(err)
	}
	status := 500
	now := time.Now().UTC()
	if err := s.UpdateProbeResult(ctx, "starcat", "app", probe.SourceDeepWiki,
		probe.ProbeResult{Source: probe.SourceDeepWiki, Status: probe.StatusError,
			URL: "https://deepwiki.com/starcat/app", HTTPStatus: &status},
		2, "upstream timeout", now.Add(-time.Minute).Format(time.RFC3339)); err != nil {
		t.Fatal(err)
	}
	stats, err := s.GetWikiOperationalStats(now, 5)
	if err != nil {
		t.Fatal(err)
	}
	if stats.Repositories != 1 || stats.TotalProbes != 3 || stats.ByStatus[string(probe.StatusError)] != 1 ||
		stats.Retryable != 1 || stats.HTTPStatuses["500"] != 1 {
		t.Fatalf("unexpected operational stats: %#v", stats)
	}
	errors, err := s.ListRecentProbeErrors(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(errors) != 1 || errors[0].Owner != "starcat" || errors[0].Attempt != 2 || errors[0].LastError != "upstream timeout" {
		t.Fatalf("unexpected errors: %#v", errors)
	}
}

package store

import (
	"database/sql"
	"testing"
)

func newTestAIModelStore(t *testing.T) (*AIModelStore, func()) {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	// Ensure we keep a single connection for :memory: databases.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	store := &AIModelStore{db: db}
	if err := store.initTables(); err != nil {
		_ = db.Close()
		t.Fatalf("initTables: %v", err)
	}

	return store, func() {
		_ = db.Close()
	}
}

func TestAIModelStoreUpdate_NewSemantics_AllowsMultiplePerProvider(t *testing.T) {
	s, cleanup := newTestAIModelStore(t)
	defer cleanup()

	userID := "user1"
	provider := "deepseek"

	if err := s.Update(userID, "modelA", provider, true, "sk-aaa", "", "deepseek-chat"); err != nil {
		t.Fatalf("Update new semantics modelA: %v", err)
	}
	if err := s.Update(userID, "modelB", provider, true, "sk-bbb", "", "deepseek-reasoner"); err != nil {
		t.Fatalf("Update new semantics modelB: %v", err)
	}

	var count int
	if err := s.db.QueryRow(
		`SELECT COUNT(*) FROM ai_models WHERE user_id = ? AND provider = ?`,
		userID,
		provider,
	).Scan(&count); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 rows for provider %q, got %d", provider, count)
	}

	found := map[string]bool{}
	rows, err := s.db.Query(
		`SELECT id, custom_model_name FROM ai_models WHERE user_id = ? AND provider = ?`,
		userID,
		provider,
	)
	if err != nil {
		t.Fatalf("query rows: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, customModelName string
		if err := rows.Scan(&id, &customModelName); err != nil {
			t.Fatalf("scan: %v", err)
		}
		found[customModelName] = true
		if id == "" {
			t.Fatalf("expected non-empty id")
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows err: %v", err)
	}

	if !found["deepseek-chat"] || !found["deepseek-reasoner"] {
		t.Fatalf("expected both custom_model_name values to exist, got: %+v", found)
	}
}

func TestAIModelStoreUpdate_NewSemantics_GeneratesIDWhenEmpty(t *testing.T) {
	s, cleanup := newTestAIModelStore(t)
	defer cleanup()

	userID := "user1"
	provider := "deepseek"

	if err := s.Update(userID, "", provider, true, "sk-aaa", "", "deepseek-chat"); err != nil {
		t.Fatalf("Update new semantics (empty id): %v", err)
	}

	var id string
	if err := s.db.QueryRow(
		`SELECT id FROM ai_models WHERE user_id = ? AND provider = ? LIMIT 1`,
		userID,
		provider,
	).Scan(&id); err != nil {
		t.Fatalf("select id: %v", err)
	}
	if id == "" {
		t.Fatalf("expected generated id to be non-empty")
	}
}

func TestAIModelStoreUpdate_Legacy_SingleRecordPerProvider(t *testing.T) {
	s, cleanup := newTestAIModelStore(t)
	defer cleanup()

	userID := "user1"
	provider := "deepseek"

	if err := s.Update(userID, provider, "", true, "sk-aaa", "", "deepseek-chat"); err != nil {
		t.Fatalf("Update legacy #1: %v", err)
	}
	if err := s.Update(userID, provider, "", true, "sk-bbb", "", "deepseek-reasoner"); err != nil {
		t.Fatalf("Update legacy #2: %v", err)
	}

	var count int
	if err := s.db.QueryRow(
		`SELECT COUNT(*) FROM ai_models WHERE user_id = ? AND provider = ?`,
		userID,
		provider,
	).Scan(&count); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row for provider %q in legacy mode, got %d", provider, count)
	}

	var id, customModelName string
	if err := s.db.QueryRow(
		`SELECT id, custom_model_name FROM ai_models WHERE user_id = ? AND provider = ? LIMIT 1`,
		userID,
		provider,
	).Scan(&id, &customModelName); err != nil {
		t.Fatalf("select id/custom_model_name: %v", err)
	}

	if id != "user1_deepseek" {
		t.Fatalf("expected legacy id to be %q, got %q", "user1_deepseek", id)
	}
	if customModelName != "deepseek-reasoner" {
		t.Fatalf("expected custom_model_name to be updated, got %q", customModelName)
	}
}

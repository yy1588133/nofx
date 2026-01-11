package store

import (
	"path/filepath"
	"testing"
)

func newTestAIModelStore(t *testing.T) (*AIModelStore, func()) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "ai_model_test.db")
	gdb, err := InitGorm(dbPath)
	if err != nil {
		t.Fatalf("InitGorm: %v", err)
	}

	s := NewAIModelStore(gdb)
	if err := s.initTables(); err != nil {
		sqlDB, _ := gdb.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		t.Fatalf("initTables: %v", err)
	}

	return s, func() {
		sqlDB, err := gdb.DB()
		if err == nil && sqlDB != nil {
			_ = sqlDB.Close()
		}
	}
}

func TestAIModelStoreUpdate_AllowsMultiplePerProvider_ByUniqueID(t *testing.T) {
	s, cleanup := newTestAIModelStore(t)
	defer cleanup()

	userID := "user1"
	provider := "deepseek"

	id1 := userID + "_1_" + provider
	id2 := userID + "_2_" + provider

	if err := s.Update(userID, id1, true, "sk-aaa", "", "deepseek-chat"); err != nil {
		t.Fatalf("Update model #1: %v", err)
	}
	if err := s.Update(userID, id2, true, "sk-bbb", "", "deepseek-reasoner"); err != nil {
		t.Fatalf("Update model #2: %v", err)
	}

	var count int64
	if err := s.db.Model(&AIModel{}).
		Where("user_id = ? AND provider = ?", userID, provider).
		Count(&count).Error; err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 rows for provider %q, got %d", provider, count)
	}

	var models []AIModel
	if err := s.db.Where("user_id = ? AND provider = ?", userID, provider).Find(&models).Error; err != nil {
		t.Fatalf("query models: %v", err)
	}

	found := map[string]bool{}
	for _, m := range models {
		found[m.CustomModelName] = true
		if m.ID == "" {
			t.Fatalf("expected non-empty model id")
		}
	}

	if !found["deepseek-chat"] || !found["deepseek-reasoner"] {
		t.Fatalf("expected both custom_model_name values to exist, got: %+v", found)
	}
}

func TestAIModelStoreUpdate_LegacyProviderKey_CreatesUserScopedID_AndPreservesAPIKey(t *testing.T) {
	s, cleanup := newTestAIModelStore(t)
	defer cleanup()

	userID := "user1"
	provider := "deepseek"

	if err := s.Update(userID, provider, true, "sk-aaa", "", "deepseek-chat"); err != nil {
		t.Fatalf("Update legacy #1: %v", err)
	}
	// apiKey 为空时，应该保持原有密钥不被覆盖
	if err := s.Update(userID, provider, true, "", "", "deepseek-reasoner"); err != nil {
		t.Fatalf("Update legacy #2: %v", err)
	}

	var model AIModel
	if err := s.db.Where("user_id = ? AND provider = ?", userID, provider).First(&model).Error; err != nil {
		t.Fatalf("select model: %v", err)
	}

	if model.ID != "user1_deepseek" {
		t.Fatalf("expected legacy id to be %q, got %q", "user1_deepseek", model.ID)
	}
	if model.CustomModelName != "deepseek-reasoner" {
		t.Fatalf("expected custom_model_name to be updated, got %q", model.CustomModelName)
	}
	if model.APIKey.String() != "sk-aaa" {
		t.Fatalf("expected API key to be preserved, got %q", model.APIKey.String())
	}
}


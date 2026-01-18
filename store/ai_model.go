package store

import (
	"errors"
	"fmt"
	"nofx/crypto"
	"nofx/logger"
	"strings"
	"time"

	"gorm.io/gorm"
)

// AIModelStore AI model storage
type AIModelStore struct {
	db *gorm.DB
}

// AIModel AI model configuration
type AIModel struct {
	ID              string                 `gorm:"primaryKey" json:"id"`
	UserID          string                 `gorm:"column:user_id;not null;default:default;index" json:"user_id"`
	Name            string                 `gorm:"not null" json:"name"`
	Provider        string                 `gorm:"not null" json:"provider"`
	Enabled         bool                   `gorm:"default:false" json:"enabled"`
	APIKey          crypto.EncryptedString `gorm:"column:api_key;default:''" json:"apiKey"`
	CustomAPIURL    string                 `gorm:"column:custom_api_url;default:''" json:"customApiUrl"`
	CustomModelName string                 `gorm:"column:custom_model_name;default:''" json:"customModelName"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

func (AIModel) TableName() string { return "ai_models" }

// NewAIModelStore creates a new AIModelStore
func NewAIModelStore(db *gorm.DB) *AIModelStore {
	return &AIModelStore{db: db}
}

func (s *AIModelStore) initTables() error {
	// For PostgreSQL with existing table, skip AutoMigrate
	if s.db.Dialector.Name() == "postgres" {
		var tableExists int64
		s.db.Raw(`SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'ai_models'`).Scan(&tableExists)
		if tableExists > 0 {
			return nil
		}
	}

	// SQLite: For legacy DBs, GORM AutoMigrate may rebuild table via ai_models__temp
	// and can fail during copy. For existing table, do safe incremental migration.
	if s.db.Dialector.Name() == "sqlite" && s.db.Migrator().HasTable(&AIModel{}) {
		return s.ensureSQLiteAIModelsCompatibility()
	}

	return s.db.AutoMigrate(&AIModel{})
}

type sqliteAIModelTableInfo struct {
	Name string `gorm:"column:name"`
}

func (s *AIModelStore) ensureSQLiteAIModelsCompatibility() error {
	var cols []sqliteAIModelTableInfo
	if err := s.db.Raw("PRAGMA table_info(ai_models)").Scan(&cols).Error; err != nil {
		return fmt.Errorf("failed to inspect ai_models table: %w", err)
	}

	colExists := make(map[string]bool, len(cols))
	for _, c := range cols {
		colExists[strings.ToLower(strings.TrimSpace(c.Name))] = true
	}

	changed := false

	// Add missing columns (keep them nullable where needed; fill data below).
	if !colExists["user_id"] {
		if err := s.db.Exec("ALTER TABLE ai_models ADD COLUMN user_id TEXT NOT NULL DEFAULT 'default'").Error; err != nil {
			return fmt.Errorf("failed to add ai_models.user_id column: %w", err)
		}
		colExists["user_id"] = true
		changed = true
	}
	if !colExists["name"] {
		if err := s.db.Exec("ALTER TABLE ai_models ADD COLUMN name TEXT NOT NULL DEFAULT 'AI Model'").Error; err != nil {
			return fmt.Errorf("failed to add ai_models.name column: %w", err)
		}
		colExists["name"] = true
		changed = true
	}
	if !colExists["provider"] {
		if err := s.db.Exec("ALTER TABLE ai_models ADD COLUMN provider TEXT NOT NULL DEFAULT 'unknown'").Error; err != nil {
			return fmt.Errorf("failed to add ai_models.provider column: %w", err)
		}
		colExists["provider"] = true
		changed = true
	}
	if !colExists["enabled"] {
		if err := s.db.Exec("ALTER TABLE ai_models ADD COLUMN enabled BOOLEAN DEFAULT 0").Error; err != nil {
			return fmt.Errorf("failed to add ai_models.enabled column: %w", err)
		}
		colExists["enabled"] = true
		changed = true
	}
	if !colExists["api_key"] {
		if err := s.db.Exec("ALTER TABLE ai_models ADD COLUMN api_key TEXT DEFAULT ''").Error; err != nil {
			return fmt.Errorf("failed to add ai_models.api_key column: %w", err)
		}
		colExists["api_key"] = true
		changed = true
	}
	if !colExists["custom_api_url"] {
		if err := s.db.Exec("ALTER TABLE ai_models ADD COLUMN custom_api_url TEXT DEFAULT ''").Error; err != nil {
			return fmt.Errorf("failed to add ai_models.custom_api_url column: %w", err)
		}
		colExists["custom_api_url"] = true
		changed = true
	}
	if !colExists["custom_model_name"] {
		if err := s.db.Exec("ALTER TABLE ai_models ADD COLUMN custom_model_name TEXT DEFAULT ''").Error; err != nil {
			return fmt.Errorf("failed to add ai_models.custom_model_name column: %w", err)
		}
		colExists["custom_model_name"] = true
		changed = true
	}
	if !colExists["created_at"] {
		if err := s.db.Exec("ALTER TABLE ai_models ADD COLUMN created_at DATETIME DEFAULT CURRENT_TIMESTAMP").Error; err != nil {
			return fmt.Errorf("failed to add ai_models.created_at column: %w", err)
		}
		colExists["created_at"] = true
		changed = true
	}
	if !colExists["updated_at"] {
		if err := s.db.Exec("ALTER TABLE ai_models ADD COLUMN updated_at DATETIME DEFAULT CURRENT_TIMESTAMP").Error; err != nil {
			return fmt.Errorf("failed to add ai_models.updated_at column: %w", err)
		}
		colExists["updated_at"] = true
		changed = true
	}

	// Data fix for legacy rows.
	if colExists["user_id"] {
		if err := s.db.Exec(`UPDATE ai_models SET user_id = 'default' WHERE user_id IS NULL OR trim(user_id) = ''`).Error; err != nil {
			return fmt.Errorf("failed to patch ai_models.user_id: %w", err)
		}
	}
	if colExists["name"] {
		if err := s.db.Exec(`UPDATE ai_models SET name = 'AI Model' WHERE name IS NULL OR trim(name) = ''`).Error; err != nil {
			return fmt.Errorf("failed to patch ai_models.name: %w", err)
		}
	}
	if colExists["provider"] {
		if err := s.db.Exec(`UPDATE ai_models SET provider = 'unknown' WHERE provider IS NULL OR trim(provider) = ''`).Error; err != nil {
			return fmt.Errorf("failed to patch ai_models.provider: %w", err)
		}
	}
	if colExists["api_key"] {
		if err := s.db.Exec(`UPDATE ai_models SET api_key = '' WHERE api_key IS NULL`).Error; err != nil {
			return fmt.Errorf("failed to patch ai_models.api_key: %w", err)
		}
	}
	if colExists["custom_api_url"] {
		if err := s.db.Exec(`UPDATE ai_models SET custom_api_url = '' WHERE custom_api_url IS NULL`).Error; err != nil {
			return fmt.Errorf("failed to patch ai_models.custom_api_url: %w", err)
		}
	}
	if colExists["custom_model_name"] {
		if err := s.db.Exec(`UPDATE ai_models SET custom_model_name = '' WHERE custom_model_name IS NULL`).Error; err != nil {
			return fmt.Errorf("failed to patch ai_models.custom_model_name: %w", err)
		}
	}
	if colExists["created_at"] {
		if err := s.db.Exec(`UPDATE ai_models SET created_at = CURRENT_TIMESTAMP WHERE created_at IS NULL`).Error; err != nil {
			return fmt.Errorf("failed to patch ai_models.created_at: %w", err)
		}
	}
	if colExists["updated_at"] {
		if err := s.db.Exec(`UPDATE ai_models SET updated_at = CURRENT_TIMESTAMP WHERE updated_at IS NULL`).Error; err != nil {
			return fmt.Errorf("failed to patch ai_models.updated_at: %w", err)
		}
	}

	// Ensure basic index exists on user_id for lookup performance.
	if colExists["user_id"] {
		if err := s.db.Exec("CREATE INDEX IF NOT EXISTS idx_ai_models_user_id ON ai_models(user_id)").Error; err != nil {
			return fmt.Errorf("failed to create index idx_ai_models_user_id on ai_models(user_id): %w", err)
		}
	}

	if changed {
		logger.Warnf("⚠️ ai_models 表已执行 SQLite 兼容迁移（补列/修复数据/索引）")
	}

	return nil
}

func (s *AIModelStore) initDefaultData() error {
	// No longer pre-populate AI models - create on demand when user configures
	return nil
}

// List retrieves user's AI model list
func (s *AIModelStore) List(userID string) ([]*AIModel, error) {
	var models []*AIModel
	err := s.db.Where("user_id = ?", userID).Order("id").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

// Get retrieves a single AI model
func (s *AIModelStore) Get(userID, modelID string) (*AIModel, error) {
	if modelID == "" {
		return nil, fmt.Errorf("model ID cannot be empty")
	}

	candidates := []string{}
	if userID != "" {
		candidates = append(candidates, userID)
	}
	if userID != "default" {
		candidates = append(candidates, "default")
	}
	if len(candidates) == 0 {
		candidates = append(candidates, "default")
	}

	for _, uid := range candidates {
		var model AIModel
		err := s.db.Where("user_id = ? AND id = ?", uid, modelID).First(&model).Error
		if err == nil {
			return &model, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	return nil, gorm.ErrRecordNotFound
}

// GetByID retrieves an AI model by ID only (for debate engine)
func (s *AIModelStore) GetByID(modelID string) (*AIModel, error) {
	if modelID == "" {
		return nil, fmt.Errorf("model ID cannot be empty")
	}

	var model AIModel
	err := s.db.Where("id = ?", modelID).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// GetDefault retrieves the default enabled AI model
func (s *AIModelStore) GetDefault(userID string) (*AIModel, error) {
	if userID == "" {
		userID = "default"
	}
	model, err := s.firstEnabled(userID)
	if err == nil {
		return model, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if userID != "default" {
		return s.firstEnabled("default")
	}
	return nil, fmt.Errorf("please configure an available AI model in the system first")
}

func (s *AIModelStore) firstEnabled(userID string) (*AIModel, error) {
	var model AIModel
	err := s.db.Where("user_id = ? AND enabled = ?", userID, true).
		Order("updated_at DESC, id ASC").
		First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// Update updates AI model, creates if not exists
// IMPORTANT: If apiKey is empty string, the existing API key will be preserved (not overwritten)
func (s *AIModelStore) Update(userID, id string, enabled bool, apiKey, customAPIURL, customModelName string) error {
	// Try exact ID match first
	var existingModel AIModel
	err := s.db.Where("user_id = ? AND id = ?", userID, id).First(&existingModel).Error
	if err == nil {
		// Update existing model
		updates := map[string]interface{}{
			"enabled":           enabled,
			"custom_api_url":    customAPIURL,
			"custom_model_name": customModelName,
			"updated_at":        time.Now().UTC(),
		}
		// If apiKey is not empty, update it (encryption handled by crypto.EncryptedString)
		if apiKey != "" {
			updates["api_key"] = crypto.EncryptedString(apiKey)
		}
		return s.db.Model(&existingModel).Updates(updates).Error
	}

	// Try legacy logic compatibility: use id as provider to search
	provider := id
	err = s.db.Where("user_id = ? AND provider = ?", userID, provider).First(&existingModel).Error
	if err == nil {
		logger.Warnf("⚠️ Using legacy provider matching to update model: %s -> %s", provider, existingModel.ID)
		updates := map[string]interface{}{
			"enabled":           enabled,
			"custom_api_url":    customAPIURL,
			"custom_model_name": customModelName,
			"updated_at":        time.Now().UTC(),
		}
		if apiKey != "" {
			updates["api_key"] = crypto.EncryptedString(apiKey)
		}
		return s.db.Model(&existingModel).Updates(updates).Error
	}

	// Create new record
	if provider == id && (provider == "deepseek" || provider == "qwen") {
		provider = id
	} else {
		parts := strings.Split(id, "_")
		if len(parts) >= 2 {
			provider = parts[len(parts)-1]
		} else {
			provider = id
		}
	}

	// Try to get name from existing model with same provider
	var refModel AIModel
	var name string
	if err := s.db.Where("provider = ?", provider).First(&refModel).Error; err == nil {
		name = refModel.Name
	} else {
		if provider == "deepseek" {
			name = "DeepSeek AI"
		} else if provider == "qwen" {
			name = "Qwen AI"
		} else {
			name = provider + " AI"
		}
	}

	newModelID := id
	if id == provider {
		newModelID = fmt.Sprintf("%s_%s", userID, provider)
	}

	logger.Infof("✓ Creating new AI model configuration: ID=%s, Provider=%s, Name=%s", newModelID, provider, name)
	newModel := &AIModel{
		ID:              newModelID,
		UserID:          userID,
		Name:            name,
		Provider:        provider,
		Enabled:         enabled,
		APIKey:          crypto.EncryptedString(apiKey),
		CustomAPIURL:    customAPIURL,
		CustomModelName: customModelName,
	}
	return s.db.Create(newModel).Error
}

// Create creates an AI model
func (s *AIModelStore) Create(userID, id, name, provider string, enabled bool, apiKey, customAPIURL string) error {
	model := &AIModel{
		ID:           id,
		UserID:       userID,
		Name:         name,
		Provider:     provider,
		Enabled:      enabled,
		APIKey:       crypto.EncryptedString(apiKey),
		CustomAPIURL: customAPIURL,
	}
	// Use FirstOrCreate to ignore if already exists
	return s.db.Where("id = ?", id).FirstOrCreate(model).Error
}

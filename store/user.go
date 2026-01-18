package store

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
	"time"

	"nofx/logger"

	"gorm.io/gorm"
)

// UserStore user storage
type UserStore struct {
	db *gorm.DB
}

// User user model
type User struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex:idx_users_email;not null" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;not null" json:"-"`
	OTPSecret    string    `gorm:"column:otp_secret" json:"-"`
	OTPVerified  bool      `gorm:"column:otp_verified;default:false" json:"otp_verified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }

// GenerateOTPSecret generates OTP secret
func GenerateOTPSecret() (string, error) {
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(secret), nil
}

// NewUserStore creates a new UserStore
func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) initTables() error {
	// For PostgreSQL with existing table, skip AutoMigrate to avoid index conflicts
	if s.db.Dialector.Name() == "postgres" {
		var tableExists int64
		s.db.Raw(`SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'users'`).Scan(&tableExists)

		if tableExists > 0 {
			// Table exists - manually ensure all columns exist
			// Core columns (should already exist)
			s.db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT NOT NULL DEFAULT ''`)
			s.db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT NOT NULL DEFAULT ''`)
			s.db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`)
			s.db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`)
			// OTP columns (added later)
			s.db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS otp_secret TEXT DEFAULT ''`)
			s.db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS otp_verified BOOLEAN DEFAULT FALSE`)

			// Ensure unique index exists on email (don't care about the name)
			var indexExists int64
			s.db.Raw(`
				SELECT COUNT(*) FROM pg_indexes
				WHERE tablename = 'users' AND indexdef LIKE '%email%' AND indexdef LIKE '%UNIQUE%'
			`).Scan(&indexExists)

			if indexExists == 0 {
				s.db.Exec("CREATE UNIQUE INDEX idx_users_email ON users(email)")
			}

			return nil
		}
	}

	// SQLite: Some legacy DBs predate the Email NOT NULL constraint, and GORM's
	// AutoMigrate may fail on ALTER TABLE (it rebuilds table via users__temp).
	// For existing users table, do a safe incremental migration instead.
	if s.db.Dialector.Name() == "sqlite" && s.db.Migrator().HasTable(&User{}) {
		return s.ensureSQLiteUsersCompatibility()
	}

	return s.db.AutoMigrate(&User{})
}

type sqliteTableInfo struct {
	Name string `gorm:"column:name"`
}

type sqliteIndexList struct {
	Name   string `gorm:"column:name"`
	Unique int    `gorm:"column:unique"`
}

type sqliteIndexInfo struct {
	Name string `gorm:"column:name"`
}

func (s *UserStore) ensureSQLiteUsersCompatibility() error {
	var cols []sqliteTableInfo
	if err := s.db.Raw("PRAGMA table_info(users)").Scan(&cols).Error; err != nil {
		return fmt.Errorf("failed to inspect users table: %w", err)
	}

	colExists := make(map[string]bool, len(cols))
	for _, c := range cols {
		colExists[strings.ToLower(strings.TrimSpace(c.Name))] = true
	}

	changed := false

	// Add missing columns (keep them nullable here; we enforce via data fix + unique index).
	// SQLite cannot easily ALTER COLUMN to add NOT NULL without rebuild.
	if !colExists["email"] {
		if err := s.db.Exec("ALTER TABLE users ADD COLUMN email TEXT").Error; err != nil {
			return fmt.Errorf("failed to add users.email column: %w", err)
		}
		colExists["email"] = true
		changed = true
	}
	if !colExists["password_hash"] {
		if err := s.db.Exec("ALTER TABLE users ADD COLUMN password_hash TEXT").Error; err != nil {
			return fmt.Errorf("failed to add users.password_hash column: %w", err)
		}
		colExists["password_hash"] = true
		changed = true
	}
	if !colExists["created_at"] {
		if err := s.db.Exec("ALTER TABLE users ADD COLUMN created_at DATETIME DEFAULT CURRENT_TIMESTAMP").Error; err != nil {
			return fmt.Errorf("failed to add users.created_at column: %w", err)
		}
		colExists["created_at"] = true
		changed = true
	}
	if !colExists["updated_at"] {
		if err := s.db.Exec("ALTER TABLE users ADD COLUMN updated_at DATETIME DEFAULT CURRENT_TIMESTAMP").Error; err != nil {
			return fmt.Errorf("failed to add users.updated_at column: %w", err)
		}
		colExists["updated_at"] = true
		changed = true
	}
	if !colExists["otp_secret"] {
		if err := s.db.Exec("ALTER TABLE users ADD COLUMN otp_secret TEXT").Error; err != nil {
			return fmt.Errorf("failed to add users.otp_secret column: %w", err)
		}
		colExists["otp_secret"] = true
		changed = true
	}
	if !colExists["otp_verified"] {
		if err := s.db.Exec("ALTER TABLE users ADD COLUMN otp_verified BOOLEAN DEFAULT 0").Error; err != nil {
			return fmt.Errorf("failed to add users.otp_verified column: %w", err)
		}
		colExists["otp_verified"] = true
		changed = true
	}

	// Data fix: ensure legacy rows won't violate NOT NULL / UNIQUE expectations.
	if colExists["email"] {
		// Prefer stable admin email if admin row exists.
		if colExists["id"] {
			if err := s.db.Exec(`UPDATE users
SET email = 'admin@localhost'
WHERE (email IS NULL OR trim(email) = '') AND id = 'admin'`).Error; err != nil {
				return fmt.Errorf("failed to patch users.email for admin: %w", err)
			}
		}

		// Fill blank email for legacy rows with deterministic unique placeholders.
		fillSQL := `UPDATE users
SET email = 'legacy_' || CAST(rowid AS TEXT) || '@local'
WHERE email IS NULL OR trim(email) = ''`
		if colExists["id"] {
			fillSQL = `UPDATE users
SET email = 'legacy_' || COALESCE(NULLIF(trim(id), ''), CAST(rowid AS TEXT)) || '@local'
WHERE email IS NULL OR trim(email) = ''`
		}
		if err := s.db.Exec(fillSQL).Error; err != nil {
			return fmt.Errorf("failed to patch legacy users.email: %w", err)
		}

		// De-duplicate any existing duplicated emails (old DBs might not enforce UNIQUE).
		var dupEmails []string
		if err := s.db.Raw(`SELECT email FROM users GROUP BY email HAVING COUNT(*) > 1`).Scan(&dupEmails).Error; err != nil {
			return fmt.Errorf("failed to detect duplicated users.email: %w", err)
		}
		for _, email := range dupEmails {
			email = strings.TrimSpace(email)
			if email == "" {
				continue
			}
			if err := s.db.Exec(`UPDATE users
SET email = email || '+' || CAST(rowid AS TEXT)
WHERE email = ? AND rowid != (SELECT MIN(rowid) FROM users WHERE email = ?)`, email, email).Error; err != nil {
				return fmt.Errorf("failed to de-duplicate users.email (%s): %w", email, err)
			}
		}
	}

	if colExists["password_hash"] {
		if err := s.db.Exec(`UPDATE users SET password_hash = '' WHERE password_hash IS NULL`).Error; err != nil {
			return fmt.Errorf("failed to patch users.password_hash: %w", err)
		}
	}
	if colExists["otp_secret"] {
		if err := s.db.Exec(`UPDATE users SET otp_secret = '' WHERE otp_secret IS NULL`).Error; err != nil {
			return fmt.Errorf("failed to patch users.otp_secret: %w", err)
		}
	}
	if colExists["otp_verified"] {
		if err := s.db.Exec(`UPDATE users SET otp_verified = 0 WHERE otp_verified IS NULL`).Error; err != nil {
			return fmt.Errorf("failed to patch users.otp_verified: %w", err)
		}
	}
	if colExists["created_at"] {
		if err := s.db.Exec(`UPDATE users SET created_at = CURRENT_TIMESTAMP WHERE created_at IS NULL`).Error; err != nil {
			return fmt.Errorf("failed to patch users.created_at: %w", err)
		}
	}
	if colExists["updated_at"] {
		if err := s.db.Exec(`UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE updated_at IS NULL`).Error; err != nil {
			return fmt.Errorf("failed to patch users.updated_at: %w", err)
		}
	}

	// Ensure unique index exists on email (don't care about the exact name, but keep idx_users_email for consistency).
	if colExists["email"] {
		uniqueEmailIndexFound := false
		var idxList []sqliteIndexList
		if err := s.db.Raw("PRAGMA index_list(users)").Scan(&idxList).Error; err != nil {
			return fmt.Errorf("failed to list users indexes: %w", err)
		}
		for _, idx := range idxList {
			if idx.Unique != 1 {
				continue
			}
			idxName := strings.TrimSpace(idx.Name)
			if idxName == "" {
				continue
			}
			var idxCols []sqliteIndexInfo
			quoted := `"` + strings.ReplaceAll(idxName, `"`, `""`) + `"`
			if err := s.db.Raw("PRAGMA index_info(" + quoted + ")").Scan(&idxCols).Error; err != nil {
				continue
			}
			if len(idxCols) == 1 && strings.EqualFold(strings.TrimSpace(idxCols[0].Name), "email") {
				uniqueEmailIndexFound = true
				break
			}
		}
		if !uniqueEmailIndexFound {
			if err := s.db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
				return fmt.Errorf("failed to create unique index idx_users_email on users(email): %w", err)
			}
			changed = true
		}
	}

	if changed {
		logger.Warnf("⚠️ users 表已执行 SQLite 兼容迁移（补列/修复数据/索引）")
	}

	return nil
}

// Create creates user
func (s *UserStore) Create(user *User) error {
	return s.db.Create(user).Error
}

// GetByEmail gets user by email
func (s *UserStore) GetByEmail(email string) (*User, error) {
	var user User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID gets user by ID
func (s *UserStore) GetByID(userID string) (*User, error) {
	var user User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Count returns the total number of users
func (s *UserStore) Count() (int, error) {
	var count int64
	err := s.db.Model(&User{}).Count(&count).Error
	return int(count), err
}

// GetAllIDs gets all user IDs
func (s *UserStore) GetAllIDs() ([]string, error) {
	var userIDs []string
	err := s.db.Model(&User{}).Order("id").Pluck("id", &userIDs).Error
	return userIDs, err
}

// UpdateOTPVerified updates OTP verification status
func (s *UserStore) UpdateOTPVerified(userID string, verified bool) error {
	return s.db.Model(&User{}).Where("id = ?", userID).Update("otp_verified", verified).Error
}

// UpdatePassword updates password
func (s *UserStore) UpdatePassword(userID, passwordHash string) error {
	return s.db.Model(&User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password_hash": passwordHash,
		"updated_at":    time.Now().UTC(),
	}).Error
}

// EnsureAdmin ensures admin user exists
func (s *UserStore) EnsureAdmin() error {
	var count int64
	s.db.Model(&User{}).Where("id = ?", "admin").Count(&count)
	if count > 0 {
		return nil
	}
	return s.Create(&User{
		ID:           "admin",
		Email:        "admin@localhost",
		PasswordHash: "",
		OTPSecret:    "",
		OTPVerified:  true,
	})
}

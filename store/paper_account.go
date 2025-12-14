package store

import (
	"database/sql"
	"fmt"
	"nofx/logger"
	"time"

	"github.com/google/uuid"
)

// PaperAccount Paper 模拟账户
type PaperAccount struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	ExchangeID      string    `json:"exchange_id"`       // paper exchange UUID
	InitialBalance  float64   `json:"initial_balance"`   // 初始资金
	CurrentBalance  float64   `json:"current_balance"`   // 当前可用余额
	TotalMarginUsed float64   `json:"total_margin_used"` // 已用保证金
	TotalPnL        float64   `json:"total_pnl"`         // 总盈亏
	SlippageRate    float64   `json:"slippage_rate"`     // 滑点率 (默认 0.0005)
	FeeRate         float64   `json:"fee_rate"`          // 手续费率 (默认 0.0004)
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PaperPosition Paper 虚拟持仓
type PaperPosition struct {
	ID               int64     `json:"id"`
	AccountID        string    `json:"account_id"`        // PaperAccount.ID
	Symbol           string    `json:"symbol"`
	Side             string    `json:"side"`              // "LONG" or "SHORT"
	Quantity         float64   `json:"quantity"`
	EntryPrice       float64   `json:"entry_price"`
	MarkPrice        float64   `json:"mark_price"`
	Leverage         int       `json:"leverage"`
	MarginUsed       float64   `json:"margin_used"`
	UnrealizedPnL    float64   `json:"unrealized_pnl"`
	LiquidationPrice float64   `json:"liquidation_price"`
	StopLoss         float64   `json:"stop_loss"`
	TakeProfit       float64   `json:"take_profit"`
	OpenTime         time.Time `json:"open_time"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// PaperOrder Paper 订单记录
type PaperOrder struct {
	ID           int64      `json:"id"`
	AccountID    string     `json:"account_id"`
	Symbol       string     `json:"symbol"`
	Side         string     `json:"side"`          // "BUY" or "SELL"
	PositionSide string     `json:"position_side"` // "LONG" or "SHORT"
	OrderType    string     `json:"order_type"`    // "MARKET", "LIMIT", "STOP_LOSS", "TAKE_PROFIT"
	Quantity     float64    `json:"quantity"`
	Price        float64    `json:"price"`
	AvgPrice     float64    `json:"avg_price"`
	Fee          float64    `json:"fee"`
	Status       string     `json:"status"`       // "NEW", "FILLED", "CANCELED"
	RealizedPnL  float64    `json:"realized_pnl"` // 平仓时的已实现盈亏
	CreatedAt    time.Time  `json:"created_at"`
	FilledAt     *time.Time `json:"filled_at"`
}

// PaperAccountStore Paper 账户存储
type PaperAccountStore struct {
	db *sql.DB
}

// NewPaperAccountStore creates paper account storage instance
func NewPaperAccountStore(db *sql.DB) *PaperAccountStore {
	return &PaperAccountStore{db: db}
}

// InitTables initializes paper trading tables
func (s *PaperAccountStore) InitTables() error {
	// Create paper_accounts table
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS paper_accounts (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			exchange_id TEXT NOT NULL UNIQUE,
			initial_balance REAL NOT NULL DEFAULT 10000,
			current_balance REAL NOT NULL DEFAULT 10000,
			total_margin_used REAL DEFAULT 0,
			total_pnl REAL DEFAULT 0,
			slippage_rate REAL DEFAULT 0.0005,
			fee_rate REAL DEFAULT 0.0004,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create paper_accounts table: %w", err)
	}

	// Create paper_positions table
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS paper_positions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id TEXT NOT NULL,
			symbol TEXT NOT NULL,
			side TEXT NOT NULL,
			quantity REAL NOT NULL,
			entry_price REAL NOT NULL,
			mark_price REAL DEFAULT 0,
			leverage INTEGER DEFAULT 10,
			margin_used REAL DEFAULT 0,
			unrealized_pnl REAL DEFAULT 0,
			liquidation_price REAL DEFAULT 0,
			stop_loss REAL DEFAULT 0,
			take_profit REAL DEFAULT 0,
			open_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(account_id, symbol, side)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create paper_positions table: %w", err)
	}

	// Create paper_orders table
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS paper_orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id TEXT NOT NULL,
			symbol TEXT NOT NULL,
			side TEXT NOT NULL,
			position_side TEXT NOT NULL,
			order_type TEXT NOT NULL,
			quantity REAL NOT NULL,
			price REAL DEFAULT 0,
			avg_price REAL DEFAULT 0,
			fee REAL DEFAULT 0,
			status TEXT DEFAULT 'FILLED',
			realized_pnl REAL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			filled_at DATETIME
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create paper_orders table: %w", err)
	}

	// Create indexes
	indices := []string{
		`CREATE INDEX IF NOT EXISTS idx_paper_positions_account ON paper_positions(account_id)`,
		`CREATE INDEX IF NOT EXISTS idx_paper_orders_account ON paper_orders(account_id)`,
	}
	for _, idx := range indices {
		if _, err := s.db.Exec(idx); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// =============================================================================
// Account Methods
// =============================================================================

// Create creates a new paper account
func (s *PaperAccountStore) Create(account *PaperAccount) error {
	now := time.Now()
	if account.ID == "" {
		account.ID = uuid.New().String()
	}
	account.CreatedAt = now
	account.UpdatedAt = now

	// Set default values
	if account.SlippageRate == 0 {
		account.SlippageRate = 0.0005
	}
	if account.FeeRate == 0 {
		account.FeeRate = 0.0004
	}
	if account.InitialBalance == 0 {
		account.InitialBalance = 10000
	}
	if account.CurrentBalance == 0 {
		account.CurrentBalance = account.InitialBalance
	}

	_, err := s.db.Exec(`
		INSERT INTO paper_accounts (
			id, user_id, exchange_id, initial_balance, current_balance,
			total_margin_used, total_pnl, slippage_rate, fee_rate,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		account.ID, account.UserID, account.ExchangeID,
		account.InitialBalance, account.CurrentBalance,
		account.TotalMarginUsed, account.TotalPnL,
		account.SlippageRate, account.FeeRate,
		now.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return fmt.Errorf("failed to create paper account: %w", err)
	}

	return nil
}

// Get gets paper account by user ID and exchange ID
func (s *PaperAccountStore) Get(userID, exchangeID string) (*PaperAccount, error) {
	var account PaperAccount
	var createdAt, updatedAt string

	err := s.db.QueryRow(`
		SELECT id, user_id, exchange_id, initial_balance, current_balance,
		       total_margin_used, total_pnl, slippage_rate, fee_rate,
		       created_at, updated_at
		FROM paper_accounts
		WHERE user_id = ? AND exchange_id = ?
	`, userID, exchangeID).Scan(
		&account.ID, &account.UserID, &account.ExchangeID,
		&account.InitialBalance, &account.CurrentBalance,
		&account.TotalMarginUsed, &account.TotalPnL,
		&account.SlippageRate, &account.FeeRate,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get paper account: %w", err)
	}

	account.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	account.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &account, nil
}

// GetByID gets paper account by ID
func (s *PaperAccountStore) GetByID(accountID string) (*PaperAccount, error) {
	var account PaperAccount
	var createdAt, updatedAt string

	err := s.db.QueryRow(`
		SELECT id, user_id, exchange_id, initial_balance, current_balance,
		       total_margin_used, total_pnl, slippage_rate, fee_rate,
		       created_at, updated_at
		FROM paper_accounts
		WHERE id = ?
	`, accountID).Scan(
		&account.ID, &account.UserID, &account.ExchangeID,
		&account.InitialBalance, &account.CurrentBalance,
		&account.TotalMarginUsed, &account.TotalPnL,
		&account.SlippageRate, &account.FeeRate,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get paper account by ID: %w", err)
	}

	account.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	account.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &account, nil
}

// Update updates paper account
func (s *PaperAccountStore) Update(account *PaperAccount) error {
	now := time.Now()
	account.UpdatedAt = now

	result, err := s.db.Exec(`
		UPDATE paper_accounts SET
			current_balance = ?,
			total_margin_used = ?,
			total_pnl = ?,
			slippage_rate = ?,
			fee_rate = ?,
			updated_at = ?
		WHERE id = ?
	`,
		account.CurrentBalance, account.TotalMarginUsed,
		account.TotalPnL, account.SlippageRate, account.FeeRate,
		now.Format("2006-01-02 15:04:05"), account.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update paper account: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("paper account not found: id=%s", account.ID)
	}

	return nil
}

// UpdateTx updates paper account within a transaction
func (s *PaperAccountStore) UpdateTx(tx *sql.Tx, account *PaperAccount) error {
	now := time.Now()
	account.UpdatedAt = now

	result, err := tx.Exec(`
		UPDATE paper_accounts SET
			current_balance = ?,
			total_margin_used = ?,
			total_pnl = ?,
			slippage_rate = ?,
			fee_rate = ?,
			updated_at = ?
		WHERE id = ?
	`,
		account.CurrentBalance, account.TotalMarginUsed,
		account.TotalPnL, account.SlippageRate, account.FeeRate,
		now.Format("2006-01-02 15:04:05"), account.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update paper account: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("paper account not found: id=%s", account.ID)
	}

	return nil
}

// Reset resets paper account to initial state with new balance
func (s *PaperAccountStore) Reset(userID, exchangeID string, newBalance float64) error {
	now := time.Now()

	// Get account first
	account, err := s.Get(userID, exchangeID)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("paper account not found: userID=%s, exchangeID=%s", userID, exchangeID)
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Delete all positions for this account
	_, err = tx.Exec(`DELETE FROM paper_positions WHERE account_id = ?`, account.ID)
	if err != nil {
		return fmt.Errorf("failed to delete positions during reset: %w", err)
	}

	// Delete all orders for this account
	_, err = tx.Exec(`DELETE FROM paper_orders WHERE account_id = ?`, account.ID)
	if err != nil {
		return fmt.Errorf("failed to delete orders during reset: %w", err)
	}

	// Reset account balance
	_, err = tx.Exec(`
		UPDATE paper_accounts SET
			initial_balance = ?,
			current_balance = ?,
			total_margin_used = 0,
			total_pnl = 0,
			updated_at = ?
		WHERE user_id = ? AND exchange_id = ?
	`,
		newBalance, newBalance, now.Format("2006-01-02 15:04:05"),
		userID, exchangeID,
	)
	if err != nil {
		return fmt.Errorf("failed to reset paper account: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// =============================================================================
// Position Methods
// =============================================================================

// CreatePosition creates a new paper position
func (s *PaperAccountStore) CreatePosition(pos *PaperPosition) error {
	now := time.Now()
	pos.OpenTime = now
	pos.UpdatedAt = now

	result, err := s.db.Exec(`
		INSERT INTO paper_positions (
			account_id, symbol, side, quantity, entry_price, mark_price,
			leverage, margin_used, unrealized_pnl, liquidation_price,
			stop_loss, take_profit, open_time, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		pos.AccountID, pos.Symbol, pos.Side, pos.Quantity, pos.EntryPrice, pos.MarkPrice,
		pos.Leverage, pos.MarginUsed, pos.UnrealizedPnL, pos.LiquidationPrice,
		pos.StopLoss, pos.TakeProfit,
		now.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return fmt.Errorf("failed to create paper position: %w", err)
	}

	id, _ := result.LastInsertId()
	pos.ID = id
	return nil
}

// CreatePositionTx creates a new paper position within a transaction
func (s *PaperAccountStore) CreatePositionTx(tx *sql.Tx, pos *PaperPosition) error {
	now := time.Now()
	pos.OpenTime = now
	pos.UpdatedAt = now

	result, err := tx.Exec(`
		INSERT INTO paper_positions (
			account_id, symbol, side, quantity, entry_price, mark_price,
			leverage, margin_used, unrealized_pnl, liquidation_price,
			stop_loss, take_profit, open_time, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		pos.AccountID, pos.Symbol, pos.Side, pos.Quantity, pos.EntryPrice, pos.MarkPrice,
		pos.Leverage, pos.MarginUsed, pos.UnrealizedPnL, pos.LiquidationPrice,
		pos.StopLoss, pos.TakeProfit,
		now.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return fmt.Errorf("failed to create paper position: %w", err)
	}

	id, _ := result.LastInsertId()
	pos.ID = id
	return nil
}

// GetPosition gets a paper position by account ID, symbol and side
func (s *PaperAccountStore) GetPosition(accountID, symbol, side string) (*PaperPosition, error) {
	var pos PaperPosition
	var openTime, updatedAt string

	err := s.db.QueryRow(`
		SELECT id, account_id, symbol, side, quantity, entry_price, mark_price,
		       leverage, margin_used, unrealized_pnl, liquidation_price,
		       stop_loss, take_profit, open_time, updated_at
		FROM paper_positions
		WHERE account_id = ? AND symbol = ? AND side = ?
	`, accountID, symbol, side).Scan(
		&pos.ID, &pos.AccountID, &pos.Symbol, &pos.Side, &pos.Quantity,
		&pos.EntryPrice, &pos.MarkPrice, &pos.Leverage, &pos.MarginUsed,
		&pos.UnrealizedPnL, &pos.LiquidationPrice, &pos.StopLoss, &pos.TakeProfit,
		&openTime, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get paper position: %w", err)
	}

	pos.OpenTime, _ = time.Parse("2006-01-02 15:04:05", openTime)
	pos.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &pos, nil
}

// ListPositions lists all positions for an account
func (s *PaperAccountStore) ListPositions(accountID string) ([]*PaperPosition, error) {
	rows, err := s.db.Query(`
		SELECT id, account_id, symbol, side, quantity, entry_price, mark_price,
		       leverage, margin_used, unrealized_pnl, liquidation_price,
		       stop_loss, take_profit, open_time, updated_at
		FROM paper_positions
		WHERE account_id = ?
		ORDER BY open_time DESC
	`, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to list paper positions: %w", err)
	}
	defer rows.Close()

	var positions []*PaperPosition
	for rows.Next() {
		var pos PaperPosition
		var openTime, updatedAt string

		err := rows.Scan(
			&pos.ID, &pos.AccountID, &pos.Symbol, &pos.Side, &pos.Quantity,
			&pos.EntryPrice, &pos.MarkPrice, &pos.Leverage, &pos.MarginUsed,
			&pos.UnrealizedPnL, &pos.LiquidationPrice, &pos.StopLoss, &pos.TakeProfit,
			&openTime, &updatedAt,
		)
		if err != nil {
			continue
		}

		pos.OpenTime, _ = time.Parse("2006-01-02 15:04:05", openTime)
		pos.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		positions = append(positions, &pos)
	}

	return positions, nil
}

// UpdatePosition updates a paper position
func (s *PaperAccountStore) UpdatePosition(pos *PaperPosition) error {
	now := time.Now()
	pos.UpdatedAt = now

	result, err := s.db.Exec(`
		UPDATE paper_positions SET
			quantity = ?,
			entry_price = ?,
			mark_price = ?,
			leverage = ?,
			margin_used = ?,
			unrealized_pnl = ?,
			liquidation_price = ?,
			stop_loss = ?,
			take_profit = ?,
			updated_at = ?
		WHERE id = ?
	`,
		pos.Quantity, pos.EntryPrice, pos.MarkPrice, pos.Leverage,
		pos.MarginUsed, pos.UnrealizedPnL, pos.LiquidationPrice,
		pos.StopLoss, pos.TakeProfit,
		now.Format("2006-01-02 15:04:05"), pos.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update paper position: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("paper position not found: id=%d", pos.ID)
	}

	return nil
}

// UpdatePositionTx updates a paper position within a transaction
func (s *PaperAccountStore) UpdatePositionTx(tx *sql.Tx, pos *PaperPosition) error {
	now := time.Now()
	pos.UpdatedAt = now

	result, err := tx.Exec(`
		UPDATE paper_positions SET
			quantity = ?,
			entry_price = ?,
			mark_price = ?,
			leverage = ?,
			margin_used = ?,
			unrealized_pnl = ?,
			liquidation_price = ?,
			stop_loss = ?,
			take_profit = ?,
			updated_at = ?
		WHERE id = ?
	`,
		pos.Quantity, pos.EntryPrice, pos.MarkPrice, pos.Leverage,
		pos.MarginUsed, pos.UnrealizedPnL, pos.LiquidationPrice,
		pos.StopLoss, pos.TakeProfit,
		now.Format("2006-01-02 15:04:05"), pos.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update paper position: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("paper position not found: id=%d", pos.ID)
	}

	return nil
}

// DeletePosition deletes a paper position
func (s *PaperAccountStore) DeletePosition(posID int64) error {
	result, err := s.db.Exec(`DELETE FROM paper_positions WHERE id = ?`, posID)
	if err != nil {
		return fmt.Errorf("failed to delete paper position: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("paper position not found: id=%d", posID)
	}

	return nil
}

// DeletePositionTx deletes a paper position within a transaction
func (s *PaperAccountStore) DeletePositionTx(tx *sql.Tx, posID int64) error {
	result, err := tx.Exec(`DELETE FROM paper_positions WHERE id = ?`, posID)
	if err != nil {
		return fmt.Errorf("failed to delete paper position: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("paper position not found: id=%d", posID)
	}

	return nil
}

// =============================================================================
// Order Methods
// =============================================================================

// CreateOrder creates a new paper order
func (s *PaperAccountStore) CreateOrder(order *PaperOrder) error {
	now := time.Now()
	order.CreatedAt = now

	result, err := s.db.Exec(`
		INSERT INTO paper_orders (
			account_id, symbol, side, position_side, order_type,
			quantity, price, avg_price, fee, status, realized_pnl,
			created_at, filled_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		order.AccountID, order.Symbol, order.Side, order.PositionSide, order.OrderType,
		order.Quantity, order.Price, order.AvgPrice, order.Fee, order.Status, order.RealizedPnL,
		now.Format("2006-01-02 15:04:05"), s.formatNullableTime(order.FilledAt),
	)
	if err != nil {
		return fmt.Errorf("failed to create paper order: %w", err)
	}

	id, _ := result.LastInsertId()
	order.ID = id
	return nil
}

// CreateOrderTx creates a new paper order within a transaction
func (s *PaperAccountStore) CreateOrderTx(tx *sql.Tx, order *PaperOrder) error {
	now := time.Now()
	order.CreatedAt = now

	result, err := tx.Exec(`
		INSERT INTO paper_orders (
			account_id, symbol, side, position_side, order_type,
			quantity, price, avg_price, fee, status, realized_pnl,
			created_at, filled_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		order.AccountID, order.Symbol, order.Side, order.PositionSide, order.OrderType,
		order.Quantity, order.Price, order.AvgPrice, order.Fee, order.Status, order.RealizedPnL,
		now.Format("2006-01-02 15:04:05"), s.formatNullableTime(order.FilledAt),
	)
	if err != nil {
		return fmt.Errorf("failed to create paper order: %w", err)
	}

	id, _ := result.LastInsertId()
	order.ID = id
	return nil
}

// GetOrder gets a paper order by account ID and order ID
func (s *PaperAccountStore) GetOrder(accountID string, orderID int64) (*PaperOrder, error) {
	var order PaperOrder
	var createdAt string
	var filledAt sql.NullString

	err := s.db.QueryRow(`
		SELECT id, account_id, symbol, side, position_side, order_type,
		       quantity, price, avg_price, fee, status, realized_pnl,
		       created_at, filled_at
		FROM paper_orders
		WHERE account_id = ? AND id = ?
	`, accountID, orderID).Scan(
		&order.ID, &order.AccountID, &order.Symbol, &order.Side,
		&order.PositionSide, &order.OrderType, &order.Quantity,
		&order.Price, &order.AvgPrice, &order.Fee, &order.Status,
		&order.RealizedPnL, &createdAt, &filledAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get paper order: %w", err)
	}

	order.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	if filledAt.Valid {
		t, _ := time.Parse("2006-01-02 15:04:05", filledAt.String)
		order.FilledAt = &t
	}

	return &order, nil
}

// ListOrders lists orders for an account with limit
func (s *PaperAccountStore) ListOrders(accountID string, limit int) ([]*PaperOrder, error) {
	rows, err := s.db.Query(`
		SELECT id, account_id, symbol, side, position_side, order_type,
		       quantity, price, avg_price, fee, status, realized_pnl,
		       created_at, filled_at
		FROM paper_orders
		WHERE account_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, accountID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list paper orders: %w", err)
	}
	defer rows.Close()

	var orders []*PaperOrder
	for rows.Next() {
		var order PaperOrder
		var createdAt string
		var filledAt sql.NullString

		err := rows.Scan(
			&order.ID, &order.AccountID, &order.Symbol, &order.Side,
			&order.PositionSide, &order.OrderType, &order.Quantity,
			&order.Price, &order.AvgPrice, &order.Fee, &order.Status,
			&order.RealizedPnL, &createdAt, &filledAt,
		)
		if err != nil {
			continue
		}

		order.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		if filledAt.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", filledAt.String)
			order.FilledAt = &t
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

// UpdateOrderStatus updates order status, average price and filled time
func (s *PaperAccountStore) UpdateOrderStatus(orderID int64, status string, avgPrice float64, filledAt time.Time) error {
	result, err := s.db.Exec(`
		UPDATE paper_orders SET
			status = ?,
			avg_price = ?,
			filled_at = ?
		WHERE id = ?
	`,
		status, avgPrice, filledAt.Format("2006-01-02 15:04:05"), orderID,
	)
	if err != nil {
		return fmt.Errorf("failed to update paper order status: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("paper order not found: id=%d", orderID)
	}

	return nil
}

// =============================================================================
// Helper Methods
// =============================================================================

// formatNullableTime formats a nullable time pointer to string or nil
func (s *PaperAccountStore) formatNullableTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format("2006-01-02 15:04:05")
}

// DeleteByExchangeID deletes paper account and all related data by exchange ID
func (s *PaperAccountStore) DeleteByExchangeID(exchangeID string) error {
	// 首先获取账户 ID
	var accountID string
	err := s.db.QueryRow(`SELECT id FROM paper_accounts WHERE exchange_id = ?`, exchangeID).Scan(&accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			// 没有关联的 paper account，不需要清理
			return nil
		}
		return fmt.Errorf("failed to find paper account: %w", err)
	}

	// 删除所有关联的持仓
	_, err = s.db.Exec(`DELETE FROM paper_positions WHERE account_id = ?`, accountID)
	if err != nil {
		return fmt.Errorf("failed to delete paper positions: %w", err)
	}

	// 删除所有关联的订单
	_, err = s.db.Exec(`DELETE FROM paper_orders WHERE account_id = ?`, accountID)
	if err != nil {
		return fmt.Errorf("failed to delete paper orders: %w", err)
	}

	// 删除账户本身
	_, err = s.db.Exec(`DELETE FROM paper_accounts WHERE id = ?`, accountID)
	if err != nil {
		return fmt.Errorf("failed to delete paper account: %w", err)
	}

	logger.Infof("🗑️ Deleted paper account and related data: exchangeID=%s, accountID=%s", exchangeID, accountID)
	return nil
}

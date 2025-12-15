package trader

import (
	"database/sql"
	"fmt"
	"nofx/logger"
	"nofx/market"
	"nofx/store"
	"strings"
	"time"
)

// PaperTrader Paper 模拟交易实现
// 实现 Trader 接口，使用真实市场数据但虚拟账户
type PaperTrader struct {
	store         *store.Store
	userID        string
	traderID      string  // 交易员ID - 隔离键
	exchangeID    string  // Paper exchange UUID
	accountID     string  // Paper account ID
	slippageRate  float64 // 滑点率 (默认 0.0005 = 0.05%)
	feeRate       float64 // 手续费率 (默认 0.0004 = 0.04%)
	isCrossMargin bool
}

// 确保 PaperTrader 实现了 Trader 接口
var _ Trader = (*PaperTrader)(nil)

// NewPaperTrader 创建 Paper Trader 实例
// traderID 用于隔离不同交易员的虚拟账户，即使使用同一个 Paper Exchange 也是独立的
func NewPaperTrader(userID, traderID, exchangeID string, initialBalance float64, st *store.Store) (*PaperTrader, error) {
	// 检查或创建 Paper Account - 使用 traderID 作为隔离键
	account, err := st.PaperAccount().GetByTraderID(traderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get paper account: %w", err)
	}

	if account == nil {
		// 创建新账户 - 使用 traderID 作为隔离键
		account = &store.PaperAccount{
			ID:             fmt.Sprintf("paper_%s", traderID),
			UserID:         userID,
			TraderID:       traderID,
			ExchangeID:     exchangeID,
			InitialBalance: initialBalance,
			CurrentBalance: initialBalance,
			SlippageRate:   0.0005,
			FeeRate:        0.0004,
		}
		if err := st.PaperAccount().Create(account); err != nil {
			return nil, fmt.Errorf("failed to create paper account: %w", err)
		}
	} else {
		// 账户已存在，检查 initialBalance 是否不匹配
		if initialBalance != account.InitialBalance {
			logger.Warnf("⚠️ [Paper] Account already exists with initial balance %.2f, requested %.2f. Use Reset API to update.", account.InitialBalance, initialBalance)
		}
	}

	return &PaperTrader{
		store:        st,
		userID:       userID,
		traderID:     traderID,
		exchangeID:   exchangeID,
		accountID:    account.ID,
		slippageRate: account.SlippageRate,
		feeRate:      account.FeeRate,
	}, nil
}

// GetBalance 返回账户余额
func (t *PaperTrader) GetBalance() (map[string]interface{}, error) {
	account, err := t.store.PaperAccount().GetByID(t.accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("paper account not found: %s", t.accountID)
	}

	// 计算未实现盈亏 - 使用实时市场价格
	positions, _ := t.store.PaperAccount().ListPositions(t.accountID)
	unrealizedPnL := 0.0
	for _, pos := range positions {
		// 获取实时市场价格
		if md, err := market.Get(pos.Symbol); err == nil {
			// 实时计算未实现盈亏
			if pos.Side == "LONG" {
				unrealizedPnL += (md.CurrentPrice - pos.EntryPrice) * pos.Quantity
			} else {
				unrealizedPnL += (pos.EntryPrice - md.CurrentPrice) * pos.Quantity
			}
		} else {
			// 如果无法获取市场价格，回退使用存储的值
			unrealizedPnL += pos.UnrealizedPnL
		}
	}

	totalEquity := account.CurrentBalance + account.TotalMarginUsed + unrealizedPnL

	return map[string]interface{}{
		"totalWalletBalance":    account.CurrentBalance + account.TotalMarginUsed,
		"totalUnrealizedProfit": unrealizedPnL,
		"availableBalance":      account.CurrentBalance,
		"total_equity":          totalEquity, // 使用 snake_case 与系统标准一致
	}, nil
}

// GetPositions 返回所有持仓
// 同时检查止损/止盈触发条件，在每个交易周期自动检查
func (t *PaperTrader) GetPositions() ([]map[string]interface{}, error) {
	positions, err := t.store.PaperAccount().ListPositions(t.accountID)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(positions))
	for _, pos := range positions {
		// 1. 获取实时市场价格
		md, err := market.Get(pos.Symbol)
		if err != nil {
			// 无法获取价格，跳过检查，使用存储的值
			posAmt := pos.Quantity
			if pos.Side == "SHORT" {
				posAmt = -pos.Quantity
			}
			result = append(result, map[string]interface{}{
				"symbol":           pos.Symbol,
				"side":             strings.ToLower(pos.Side),
				"positionAmt":      posAmt,
				"entryPrice":       pos.EntryPrice,
				"markPrice":        pos.MarkPrice,
				"unRealizedProfit": pos.UnrealizedPnL,
				"leverage":         float64(pos.Leverage),
				"liquidationPrice": pos.LiquidationPrice,
			})
			continue
		}

		currentPrice := md.CurrentPrice

		// 2. 首先检查强平（优先级最高）
		// LONG: 市场价 ≤ 强平价 → 触发（价格暴跌导致爆仓）
		// SHORT: 市场价 ≥ 强平价 → 触发（价格暴涨导致爆仓）
		if pos.LiquidationPrice > 0 {
			liquidated := false
			if pos.Side == "LONG" && currentPrice <= pos.LiquidationPrice {
				liquidated = true
			} else if pos.Side == "SHORT" && currentPrice >= pos.LiquidationPrice {
				liquidated = true
			}
			if liquidated {
				logger.Warnf("💀 [Paper] LIQUIDATION triggered for %s %s: price=%.4f, liqPrice=%.4f",
					pos.Side, pos.Symbol, currentPrice, pos.LiquidationPrice)
				// 强平：使用强平价格执行，而非当前市场价
				if err := t.executeAutoClose(pos, pos.LiquidationPrice, "LIQUIDATION"); err != nil {
					logger.Warnf("⚠️ [Paper] Auto-close failed, position still exists: %v", err)
					// 继续将此持仓添加到结果（不 continue）
				} else {
					continue // 仓位已成功平仓，不添加到结果
				}
			}
		}

		// 3. 检查止损触发
		// LONG: 市场价 ≤ 止损价 → 触发（价格下跌到止损）
		// SHORT: 市场价 ≥ 止损价 → 触发（价格上涨到止损）
		if pos.StopLoss > 0 {
			triggered := false
			if pos.Side == "LONG" && currentPrice <= pos.StopLoss {
				triggered = true
			} else if pos.Side == "SHORT" && currentPrice >= pos.StopLoss {
				triggered = true
			}
			if triggered {
				logger.Infof("🛑 [Paper] Stop loss triggered for %s %s: price=%.4f, stopLoss=%.4f",
					pos.Side, pos.Symbol, currentPrice, pos.StopLoss)
				if err := t.executeAutoClose(pos, currentPrice, "STOP_LOSS"); err != nil {
					logger.Warnf("⚠️ [Paper] Auto-close failed, position still exists: %v", err)
					// 继续将此持仓添加到结果（不 continue）
				} else {
					continue // 仓位已成功平仓，不添加到结果
				}
			}
		}

		// 4. 检查止盈触发
		// LONG: 市场价 ≥ 止盈价 → 触发（价格上涨到止盈）
		// SHORT: 市场价 ≤ 止盈价 → 触发（价格下跌到止盈）
		if pos.TakeProfit > 0 {
			triggered := false
			if pos.Side == "LONG" && currentPrice >= pos.TakeProfit {
				triggered = true
			} else if pos.Side == "SHORT" && currentPrice <= pos.TakeProfit {
				triggered = true
			}
			if triggered {
				logger.Infof("🎯 [Paper] Take profit triggered for %s %s: price=%.4f, takeProfit=%.4f",
					pos.Side, pos.Symbol, currentPrice, pos.TakeProfit)
				if err := t.executeAutoClose(pos, currentPrice, "TAKE_PROFIT"); err != nil {
					logger.Warnf("⚠️ [Paper] Auto-close failed, position still exists: %v", err)
					// 继续将此持仓添加到结果（不 continue）
				} else {
					continue // 仓位已成功平仓，不添加到结果
				}
			}
		}

		// 5. 更新 mark price 和未实现盈亏
		pos.MarkPrice = currentPrice
		pos.UnrealizedPnL = t.calculateUnrealizedPnL(pos)
		if err := t.store.PaperAccount().UpdatePosition(pos); err != nil {
			logger.Warnf("⚠️ [Paper] Failed to update position mark price: %v", err)
		}

		posAmt := pos.Quantity
		if pos.Side == "SHORT" {
			posAmt = -pos.Quantity
		}

		result = append(result, map[string]interface{}{
			"symbol":           pos.Symbol,
			"side":             strings.ToLower(pos.Side),
			"positionAmt":      posAmt,
			"entryPrice":       pos.EntryPrice,
			"markPrice":        pos.MarkPrice,
			"unRealizedProfit": pos.UnrealizedPnL,
			"leverage":         float64(pos.Leverage),
			"liquidationPrice": pos.LiquidationPrice,
		})
	}
	return result, nil
}

// OpenLong 开多仓
func (t *PaperTrader) OpenLong(symbol string, quantity float64, leverage int) (map[string]interface{}, error) {
	return t.openPosition(symbol, "LONG", quantity, leverage)
}

// OpenShort 开空仓
func (t *PaperTrader) OpenShort(symbol string, quantity float64, leverage int) (map[string]interface{}, error) {
	return t.openPosition(symbol, "SHORT", quantity, leverage)
}

func (t *PaperTrader) openPosition(symbol, side string, quantity float64, leverage int) (map[string]interface{}, error) {
	// 1. 获取市场价格（不需要在事务中）
	md, err := market.Get(symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get market price: %w", err)
	}

	// 2. 应用滑点
	execPrice := t.applySlippage(md.CurrentPrice, side, true)

	// 3. 计算保证金和手续费
	notional := execPrice * quantity
	margin := notional / float64(leverage)
	fee := notional * t.feeRate

	// 4. 预先检查账户余额（在事务外）
	account, err := t.store.PaperAccount().GetByID(t.accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("paper account not found: %s", t.accountID)
	}

	if margin+fee > account.CurrentBalance {
		return nil, fmt.Errorf("insufficient balance: need %.2f, available %.2f", margin+fee, account.CurrentBalance)
	}

	// 5. 检查现有持仓（在事务外）
	existingPos, _ := t.store.PaperAccount().GetPosition(t.accountID, symbol, side)

	// 6. 准备订单对象
	orderSide := "BUY"
	if side == "SHORT" {
		orderSide = "SELL"
	}
	filledAt := time.Now()
	order := &store.PaperOrder{
		AccountID:    t.accountID,
		Symbol:       symbol,
		Side:         orderSide,
		PositionSide: side,
		OrderType:    "MARKET",
		Quantity:     quantity,
		Price:        md.CurrentPrice,
		AvgPrice:     execPrice,
		Fee:          fee,
		Status:       "FILLED",
		FilledAt:     &filledAt,
	}

	// 7. 使用事务执行所有数据库更新
	err = t.store.Transaction(func(tx *sql.Tx) error {
		// 7.1 更新账户余额
		account.CurrentBalance -= (margin + fee)
		account.TotalMarginUsed += margin
		if err := t.store.PaperAccount().UpdateTx(tx, account); err != nil {
			return fmt.Errorf("failed to update account: %w", err)
		}

		// 7.2 创建或更新持仓
		if existingPos != nil {
			// 加仓
			newQty := existingPos.Quantity + quantity
			existingPos.EntryPrice = (existingPos.EntryPrice*existingPos.Quantity + execPrice*quantity) / newQty
			existingPos.Quantity = newQty
			existingPos.MarginUsed += margin
			existingPos.Leverage = leverage
			existingPos.LiquidationPrice = t.calculateLiquidationPrice(existingPos.EntryPrice, leverage, side)
			if err := t.store.PaperAccount().UpdatePositionTx(tx, existingPos); err != nil {
				return fmt.Errorf("failed to update position: %w", err)
			}
		} else {
			// 新建持仓
			pos := &store.PaperPosition{
				AccountID:        t.accountID,
				Symbol:           symbol,
				Side:             side,
				Quantity:         quantity,
				EntryPrice:       execPrice,
				MarkPrice:        execPrice,
				Leverage:         leverage,
				MarginUsed:       margin,
				LiquidationPrice: t.calculateLiquidationPrice(execPrice, leverage, side),
				OpenTime:         time.Now(),
			}
			if err := t.store.PaperAccount().CreatePositionTx(tx, pos); err != nil {
				return fmt.Errorf("failed to create position: %w", err)
			}
		}

		// 7.3 创建订单记录
		if err := t.store.PaperAccount().CreateOrderTx(tx, order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Infof("📝 [Paper] Opened %s %s: qty=%.4f, price=%.4f, fee=%.4f", side, symbol, quantity, execPrice, fee)

	return map[string]interface{}{
		"orderId": order.ID,
		"status":  "FILLED",
	}, nil
}

// CloseLong 平多仓
func (t *PaperTrader) CloseLong(symbol string, quantity float64) (map[string]interface{}, error) {
	return t.closePosition(symbol, "LONG", quantity)
}

// CloseShort 平空仓
func (t *PaperTrader) CloseShort(symbol string, quantity float64) (map[string]interface{}, error) {
	return t.closePosition(symbol, "SHORT", quantity)
}

func (t *PaperTrader) closePosition(symbol, side string, quantity float64) (map[string]interface{}, error) {
	// 1. 获取持仓（在事务外）
	pos, err := t.store.PaperAccount().GetPosition(t.accountID, symbol, side)
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}
	if pos == nil {
		return nil, fmt.Errorf("no %s position for %s", side, symbol)
	}

	if quantity <= 0 || quantity > pos.Quantity {
		quantity = pos.Quantity // 平仓全部
	}

	// 2. 获取市场价格（在事务外）
	md, err := market.Get(symbol)
	if err != nil {
		return nil, err
	}

	execPrice := t.applySlippage(md.CurrentPrice, side, false)
	notional := execPrice * quantity
	fee := notional * t.feeRate

	// 3. 计算已实现盈亏
	var realizedPnL float64
	if side == "LONG" {
		realizedPnL = (execPrice - pos.EntryPrice) * quantity
	} else {
		realizedPnL = (pos.EntryPrice - execPrice) * quantity
	}
	realizedPnL -= fee

	// 4. 释放保证金
	marginReleased := pos.MarginUsed * (quantity / pos.Quantity)

	// 5. 获取账户（在事务外）
	account, err := t.store.PaperAccount().GetByID(t.accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("paper account not found: %s", t.accountID)
	}

	// 6. 准备订单对象
	orderSide := "SELL"
	if side == "SHORT" {
		orderSide = "BUY"
	}
	filledAt := time.Now()
	order := &store.PaperOrder{
		AccountID:    t.accountID,
		Symbol:       symbol,
		Side:         orderSide,
		PositionSide: side,
		OrderType:    "MARKET",
		Quantity:     quantity,
		Price:        md.CurrentPrice,
		AvgPrice:     execPrice,
		Fee:          fee,
		Status:       "FILLED",
		RealizedPnL:  realizedPnL,
		FilledAt:     &filledAt,
	}

	// 7. 判断是全部平仓还是部分平仓
	isFullClose := quantity >= pos.Quantity

	// 8. 使用事务执行所有数据库更新
	err = t.store.Transaction(func(tx *sql.Tx) error {
		// 8.1 更新账户余额
		account.CurrentBalance += marginReleased + realizedPnL
		account.TotalMarginUsed -= marginReleased
		account.TotalPnL += realizedPnL
		if err := t.store.PaperAccount().UpdateTx(tx, account); err != nil {
			return fmt.Errorf("failed to update account: %w", err)
		}

		// 8.2 更新或删除持仓
		if isFullClose {
			if err := t.store.PaperAccount().DeletePositionTx(tx, pos.ID); err != nil {
				return fmt.Errorf("failed to delete position: %w", err)
			}
		} else {
			pos.Quantity -= quantity
			pos.MarginUsed -= marginReleased
			if err := t.store.PaperAccount().UpdatePositionTx(tx, pos); err != nil {
				return fmt.Errorf("failed to update position: %w", err)
			}
		}

		// 8.3 创建订单记录
		if err := t.store.PaperAccount().CreateOrderTx(tx, order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Infof("📝 [Paper] Closed %s %s: qty=%.4f, price=%.4f, pnl=%.4f", side, symbol, quantity, execPrice, realizedPnL)

	return map[string]interface{}{
		"orderId": order.ID,
		"status":  "FILLED",
	}, nil
}

// SetLeverage 设置杠杆
// 注意：Paper Trading 的杠杆通过 OpenLong/OpenShort 参数直接传递
// 此方法仅记录日志，与真实交易所 API 保持接口一致
func (t *PaperTrader) SetLeverage(symbol string, leverage int) error {
	logger.Infof("📝 [Paper] Set leverage for %s: %dx", symbol, leverage)
	return nil
}

// SetMarginMode 设置保证金模式
func (t *PaperTrader) SetMarginMode(symbol string, isCrossMargin bool) error {
	t.isCrossMargin = isCrossMargin
	return nil
}

// GetMarketPrice 获取市场价格
func (t *PaperTrader) GetMarketPrice(symbol string) (float64, error) {
	md, err := market.Get(symbol)
	if err != nil {
		return 0, err
	}
	return md.CurrentPrice, nil
}

// SetStopLoss 设置止损
func (t *PaperTrader) SetStopLoss(symbol string, positionSide string, quantity, stopPrice float64) error {
	side := strings.ToUpper(positionSide)
	pos, err := t.store.PaperAccount().GetPosition(t.accountID, symbol, side)
	if err != nil {
		return fmt.Errorf("failed to get position: %w", err)
	}
	if pos == nil {
		return fmt.Errorf("no %s position for %s", side, symbol)
	}
	pos.StopLoss = stopPrice
	if err := t.store.PaperAccount().UpdatePosition(pos); err != nil {
		return fmt.Errorf("failed to update position: %w", err)
	}
	return nil
}

// SetTakeProfit 设置止盈
func (t *PaperTrader) SetTakeProfit(symbol string, positionSide string, quantity, takeProfitPrice float64) error {
	side := strings.ToUpper(positionSide)
	pos, err := t.store.PaperAccount().GetPosition(t.accountID, symbol, side)
	if err != nil {
		return fmt.Errorf("failed to get position: %w", err)
	}
	if pos == nil {
		return fmt.Errorf("no %s position for %s", side, symbol)
	}
	pos.TakeProfit = takeProfitPrice
	if err := t.store.PaperAccount().UpdatePosition(pos); err != nil {
		return fmt.Errorf("failed to update position: %w", err)
	}
	return nil
}

// CancelStopLossOrders Paper trading 没有挂单
func (t *PaperTrader) CancelStopLossOrders(symbol string) error {
	return nil // Paper trading doesn't have pending orders
}

// CancelTakeProfitOrders Paper trading 没有挂单
func (t *PaperTrader) CancelTakeProfitOrders(symbol string) error {
	return nil
}

// CancelAllOrders Paper trading 没有挂单
func (t *PaperTrader) CancelAllOrders(symbol string) error {
	return nil
}

// CancelStopOrders Paper trading 没有挂单
func (t *PaperTrader) CancelStopOrders(symbol string) error {
	return nil
}

// FormatQuantity 格式化数量精度
func (t *PaperTrader) FormatQuantity(symbol string, quantity float64) (string, error) {
	// 使用通用精度格式化
	return fmt.Sprintf("%.4f", quantity), nil
}

// GetOrderStatus 获取订单状态
func (t *PaperTrader) GetOrderStatus(symbol string, orderID string) (map[string]interface{}, error) {
	// Paper orders 都是立即成交
	return map[string]interface{}{
		"orderId": orderID,
		"status":  "FILLED",
	}, nil
}

// GetClosedPnL 获取已平仓盈亏记录
// Paper Trading 限制说明:
// - EntryPrice: 平仓订单不包含入场价，返回 0
// - EntryTime: 平仓订单不包含原始入场时间，返回零值
// - Leverage: 平仓订单不包含杠杆信息，返回 0
// 这些信息在真实交易所 API 中也通常不包含在平仓记录里
func (t *PaperTrader) GetClosedPnL(startTime time.Time, limit int) ([]ClosedPnLRecord, error) {
	orders, err := t.store.PaperAccount().ListOrders(t.accountID, limit)
	if err != nil {
		return nil, err
	}

	var records []ClosedPnLRecord
	for _, order := range orders {
		if order.RealizedPnL != 0 && order.FilledAt != nil && order.FilledAt.After(startTime) {
			records = append(records, ClosedPnLRecord{
				Symbol:      order.Symbol,
				Side:        strings.ToLower(order.PositionSide),
				EntryPrice:  0,           // Paper Trading: 平仓记录不包含入场价（设计限制）
				ExitPrice:   order.AvgPrice,
				Quantity:    order.Quantity,
				RealizedPnL: order.RealizedPnL,
				Fee:         order.Fee,
				Leverage:    0,           // Paper Trading: 平仓记录不包含杠杆（设计限制）
				EntryTime:   time.Time{}, // Paper Trading: 平仓记录不包含入场时间（设计限制）
				ExitTime:    *order.FilledAt,
				OrderID:     fmt.Sprintf("%d", order.ID),
				CloseType:   determineCloseType(order.OrderType),
				ExchangeID:  fmt.Sprintf("paper_%d", order.ID),
			})
		}
	}
	return records, nil
}

// determineCloseType 根据订单类型判断平仓类型
func determineCloseType(orderType string) string {
	switch orderType {
	case "STOP_LOSS":
		return "stop_loss"
	case "TAKE_PROFIT":
		return "take_profit"
	case "LIQUIDATION":
		return "liquidation"
	default:
		return "manual"
	}
}

// =============================================================================
// 辅助方法
// =============================================================================

// applySlippage 应用滑点
// 滑点模拟真实市场中的成交价偏差，总是对交易者不利：
// - LONG 开仓（买入）：成交价更高 → 入场成本增加
// - LONG 平仓（卖出）：成交价更低 → 卖出收益减少
// - SHORT 开仓（卖出建仓）：成交价更低 → 建仓价格不利（空头希望高价卖出）
// - SHORT 平仓（买入平仓）：成交价更高 → 平仓成本增加
func (t *PaperTrader) applySlippage(price float64, side string, isOpen bool) float64 {
	adjust := 1.0
	if side == "LONG" {
		if isOpen {
			adjust += t.slippageRate // 买入时价格略高（对用户不利）
		} else {
			adjust -= t.slippageRate // 卖出时价格略低（对用户不利）
		}
	} else { // SHORT
		if isOpen {
			adjust -= t.slippageRate // 做空开仓时价格略低（对用户不利，因为做空是卖出建仓，价格低意味着卖便宜了）
		} else {
			adjust += t.slippageRate // 做空平仓时价格略高（对用户不利，因为平仓是买回，价格高意味着买贵了）
		}
	}
	return price * adjust
}

// calculateLiquidationPrice 计算强平价格
func (t *PaperTrader) calculateLiquidationPrice(entry float64, leverage int, side string) float64 {
	// 防止除零：杠杆必须大于 0
	if leverage <= 0 {
		logger.Warnf("⚠️ [Paper] Invalid leverage %d for liquidation price calculation, using 1", leverage)
		leverage = 1
	}
	lev := float64(leverage)
	if side == "LONG" {
		return entry * (1.0 - 0.9/lev) // 保留 10% 安全边际
	}
	return entry * (1.0 + 0.9/lev)
}

// calculateUnrealizedPnL 计算未实现盈亏
func (t *PaperTrader) calculateUnrealizedPnL(pos *store.PaperPosition) float64 {
	if pos.Side == "LONG" {
		return (pos.MarkPrice - pos.EntryPrice) * pos.Quantity
	}
	return (pos.EntryPrice - pos.MarkPrice) * pos.Quantity
}

// executeAutoClose 执行自动平仓（止损/止盈/强平通用）
// 参数:
//   - pos: 要平仓的持仓
//   - execPrice: 执行价格（触发时的市场价）
//   - orderType: "STOP_LOSS", "TAKE_PROFIT", "LIQUIDATION"
//
// 返回:
//   - error: 平仓失败时返回错误，成功返回 nil
func (t *PaperTrader) executeAutoClose(pos *store.PaperPosition, execPrice float64, orderType string) error {
	quantity := pos.Quantity
	notional := execPrice * quantity
	fee := notional * t.feeRate

	// 1. 计算已实现盈亏
	var realizedPnL float64
	if pos.Side == "LONG" {
		realizedPnL = (execPrice - pos.EntryPrice) * quantity
	} else {
		realizedPnL = (pos.EntryPrice - execPrice) * quantity
	}
	realizedPnL -= fee

	// 2. 释放保证金
	marginReleased := pos.MarginUsed

	// 3. 获取账户（在事务外）
	account, err := t.store.PaperAccount().GetByID(t.accountID)
	if err != nil || account == nil {
		logger.Errorf("❌ [Paper] Failed to get account for auto-close: %v", err)
		return fmt.Errorf("failed to get account for auto-close: %w", err)
	}

	// 4. 准备订单对象
	orderSide := "SELL"
	if pos.Side == "SHORT" {
		orderSide = "BUY"
	}
	filledAt := time.Now()
	order := &store.PaperOrder{
		AccountID:    t.accountID,
		Symbol:       pos.Symbol,
		Side:         orderSide,
		PositionSide: pos.Side,
		OrderType:    orderType, // "STOP_LOSS", "TAKE_PROFIT", or "LIQUIDATION"
		Quantity:     quantity,
		Price:        execPrice,
		AvgPrice:     execPrice,
		Fee:          fee,
		Status:       "FILLED",
		RealizedPnL:  realizedPnL,
		FilledAt:     &filledAt,
	}

	// 5. 使用事务执行所有数据库更新
	err = t.store.Transaction(func(tx *sql.Tx) error {
		// 5.1 更新账户余额
		account.CurrentBalance += marginReleased + realizedPnL
		account.TotalMarginUsed -= marginReleased
		account.TotalPnL += realizedPnL
		if err := t.store.PaperAccount().UpdateTx(tx, account); err != nil {
			return fmt.Errorf("failed to update account for auto-close: %w", err)
		}

		// 5.2 删除仓位
		if err := t.store.PaperAccount().DeletePositionTx(tx, pos.ID); err != nil {
			return fmt.Errorf("failed to delete position for auto-close: %w", err)
		}

		// 5.3 创建订单记录
		if err := t.store.PaperAccount().CreateOrderTx(tx, order); err != nil {
			return fmt.Errorf("failed to create order for auto-close: %w", err)
		}

		return nil
	})

	if err != nil {
		logger.Errorf("❌ [Paper] Auto-close transaction failed: %v", err)
		return err
	}

	logger.Infof("📝 [Paper] Auto-closed %s %s (%s): qty=%.4f, price=%.4f, pnl=%.4f",
		pos.Side, pos.Symbol, orderType, quantity, execPrice, realizedPnL)
	return nil
}

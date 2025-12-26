import { Shield, AlertTriangle } from 'lucide-react'
import type { RiskControlConfig } from '../../types'

interface RiskControlEditorProps {
  config: RiskControlConfig
  onChange: (config: RiskControlConfig) => void
  disabled?: boolean
  language: string
}

export function RiskControlEditor({
  config,
  onChange,
  disabled,
  language,
}: RiskControlEditorProps) {
  const t = (key: string) => {
    const translations: Record<string, Record<string, string>> = {
      positionLimits: { zh: '仓位限制', en: 'Position Limits' },
      maxPositions: { zh: '最大持仓数量', en: 'Max Positions' },
      maxPositionsDesc: { zh: '同时持有的最大币种数量', en: 'Maximum coins held simultaneously' },
      // Trading leverage (exchange leverage)
      tradingLeverage: { zh: '交易杠杆（交易所杠杆）', en: 'Trading Leverage (Exchange)' },
      btcEthLeverage: { zh: 'BTC/ETH 交易杠杆', en: 'BTC/ETH Trading Leverage' },
      btcEthLeverageDesc: { zh: '交易所开仓使用的杠杆倍数', en: 'Exchange leverage for opening positions' },
      altcoinLeverage: { zh: '山寨币交易杠杆', en: 'Altcoin Trading Leverage' },
      altcoinLeverageDesc: { zh: '交易所开仓使用的杠杆倍数', en: 'Exchange leverage for opening positions' },
      // Position value ratio (risk control) - CODE ENFORCED
      positionValueRatio: { zh: '仓位价值比例（代码强制）', en: 'Position Value Ratio (CODE ENFORCED)' },
      positionValueRatioDesc: { zh: '单仓位名义价值 / 账户净值，由代码强制执行', en: 'Position notional value / equity, enforced by code' },
      btcEthPositionValueRatio: { zh: 'BTC/ETH 仓位价值比例', en: 'BTC/ETH Position Value Ratio' },
      btcEthPositionValueRatioDesc: { zh: '单仓最大名义价值 = 净值 × 此值（代码强制）', en: 'Max position value = equity × this ratio (CODE ENFORCED)' },
      altcoinPositionValueRatio: { zh: '山寨币仓位价值比例', en: 'Altcoin Position Value Ratio' },
      altcoinPositionValueRatioDesc: { zh: '单仓最大名义价值 = 净值 × 此值（代码强制）', en: 'Max position value = equity × this ratio (CODE ENFORCED)' },
      riskParameters: { zh: '风险参数', en: 'Risk Parameters' },
      minRiskReward: { zh: '最小风险回报比', en: 'Min Risk/Reward Ratio' },
      minRiskRewardDesc: { zh: '开仓要求的最低盈亏比', en: 'Minimum profit ratio for opening' },
      maxMarginUsage: { zh: '最大保证金使用率（代码强制）', en: 'Max Margin Usage (CODE ENFORCED)' },
      maxMarginUsageDesc: { zh: '保证金使用率上限，由代码强制执行', en: 'Maximum margin utilization, enforced by code' },
      frequencyCostControls: { zh: '频率与成本约束（代码强制）', en: 'Frequency & Cost Controls (CODE ENFORCED)' },
      minHoldMinutes: { zh: '最短持仓时间（分钟）', en: 'Min Hold Time (minutes)' },
      minHoldMinutesDesc: {
        zh: 'AI 主动平仓前必须至少持仓该分钟数（0=关闭）',
        en: 'Minimum holding time before AI can close (0=disabled)',
      },
      cooldownAfterCloseMinutes: { zh: '平仓冷却（分钟）', en: 'Cooldown After Close (minutes)' },
      cooldownAfterCloseMinutesDesc: {
        zh: '平仓后该分钟数内禁止同币种再次开仓（0=关闭）',
        en: 'Block reopening the same symbol for N minutes (0=disabled)',
      },
      minTPCostMultiplier: { zh: '止盈距离/成本倍数', en: 'TP Distance / Cost Multiplier' },
      minTPCostMultiplierDesc: {
        zh: 'PAPER 模式下要求 TP 距离 ≥ round-trip 成本 × 倍数（0=关闭）',
        en: 'In PAPER mode, require TP distance ≥ round-trip cost × multiplier (0=disabled)',
      },
      minStopLossDistancePct: { zh: '最小止损距离（%）', en: 'Min Stop Loss Distance (%)' },
      minStopLossDistancePctDesc: {
        zh: '开仓要求止损距离 ≥ 该百分比（0=关闭，例如 0.25=0.25%）',
        en: 'Require SL distance ≥ this percent (0=disabled, e.g. 0.25=0.25%)',
      },
      entryRequirements: { zh: '开仓要求', en: 'Entry Requirements' },
      minPositionSize: { zh: '最小开仓金额', en: 'Min Position Size' },
      minPositionSizeDesc: { zh: 'USDT 最小名义价值', en: 'Minimum notional value in USDT' },
      minConfidence: { zh: '最小信心度', en: 'Min Confidence' },
      minConfidenceDesc: { zh: 'AI 开仓信心度阈值', en: 'AI confidence threshold for entry' },
    }
    return translations[key]?.[language] || key
  }

  const updateField = <K extends keyof RiskControlConfig>(
    key: K,
    value: RiskControlConfig[K]
  ) => {
    if (!disabled) {
      onChange({ ...config, [key]: value })
    }
  }

  return (
    <div className="space-y-6">
      {/* Position Limits */}
      <div>
        <div className="flex items-center gap-2 mb-4">
          <Shield className="w-5 h-5" style={{ color: '#F0B90B' }} />
          <h3 className="font-medium" style={{ color: '#EAECEF' }}>
            {t('positionLimits')}
          </h3>
        </div>

        <div className="grid grid-cols-1 gap-4 mb-4">
          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('maxPositions')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('maxPositionsDesc')}
            </p>
            <input
              type="number"
              value={config.max_positions ?? 3}
              onChange={(e) =>
                updateField('max_positions', parseInt(e.target.value) || 3)
              }
              disabled={disabled}
              min={1}
              max={10}
              className="w-32 px-3 py-2 rounded"
              style={{
                background: '#1E2329',
                border: '1px solid #2B3139',
                color: '#EAECEF',
              }}
            />
          </div>
        </div>

        {/* Trading Leverage (Exchange) */}
        <div className="mb-2">
          <p className="text-xs font-medium mb-2" style={{ color: '#F0B90B' }}>
            {t('tradingLeverage')}
          </p>
        </div>
        <div className="grid grid-cols-2 gap-4 mb-4">
          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('btcEthLeverage')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('btcEthLeverageDesc')}
            </p>
            <div className="flex items-center gap-2">
              <input
                type="range"
                value={config.btc_eth_max_leverage ?? 5}
                onChange={(e) =>
                  updateField('btc_eth_max_leverage', parseInt(e.target.value))
                }
                disabled={disabled}
                min={1}
                max={20}
                className="flex-1 accent-yellow-500"
              />
              <span
                className="w-12 text-center font-mono"
                style={{ color: '#F0B90B' }}
              >
                {config.btc_eth_max_leverage ?? 5}x
              </span>
            </div>
          </div>

          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('altcoinLeverage')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('altcoinLeverageDesc')}
            </p>
            <div className="flex items-center gap-2">
              <input
                type="range"
                value={config.altcoin_max_leverage ?? 5}
                onChange={(e) =>
                  updateField('altcoin_max_leverage', parseInt(e.target.value))
                }
                disabled={disabled}
                min={1}
                max={20}
                className="flex-1 accent-yellow-500"
              />
              <span
                className="w-12 text-center font-mono"
                style={{ color: '#F0B90B' }}
              >
                {config.altcoin_max_leverage ?? 5}x
              </span>
            </div>
          </div>
        </div>

        {/* Position Value Ratio (Risk Control - CODE ENFORCED) */}
        <div className="mb-2">
          <p className="text-xs font-medium" style={{ color: '#0ECB81' }}>
            {t('positionValueRatio')}
          </p>
          <p className="text-xs mt-1" style={{ color: '#848E9C' }}>
            {t('positionValueRatioDesc')}
          </p>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #0ECB81' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('btcEthPositionValueRatio')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('btcEthPositionValueRatioDesc')}
            </p>
            <div className="flex items-center gap-2">
              <input
                type="range"
                value={config.btc_eth_max_position_value_ratio ?? 5}
                onChange={(e) =>
                  updateField('btc_eth_max_position_value_ratio', parseFloat(e.target.value))
                }
                disabled={disabled}
                min={0.5}
                max={10}
                step={0.5}
                className="flex-1 accent-green-500"
              />
              <span
                className="w-12 text-center font-mono"
                style={{ color: '#0ECB81' }}
              >
                {config.btc_eth_max_position_value_ratio ?? 5}x
              </span>
            </div>
          </div>

          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #0ECB81' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('altcoinPositionValueRatio')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('altcoinPositionValueRatioDesc')}
            </p>
            <div className="flex items-center gap-2">
              <input
                type="range"
                value={config.altcoin_max_position_value_ratio ?? 1}
                onChange={(e) =>
                  updateField('altcoin_max_position_value_ratio', parseFloat(e.target.value))
                }
                disabled={disabled}
                min={0.5}
                max={10}
                step={0.5}
                className="flex-1 accent-green-500"
              />
              <span
                className="w-12 text-center font-mono"
                style={{ color: '#0ECB81' }}
              >
                {config.altcoin_max_position_value_ratio ?? 1}x
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Risk Parameters */}
      <div>
        <div className="flex items-center gap-2 mb-4">
          <AlertTriangle className="w-5 h-5" style={{ color: '#F6465D' }} />
          <h3 className="font-medium" style={{ color: '#EAECEF' }}>
            {t('riskParameters')}
          </h3>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('minRiskReward')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('minRiskRewardDesc')}
            </p>
            <div className="flex items-center">
              <span style={{ color: '#848E9C' }}>1:</span>
              <input
                type="number"
                value={config.min_risk_reward_ratio ?? 3}
                onChange={(e) =>
                  updateField('min_risk_reward_ratio', parseFloat(e.target.value) || 3)
                }
                disabled={disabled}
                min={1}
                max={10}
                step={0.5}
                className="w-20 px-3 py-2 rounded ml-2"
                style={{
                  background: '#1E2329',
                  border: '1px solid #2B3139',
                  color: '#EAECEF',
                }}
              />
            </div>
          </div>

          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #0ECB81' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('maxMarginUsage')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('maxMarginUsageDesc')}
            </p>
            <div className="flex items-center gap-2">
              <input
                type="range"
                value={(config.max_margin_usage ?? 0.9) * 100}
                onChange={(e) =>
                  updateField('max_margin_usage', parseInt(e.target.value) / 100)
                }
                disabled={disabled}
                min={10}
                max={100}
                className="flex-1 accent-green-500"
              />
              <span className="w-12 text-center font-mono" style={{ color: '#0ECB81' }}>
                {Math.round((config.max_margin_usage ?? 0.9) * 100)}%
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Frequency & Cost Controls */}
      <div>
        <div className="flex items-center gap-2 mb-4">
          <Shield className="w-5 h-5" style={{ color: '#F0B90B' }} />
          <h3 className="font-medium" style={{ color: '#EAECEF' }}>
            {t('frequencyCostControls')}
          </h3>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('minHoldMinutes')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('minHoldMinutesDesc')}
            </p>
            <div className="flex items-center">
              <input
                type="number"
                value={config.min_hold_minutes ?? 0}
                onChange={(e) =>
                  updateField('min_hold_minutes', parseInt(e.target.value) || 0)
                }
                disabled={disabled}
                min={0}
                max={1440}
                className="w-24 px-3 py-2 rounded"
                style={{
                  background: '#1E2329',
                  border: '1px solid #2B3139',
                  color: '#EAECEF',
                }}
              />
              <span className="ml-2" style={{ color: '#848E9C' }}>
                min
              </span>
            </div>
          </div>

          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('cooldownAfterCloseMinutes')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('cooldownAfterCloseMinutesDesc')}
            </p>
            <div className="flex items-center">
              <input
                type="number"
                value={config.cooldown_after_close_minutes ?? 0}
                onChange={(e) =>
                  updateField(
                    'cooldown_after_close_minutes',
                    parseInt(e.target.value) || 0
                  )
                }
                disabled={disabled}
                min={0}
                max={1440}
                className="w-24 px-3 py-2 rounded"
                style={{
                  background: '#1E2329',
                  border: '1px solid #2B3139',
                  color: '#EAECEF',
                }}
              />
              <span className="ml-2" style={{ color: '#848E9C' }}>
                min
              </span>
            </div>
          </div>

          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('minTPCostMultiplier')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('minTPCostMultiplierDesc')}
            </p>
            <div className="flex items-center">
              <input
                type="number"
                value={config.min_tp_cost_multiplier ?? 0}
                onChange={(e) =>
                  updateField('min_tp_cost_multiplier', parseFloat(e.target.value) || 0)
                }
                disabled={disabled}
                min={0}
                max={20}
                step={0.5}
                className="w-24 px-3 py-2 rounded"
                style={{
                  background: '#1E2329',
                  border: '1px solid #2B3139',
                  color: '#EAECEF',
                }}
              />
              <span className="ml-2" style={{ color: '#848E9C' }}>
                x
              </span>
            </div>
          </div>

          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('minStopLossDistancePct')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('minStopLossDistancePctDesc')}
            </p>
            <div className="flex items-center">
              <input
                type="number"
                value={config.min_stop_loss_distance_pct ?? 0}
                onChange={(e) =>
                  updateField(
                    'min_stop_loss_distance_pct',
                    parseFloat(e.target.value) || 0
                  )
                }
                disabled={disabled}
                min={0}
                max={10}
                step={0.05}
                className="w-24 px-3 py-2 rounded"
                style={{
                  background: '#1E2329',
                  border: '1px solid #2B3139',
                  color: '#EAECEF',
                }}
              />
              <span className="ml-2" style={{ color: '#848E9C' }}>
                %
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Entry Requirements */}
      <div>
        <div className="flex items-center gap-2 mb-4">
          <Shield className="w-5 h-5" style={{ color: '#0ECB81' }} />
          <h3 className="font-medium" style={{ color: '#EAECEF' }}>
            {t('entryRequirements')}
          </h3>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('minPositionSize')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('minPositionSizeDesc')}
            </p>
            <div className="flex items-center">
              <input
                type="number"
                value={config.min_position_size ?? 12}
                onChange={(e) =>
                  updateField('min_position_size', parseFloat(e.target.value) || 12)
                }
                disabled={disabled}
                min={10}
                max={1000}
                className="w-24 px-3 py-2 rounded"
                style={{
                  background: '#1E2329',
                  border: '1px solid #2B3139',
                  color: '#EAECEF',
                }}
              />
              <span className="ml-2" style={{ color: '#848E9C' }}>
                USDT
              </span>
            </div>
          </div>

          <div
            className="p-4 rounded-lg"
            style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
          >
            <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
              {t('minConfidence')}
            </label>
            <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
              {t('minConfidenceDesc')}
            </p>
            <div className="flex items-center gap-2">
              <input
                type="range"
                value={config.min_confidence ?? 75}
                onChange={(e) =>
                  updateField('min_confidence', parseInt(e.target.value))
                }
                disabled={disabled}
                min={50}
                max={100}
                className="flex-1 accent-green-500"
              />
              <span className="w-12 text-center font-mono" style={{ color: '#0ECB81' }}>
                {config.min_confidence ?? 75}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

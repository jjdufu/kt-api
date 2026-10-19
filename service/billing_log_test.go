package service

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/stretchr/testify/assert"
)

func TestSettleBillingLogLine(t *testing.T) {
	general := operation_setting.GetGeneralSetting()
	previousDisplay := general.QuotaDisplayType
	t.Cleanup(func() {
		general.QuotaDisplayType = previousDisplay
		common.SetLogLanguage("")
	})
	general.QuotaDisplayType = operation_setting.QuotaDisplayTypeUSD

	task := &relaycommon.RelayInfo{TaskRelayInfo: &relaycommon.TaskRelayInfo{}, OriginModelName: "doubao-seedance-2-0-mini-260615"}
	perRequestTask := &relaycommon.RelayInfo{TaskRelayInfo: &relaycommon.TaskRelayInfo{}}
	perRequestTask.PriceData.UsePrice = true
	chat := &relaycommon.RelayInfo{}

	t.Run("a usage-billed task pre-deduct is logged as pending settlement without a currency", func(t *testing.T) {
		line := settleBillingLogLine(task, 444789, 444789)
		assert.Equal(t, "task pre-deducted 444789 quota points, pending settlement: the final amount is settled from actual usage when the task completes", line)
		assert.NotContains(t, line, "per-request")
		assert.NotContains(t, line, "$")
		assert.NotContains(t, line, "＄")
	})
	t.Run("zh-CN", func(t *testing.T) {
		common.SetLogLanguage("zh-CN")
		defer common.SetLogLanguage("")
		assert.Equal(t, "任务已预扣 444789 点额度，待结算：任务完成后按实际用量结算最终金额", settleBillingLogLine(task, 444789, 444789))
		assert.Equal(t, "预扣费与实际消耗一致，无需调整：100 点额度（按次计费）", settleBillingLogLine(perRequestTask, 100, 100))
		assert.Equal(t, "预扣费与实际消耗一致，无需调整：100 点额度", settleBillingLogLine(chat, 100, 100))
	})
	t.Run("per-request and non-task wording", func(t *testing.T) {
		assert.Equal(t, "pre-consumed quota matches the actual quota, no adjustment: 100 quota points (per-request billing)", settleBillingLogLine(perRequestTask, 100, 100))
		assert.Equal(t, "pre-consumed quota matches the actual quota, no adjustment: 100 quota points", settleBillingLogLine(chat, 100, 100))
	})
	t.Run("differences are charged or refunded", func(t *testing.T) {
		assert.Equal(t, "charging the difference after pre-consume: 25046 quota points (actual 469835 quota points, pre-consumed 444789 quota points)", settleBillingLogLine(task, 469835, 444789))
		assert.Equal(t, "refunding the difference after pre-consume: 10 quota points (actual 90 quota points, pre-consumed 100 quota points)", settleBillingLogLine(chat, 90, 100))
	})
	t.Run("a configured CNY display keeps its symbol", func(t *testing.T) {
		general.QuotaDisplayType = operation_setting.QuotaDisplayTypeCNY
		defer func() { general.QuotaDisplayType = operation_setting.QuotaDisplayTypeUSD }()
		assert.Contains(t, settleBillingLogLine(task, 500000, 500000), "¥")
	})
}

package service

import (
	"fmt"
	"net/http"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/logger"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/types"
	"github.com/gin-gonic/gin"
)

const (
	BillingSourceWallet       = "wallet"
	BillingSourceSubscription = "subscription"
)

// PreConsumeBilling 根据用户计费偏好创建 BillingSession 并执行预扣费。
// 会话存储在 relayInfo.Billing 上，供后续 Settle / Refund 使用。
func PreConsumeBilling(c *gin.Context, preConsumedQuota int, relayInfo *relaycommon.RelayInfo) *types.NewAPIError {
	if relayInfo != nil && relayInfo.QuotaClamp != nil {
		return types.NewErrorWithStatusCode(
			relayInfo.QuotaClamp,
			types.ErrorCodeModelPriceError,
			http.StatusBadRequest,
			types.ErrOptionWithSkipRetry(),
		)
	}
	if preConsumedQuota < 0 {
		return types.NewErrorWithStatusCode(
			fmt.Errorf("pre-consume quota cannot be negative: %d", preConsumedQuota),
			types.ErrorCodeModelPriceError,
			http.StatusBadRequest,
			types.ErrOptionWithSkipRetry(),
		)
	}
	session, apiErr := NewBillingSession(c, relayInfo, preConsumedQuota)
	if apiErr != nil {
		return apiErr
	}
	relayInfo.Billing = session
	return nil
}

// ---------------------------------------------------------------------------
// SettleBilling — 后结算辅助函数
// ---------------------------------------------------------------------------

// taskSettlesOnCompletion reports whether a task submission's quota is only
// a pre-deduction: unless the model is billed per request, the final amount is
// settled from actual usage when the task completes (usage expression, adaptor
// adjustment or token recount; see settleTaskBillingOnComplete).
func taskSettlesOnCompletion(relayInfo *relaycommon.RelayInfo) bool {
	if relayInfo == nil || relayInfo.TaskRelayInfo == nil {
		return false
	}
	return !billedPerRequest(relayInfo)
}

func billedPerRequest(relayInfo *relaycommon.RelayInfo) bool {
	return relayInfo.PriceData.UsePrice || common.StringsContains(constant.TaskPricePatches, relayInfo.OriginModelName)
}

// settleBillingLogLine describes a settlement in the log language, with
// quotas formatted by the site's display setting (logger.LogQuota).
func settleBillingLogLine(relayInfo *relaycommon.RelayInfo, actualQuota, preConsumed int) string {
	delta := actualQuota - preConsumed
	switch {
	case delta > 0:
		return common.LogText("charging the difference after pre-consume: %s (actual %s, pre-consumed %s)",
			logger.LogQuota(delta), logger.LogQuota(actualQuota), logger.LogQuota(preConsumed))
	case delta < 0:
		return common.LogText("refunding the difference after pre-consume: %s (actual %s, pre-consumed %s)",
			logger.LogQuota(-delta), logger.LogQuota(actualQuota), logger.LogQuota(preConsumed))
	case taskSettlesOnCompletion(relayInfo):
		return common.LogText("task pre-deducted %s, pending settlement: the final amount is settled from actual usage when the task completes",
			logger.LogQuota(actualQuota))
	case billedPerRequest(relayInfo):
		return common.LogText("pre-consumed quota matches the actual quota, no adjustment: %s (per-request billing)",
			logger.LogQuota(actualQuota))
	default:
		return common.LogText("pre-consumed quota matches the actual quota, no adjustment: %s",
			logger.LogQuota(actualQuota))
	}
}

// SettleBilling 执行计费结算。如果 RelayInfo 上有 BillingSession 则通过 session 结算，
// 否则回退到旧的 PostConsumeQuota 路径（兼容按次计费等场景）。
func SettleBilling(ctx *gin.Context, relayInfo *relaycommon.RelayInfo, actualQuota int) error {
	if relayInfo.Billing != nil {
		preConsumed := relayInfo.Billing.GetPreConsumedQuota()
		logger.LogInfo(ctx, settleBillingLogLine(relayInfo, actualQuota, preConsumed))

		if err := relayInfo.Billing.Settle(actualQuota); err != nil {
			return err
		}

		// 发送额度通知（订阅计费使用订阅剩余额度）
		if actualQuota != 0 {
			if relayInfo.BillingSource == BillingSourceSubscription {
				checkAndSendSubscriptionQuotaNotify(relayInfo)
			} else {
				checkAndSendQuotaNotify(relayInfo, actualQuota-preConsumed, preConsumed)
			}
		}
		return nil
	}

	// 回退：无 BillingSession 时使用旧路径
	quotaDelta := actualQuota - relayInfo.FinalPreConsumedQuota
	if quotaDelta != 0 {
		return PostConsumeQuota(relayInfo, quotaDelta, relayInfo.FinalPreConsumedQuota, true)
	}
	return nil
}

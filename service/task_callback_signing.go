package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/setting/system_setting"
)

// Optional task callback signing (kt-api extension, off by default; the
// supplier and Ark documents define no signature).
const (
	TaskCallbackHeaderTimestamp = "X-KtApi-Timestamp"
	TaskCallbackHeaderSignature = "X-KtApi-Signature"
	taskCallbackSignatureScheme = "v1"
	taskCallbackSecretContext   = "kt-api-task-callback-signing:"
)

// TaskCallbackUserSecret derives a user's signing secret from the master
// signing secret: hex(HMAC-SHA256(master, context + user id)). Nothing is
// stored per user; rotating the master secret rotates every user secret.
func TaskCallbackUserSecret(masterSecret string, userID int) string {
	if strings.TrimSpace(masterSecret) == "" || userID <= 0 {
		return ""
	}
	mac := hmac.New(sha256.New, []byte(masterSecret))
	mac.Write([]byte(taskCallbackSecretContext + strconv.Itoa(userID)))
	return hex.EncodeToString(mac.Sum(nil))
}

// TaskCallbackSignature returns v1=hex(HMAC-SHA256(secret, timestamp + "." +
// event_id + "." + body)).
func TaskCallbackSignature(secret string, timestamp int64, eventID string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(timestamp, 10)))
	mac.Write([]byte("."))
	mac.Write([]byte(eventID))
	mac.Write([]byte("."))
	mac.Write(body)
	return taskCallbackSignatureScheme + "=" + hex.EncodeToString(mac.Sum(nil))
}

// VerifyTaskCallbackSignature checks a signature header in constant time
// (reference implementation for customers and tests).
func VerifyTaskCallbackSignature(secret string, timestamp int64, eventID string, body []byte, header string) bool {
	expected := TaskCallbackSignature(secret, timestamp, eventID, body)
	for _, candidate := range strings.Split(header, ",") {
		if hmac.Equal([]byte(strings.TrimSpace(candidate)), []byte(expected)) {
			return true
		}
	}
	return false
}

// CurrentTaskCallbackUserSecret returns the user's signing secret, or ""
// when signing is disabled or no master secret is configured.
func CurrentTaskCallbackUserSecret(userID int) string {
	setting := system_setting.GetTaskCallbackSetting()
	if !setting.SigningEnabled {
		return ""
	}
	return TaskCallbackUserSecret(setting.SigningSecret, userID)
}

// signTaskCallbackRequest adds the timestamp and signature headers when
// signing is enabled and a master secret is configured.
func signTaskCallbackRequest(request *http.Request, setting *system_setting.TaskCallbackSetting, userID int, eventID string, body []byte, timestamp int64) {
	if !setting.SigningEnabled {
		return
	}
	secret := TaskCallbackUserSecret(setting.SigningSecret, userID)
	if secret == "" {
		return
	}
	request.Header.Set(TaskCallbackHeaderTimestamp, strconv.FormatInt(timestamp, 10))
	request.Header.Set(TaskCallbackHeaderSignature, TaskCallbackSignature(secret, timestamp, eventID, body))
}

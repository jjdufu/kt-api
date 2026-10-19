package setting

import (
	"encoding/json"

	"github.com/QuantumNous/new-api/common"
)

// Chats is the list of chat client presets shown to users. It is empty by
// default; administrators add the clients they want to offer.
var Chats = []map[string]string{}

func UpdateChatsByJsonString(jsonString string) error {
	Chats = make([]map[string]string, 0)
	return json.Unmarshal([]byte(jsonString), &Chats)
}

func Chats2JsonString() string {
	jsonBytes, err := json.Marshal(Chats)
	if err != nil {
		common.SysLog(common.LogText("error marshalling chats: %s", err.Error()))
		return "[]"
	}
	return string(jsonBytes)
}

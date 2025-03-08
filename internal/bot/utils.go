package bot

import (
	"strings"
	"time"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func currentDateString() string {
	return time.Now().Format("2006-01-02")
}

func parseNumber(jid types.JID) string {
	number := strings.Split(jid.ToNonAD().String(), "@")[0]
	return "+" + number
}

// extractText picks up the text from the different places a message might hold it
func extractText(msgEvt *events.Message) string {
	if msgEvt.Message.GetConversation() != "" {
		return msgEvt.Message.GetConversation()
	} else if msgEvt.Message.ExtendedTextMessage != nil {
		return msgEvt.Message.ExtendedTextMessage.GetText()
	}
	return ""
}

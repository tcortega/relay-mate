package bot

import (
	"log"

	"go.mau.fi/whatsmeow/types/events"
)

// eventHandler routes incoming events from the WhatsApp client.
func (b *Bot) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		b.handleMessage(v)
	}
}

func (b *Bot) handleMessage(msgEvt *events.Message) {
	if msgEvt.Info.IsGroup {
		return
	}

	sender := msgEvt.Info.Sender
	if sender == b.LoggedInJID || sender == b.ForwardToJID {
		return
	}
	text := extractText(msgEvt)

	// Collect contact name info
	contactName := ""
	pushName := msgEvt.Info.PushName
	contactInfo, err := b.Client.Store.Contacts.GetContact(sender)
	if err != nil {
		log.Printf("Error getting contact info: %v", err)
	} else {
		contactName = contactInfo.FullName
		if pushName == "" {
			pushName = contactInfo.PushName
		}
	}

	// Possibly send the "starter" message if needed
	if err := b.sendStarterIfNeeded(sender, msgEvt.Info.Timestamp); err != nil {
		log.Printf("Error sending starter message: %v", err)
	}

	contactNameLine := ""
	if contactName != "" {
		contactNameLine = "👤 *Contato:* " + contactName + "\n"
	}
	pushNameLine := ""
	if pushName != "" {
		pushNameLine = "🔖 *Nome:* " + pushName + "\n"
	}

	forwardText := "📬 *Nova mensagem recebida!*\n\n> " + text + "\n\n" +
		contactNameLine + pushNameLine +
		"📱 *Número:* " + parseNumber(sender)

	// Try to get the sender's profile picture
	log.Printf("Getting profile picture info for %s", sender)
	profilePictureInfo, err := b.Client.GetProfilePictureInfo(sender, nil)
	if err != nil {
		log.Printf("Error getting profile picture info: %v", err)
	}

	// If a profile picture exists, forward it as an image message
	if profilePictureInfo != nil && profilePictureInfo.URL != "" {
		if err := b.sendImageMessage(b.ForwardToJID, forwardText, profilePictureInfo.URL); err != nil {
			log.Printf("Error sending image message: %v", err)
		}
	} else {
		if err := b.sendTextMessage(b.ForwardToJID, forwardText); err != nil {
			log.Printf("Error sending text message: %v", err)
		}
	}
}

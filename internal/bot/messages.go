package bot

import (
	"context"
	"database/sql"
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net/http"
	"time"
)

// sendStarterIfNeeded checks if we've already sent the starter message to this
// contact today. If not, it sends the message and updates the DB.
// Messages older than 5 minutes are ignored.
func (b *Bot) sendStarterIfNeeded(jid types.JID, timestamp time.Time) error {
	currentTime := time.Now()
	if currentTime.Sub(timestamp) > 5*time.Minute {
		return nil
	}

	_, err := b.DB.Exec(`
		INSERT OR IGNORE INTO contacts (jid, last_starter_sent)
		VALUES (?, NULL)
	`, jid.String())
	if err != nil {
		return err
	}

	// Check when we last sent the starter
	var lastSentVal sql.NullString
	err = b.DB.QueryRow(`
		SELECT last_starter_sent
		FROM contacts
		WHERE jid = ?
	`, jid.String()).Scan(&lastSentVal)

	// If some other error
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	today := currentDateString()

	// If lastSent is nil or different from today, send the starter
	if !lastSentVal.Valid || lastSentVal.String != today {
		if err := b.sendTextMessage(jid, b.Config.StarterMessage); err != nil {
			return err
		}
		_, err = b.DB.Exec(`
			UPDATE contacts
			SET last_starter_sent = ?
			WHERE jid = ?
		`, today, jid.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) sendTextMessage(to types.JID, message string) error {
	_, err := b.Client.SendMessage(context.Background(), to, &waE2E.Message{
		Conversation: proto.String(message),
	})
	return err
}

func (b *Bot) sendImageMessage(to types.JID, caption string, imageUrl string) error {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return fmt.Errorf("failed to download image from %s: %w", imageUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch image from %s, response code: %d", imageUrl, resp.StatusCode)
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read image body: %w", err)
	}

	return b.sendImageMessageBytes(to, caption, imageBytes)
}

func (b *Bot) sendImageMessageBytes(to types.JID, caption string, imageBytes []byte) error {
	if len(imageBytes) == 0 {
		return fmt.Errorf("imageBytes is empty")
	}

	mimeType := http.DetectContentType(imageBytes)
	uploadResp, err := b.Client.Upload(context.Background(), imageBytes, whatsmeow.MediaImage)
	if err != nil {
		log.Printf("Failed to upload image: %v", err)
		return err
	}

	imageMsg := &waE2E.ImageMessage{
		Caption:       proto.String(caption),
		Mimetype:      proto.String(mimeType),
		URL:           &uploadResp.URL,
		DirectPath:    &uploadResp.DirectPath,
		MediaKey:      uploadResp.MediaKey,
		FileEncSHA256: uploadResp.FileEncSHA256,
		FileSHA256:    uploadResp.FileSHA256,
		FileLength:    &uploadResp.FileLength,
	}

	resp, err := b.Client.SendMessage(context.Background(), to, &waE2E.Message{
		ImageMessage: imageMsg,
	})
	if err != nil {
		log.Printf("Failed to send image message: %v", err)
		return err
	}

	log.Printf("Image message sent successfully: %+v", resp)
	return nil
}

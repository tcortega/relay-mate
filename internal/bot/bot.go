package bot

import (
	"database/sql"
	"fmt"
	"log"
	"relay-mate/internal/config"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Bot struct {
	Client       *whatsmeow.Client
	DB           *sql.DB
	Config       config.Config
	LoggedInJID  types.JID
	ForwardToJID types.JID
}

// NewBot constructs a Bot. It sets up the client but does NOT connect yet.
func NewBot(cfg config.Config, db *sql.DB, deviceStore *store.Device) (*Bot, error) {
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	forwardJID, err := types.ParseJID(cfg.ForwardTo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse forward_to JID: %w", err)
	}

	bot := &Bot{
		Client:       client,
		DB:           db,
		Config:       cfg,
		ForwardToJID: forwardJID,
	}

	bot.Client.AddEventHandler(bot.eventHandler)

	return bot, nil
}

// Start connects the client to WhatsApp. If client.Store.ID == nil, the user
// needs to scan a QR code (handled by main).
func (b *Bot) Start() error {
	err := b.Client.Connect()
	if err != nil {
		return err
	}

	if b.Client.Store.ID != nil {
		b.LoggedInJID = b.Client.Store.ID.ToNonAD()
	}
	return nil
}

// Stop cleanly disconnects the client.
func (b *Bot) Stop() {
	b.Client.Disconnect()
	log.Println("Client disconnected")
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"relay-mate/internal/bot"
	"relay-mate/internal/config"
	"relay-mate/internal/database"
	"syscall"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	db, err := database.InitDB("file:botstore.db?_foreign_keys=on")
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:botstore.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatalf("Error creating SQLStore container: %v", err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Fatalf("Error getting device store: %v", err)
	}

	wabot, err := bot.NewBot(cfg, db, deviceStore)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	if deviceStore.ID == nil {
		qrChan, _ := wabot.Client.GetQRChannel(context.Background())
		if err := wabot.Start(); err != nil {
			log.Fatalf("Failed to start bot: %v", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Printf("QR Event: %s\n", evt.Event)
			}
		}
	} else {
		if err := wabot.Start(); err != nil {
			log.Fatalf("Failed to start bot: %v", err)
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	fmt.Println("Shutting down bot...")
	wabot.Stop()
}

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/fsnotify.v1"
	"gopkg.in/gomail.v2"
)

func sendEmailWithAttachment(filePath string) error {
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	toEmail := os.Getenv("EMAIL_TO")

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", username)
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", fmt.Sprintf("New File: %s", filepath.Base(filePath)))
	mailer.SetBody("text/plain", fmt.Sprintf("A new file has been added: %s", filepath.Base(filePath)))

	mailer.Attach(filePath)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, username, password)

	if err := dialer.DialAndSend(mailer); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

func main() {
	dirToWatch := os.Getenv("/books")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
	}
	defer watcher.Close()

	err = watcher.Add(dirToWatch)
	if err != nil {
		log.Fatal("Error adding directory to watcher:", err)
	}

	fmt.Printf("Watching for new files in %s...\n", dirToWatch)

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("New file created:", event.Name)

				time.Sleep(2 * time.Second)

				err := sendEmailWithAttachment(event.Name)
				if err != nil {
					log.Println("Error sending email:", err)
				} else {
					fmt.Printf("Email sent with attachment: %s\n", event.Name)
				}
			}
		case err := <-watcher.Errors:
			log.Println("Watcher error:", err)
		}
			time.Sleep(20 * time.Second)
	}
}

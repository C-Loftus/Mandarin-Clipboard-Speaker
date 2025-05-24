package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"unicode"

	"golang.design/x/clipboard"
)

func getUserHomeDir() string {
	// Attempt to get the home directory using os.UserHomeDir()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to using the current user information
		currentUser, err := user.Current()
		if err != nil || currentUser.HomeDir == "" {
			log.Fatalf("failed to get user home directory: %v", err)
		}
		homeDir = currentUser.HomeDir
	}
	return homeDir
}

// Download piper and the chinese model if it doesn't exist
// we use the same directory as QuickPiperAudiobook in order to share models
func downloadPiper() {
	homeDir := getUserHomeDir()
	piperPath := filepath.Join(homeDir, ".config", "QuickPiperAudiobook")
	modelFile := filepath.Join(piperPath, "zh_CN-huayan-medium.onnx")
	modelJSON := filepath.Join(piperPath, "zh_CN-huayan-medium.onnx.json")
	piperBinary := filepath.Join(piperPath, "piper", "piper")

	// Check if model and binary exist
	if _, err := os.Stat(modelFile); err == nil {
		if _, err := os.Stat(piperBinary); err == nil {
			log.Println("Piper and model already exist, skipping download.")
			return
		}
	}

	log.Print("Downloading Piper for text-to-speech...")

	// Create piper directory if it doesn't exist
	if err := os.MkdirAll(piperPath, 0755); err != nil {
		log.Fatalf("failed to create piper directory: %v", err)
	}

	// Download Piper tar.gz
	tarPath := filepath.Join(piperPath, "piper_amd64.tar.gz")
	cmd := exec.Command("wget", "-O", tarPath, "https://github.com/rhasspy/piper/releases/download/v1.2.0/piper_amd64.tar.gz")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to download Piper binary: %v", err)
	}

	// Extract Piper tar.gz
	cmd = exec.Command("tar", "-xvf", tarPath, "-C", piperPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to extract Piper binary: %v", err)
	}

	// Download model files
	modelURL := "https://huggingface.co/rhasspy/piper-voices/resolve/main/zh/zh_CN/huayan/medium/zh_CN-huayan-medium.onnx?download=true"
	cmd = exec.Command("wget", "-O", modelFile, modelURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to download ONNX model: %v", err)
	}

	modelJSONURL := "https://huggingface.co/rhasspy/piper-voices/resolve/main/zh/zh_CN/huayan/medium/zh_CN-huayan-medium.onnx.json?download=true"
	cmd = exec.Command("wget", "-O", modelJSON, modelJSONURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to download ONNX model JSON: %v", err)
	}

	log.Println("Piper setup completed successfully.")
}

func main() {

	downloadPiper()

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("failed to get user home directory: %w", err)
	}
	piperPath := filepath.Join(homedir, "piper")

	go func() {
		err := clipboard.Init()
		if err != nil {
			panic(err)
		}

		err = os.Chdir(piperPath)
		if err != nil {
			log.Fatalf("Failed to change directory: %v", err)
		}

		var cancelCurrent context.CancelFunc

		changed := clipboard.Watch(context.Background(), clipboard.FmtText)
		for range changed {
			text := clipboard.Read(clipboard.FmtText)
			textWithJustHan := ""

			for _, char := range string(text) {
				if unicode.Is(unicode.Han, char) {
					textWithJustHan += string(char)
				}
			}
			if textWithJustHan == "" {
				log.Println("clipboard changed, but no han characters found to speak")
				continue
			}

			// Cancel the previous speaking process
			if cancelCurrent != nil {
				cancelCurrent()
			}

			// Create a new context for the current speaking process
			ctx, cancel := context.WithCancel(context.Background())
			cancelCurrent = cancel

			// Launch speaking process
			go func(ctx context.Context, hanText string) {
				log.Printf("Speaking: %s", hanText)

				piperCmd := exec.CommandContext(ctx, "./piper/piper", "--model", "zh_CN-huayan-medium.onnx", "--output-raw")
				piperIn, err := piperCmd.StdinPipe()
				if err != nil {
					log.Printf("Failed to get Piper stdin: %v", err)
					return
				}
				piperOut, err := piperCmd.StdoutPipe()
				if err != nil {
					log.Printf("Failed to get Piper stdout: %v", err)
					return
				}

				aplayCmd := exec.CommandContext(ctx, "aplay", "-r", "22050", "-f", "S16_LE", "-t", "raw", "-")
				aplayIn, err := aplayCmd.StdinPipe()
				if err != nil {
					log.Printf("Failed to get aplay stdin: %v", err)
					return
				}

				if err := piperCmd.Start(); err != nil {
					log.Printf("Failed to start Piper: %v", err)
					return
				}
				if err := aplayCmd.Start(); err != nil {
					log.Printf("Failed to start aplay: %v", err)
					return
				}

				// Send text to Piper
				_, err = fmt.Fprintln(piperIn, hanText)
				if err != nil {
					log.Printf("Error writing to Piper stdin: %v", err)
				}
				piperIn.Close()

				// Pipe Piper output to aplay
				_, err = bufio.NewReader(piperOut).WriteTo(aplayIn)
				if err != nil {
					log.Printf("Error piping Piper output to aplay: %v", err)
				}
				aplayIn.Close()

				// Wait for processes to finish
				piperCmd.Wait()
				aplayCmd.Wait()

			}(ctx, textWithJustHan)
		}
	}()

	fmt.Print("Listing on clipboard. Press Ctrl+C to exit.\n")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("\nExiting.")

}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"
	"youchat/chat/config"

	"github.com/briandowns/spinner"
	"google.golang.org/genai"
)

var allowedBlobMIMETypes = map[string]bool{
	"application/pdf": true,
	"audio/mpeg":      true,
	"audio/mp3":       true,
	"audio/wav":       true,
	"image/png":       true,
	"image/jpeg":      true,
	"image/webp":      true,
	"video/mov":       true,
	"video/mpeg":      true,
	"video/mp4":       true,
	"video/mpg":       true,
	"video/avi":       true,
	"video/wmv":       true,
	"video/mpegps":    true,
	"video/flv":       true,
}

func main() {
	// Define flags for file input and direct question
	filePath := flag.String("file", "", "Path to a file (text, image, etc.) to analyze")
	directQuestion := flag.String("q", "", "Direct question to ask")

	flag.Parse()

	var userParts []*genai.Part
	var userQuestionText string

	if *filePath != "" {
		fileBytes, err := os.ReadFile(*filePath)
		if err != nil {
			log.Fatalf("Error reading file %s: %v", *filePath, err)
		}

		mimeType := mime.TypeByExtension(filepath.Ext(*filePath))
		if mimeType == "" {
			fmt.Printf("Warning: Could not determine MIME type for %s. Defaulting to text/plain.\n", *filePath)
			mimeType = "text/plain"
		}

		fmt.Printf("Analyzing file: %s (MIME type: %s)\n", *filePath, mimeType)

		// Check if the current MIME type is in our allowedBlobMIMETypes map
		if allowedBlobMIMETypes[mimeType] {
			blobPart := &genai.Part{
				InlineData: &genai.Blob{
					MIMEType: mimeType,
					Data:     fileBytes,
				},
			}
			userParts = append(userParts, blobPart)
		} else {
			// It's NOT an allowed blob type, so treat it as text and format it
			fileContentString := string(fileBytes)
			fileName := filepath.Base(*filePath)

			// Format the text as "file - [filename], [filecontent]"
			formattedFileText := fmt.Sprintf("file - [%s], [%s]", fileName, fileContentString)
			userParts = append(userParts, genai.NewPartFromText(formattedFileText))
		}

		if *directQuestion != "" {
			userQuestionText = *directQuestion
		} else {
			userQuestionText = "What is in this file?"
		}

	} else if *directQuestion != "" {
		userQuestionText = *directQuestion
	} else {
		if len(flag.Args()) > 0 {
			userQuestionText = strings.Join(flag.Args(), " ")
		} else {
			fmt.Println("Please provide a question to ask, a file using --file, or use the -q flag.")
			printHelp()
			os.Exit(1)
		}
	}

	if userQuestionText != "" {
		userParts = append(userParts, genai.NewPartFromText(userQuestionText))
	}

	// Ensure we have at least one part (either text or file)
	if len(userParts) == 0 {
		fmt.Println("No question or file content provided. Exiting.")
		printHelp()
		os.Exit(1)
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(userParts, genai.RoleUser),
	}

	if len(flag.Args()) > 0 && flag.Args()[0] == "help" && *filePath == "" && *directQuestion == "" {
		printHelp()
		return
	} else if flag.Args()[0] == "set" {
		if len(flag.Args()) < 2 {
			fmt.Println("Usage:")
			fmt.Println("  set <your-api-key> - obtain one from https://aistudio.google.com (free)")
		} else {
			api_key := flag.Args()[1]
			cfg, err := config.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				os.Exit(1)
			}
			cfg.APIKey = api_key

			if err := config.SaveConfig(cfg); err != nil {
				fmt.Printf("Error saving API key: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("API key set successfully.")

			// Get the config path again to display the full file location
			configDirPath, err := config.GetConfigPath()
			if err != nil {
				fmt.Printf("Error getting config directory path: %v\n", err)
				os.Exit(1)
			}
			configFilePath := filepath.Join(configDirPath, config.ConfigFileName)
			fmt.Printf("Config file location: %s\n", configFilePath)
		}
		return
	}

	handleAskCommand(contents)
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  chat [options] <your_question>")
	fmt.Println("  chat --file <path_to_file> [--q \"Your question about the file\"]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  --file <path>   Provides a file (text, image, etc.) for analysis. Required for multimodal input.")
	fmt.Println("  -q <question>   Directly provides a question (alternative to just typing the question).")
	fmt.Println("                  When used with --file, this question is asked *about* the file.")
	fmt.Println("")
	fmt.Println("Available Commands:")
	fmt.Println("  <your_question>   Asks a question to the Gemini 2.5 Flash Lite model.")
	fmt.Println("                    Example: chat \"What is the capital of France?\"")
	fmt.Println("  --file            Uploads a file to the model for analysis. Can be combined with a question.")
	fmt.Println("                    Example: chat --file ./my_document.txt \"Summarize this document.\"")
	fmt.Println("                    Example: chat --file ./image.jpg \"What is in this picture?\"")
	fmt.Println("  help              Shows this help message.")
	fmt.Println("  set <api-key>     Set your gemini api key. obtain one from https://aistudio.google.com (free)")
}

func handleAskCommand(contents []*genai.Content) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading API key: %v\n", err)
		os.Exit(1)
	}
	if cfg.APIKey == "" {
		fmt.Println("API key not set.")
		os.Exit(1)
	}
	apiKey := cfg.APIKey

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	modelName := "gemini-2.5-flash-lite-preview-06-17"

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText("You are an assistant in the terminal that provides direct answers to users' questions without any further context or filler text. You are to generate short and precise answers to users' questions. Aim for less than 100 words. Also, remember you're in the terminal, so avoid markdown or similar formatting. If no specific context is provided, assume it's a question relating to something in the command line. Also note that you're a one-way chat, meaning the user cannot provide additional context after your first reply. When you're given a file though, be consice with your reply, but not soo that you are not given any meaning to the user. Also the user might request a specific thing for thier output, give it to them as requested. you're also a bit of savage when the user is not requesting serious and it just casual conversation 😂", genai.RoleUser),
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	stream := client.Models.GenerateContentStream(
		ctxWithTimeout,
		modelName,
		contents,
		config,
	)

	s := spinner.New(spinner.CharSets[28], 100*time.Millisecond)
	s.Start()

	answer := ""
	for chunk, _ := range stream {
		if len(chunk.Candidates) > 0 {
			part := chunk.Candidates[0].Content.Parts[0]
			answer += part.Text

			s.Suffix = answer
		}

	}
	s.Stop()
	fmt.Println(answer)
}

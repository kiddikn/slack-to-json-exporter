package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

func main() {
	var (
		token     = flag.String("token", "", "Slack Bot OAuth Token (required)")
		channels  = flag.String("channels", "", "Comma-separated list of channel IDs (required)")
		startDate = flag.String("startDate", "", "Start date in YYYY-MM-DD format (required)")
		endDate   = flag.String("endDate", "", "End date in YYYY-MM-DD format (required)")
		outputDir = flag.String("outputDir", "./slack_logs", "Output directory for JSON files")
	)
	flag.Parse()

	if *token == "" || *channels == "" || *startDate == "" || *endDate == "" {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	startTimestamp, err := parseDate(*startDate, true)
	if err != nil {
		log.Fatalf("Error parsing start date: %v", err)
	}

	endTimestamp, err := parseDate(*endDate, false)
	if err != nil {
		log.Fatalf("Error parsing end date: %v", err)
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	api := slack.New(*token)

	channelList := strings.Split(*channels, ",")
	for _, channelID := range channelList {
		channelID = strings.TrimSpace(channelID)
		if channelID == "" {
			continue
		}

		fmt.Printf("Fetching logs for channel %s...\n", channelID)

		messages, err := fetchAllMessages(api, channelID, startTimestamp, endTimestamp)
		if err != nil {
			log.Printf("Error fetching messages for channel %s: %v", channelID, err)
			continue
		}

		filename := filepath.Join(*outputDir, channelID+".json")
		if err := saveMessagesToFile(messages, filename); err != nil {
			log.Printf("Error saving messages to file %s: %v", filename, err)
			continue
		}

		fmt.Printf("Successfully saved logs to %s\n", filename)
	}
}

func parseDate(dateStr string, isStartOfDay bool) (string, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	if isStartOfDay {
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	} else {
		t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
	}

	return fmt.Sprintf("%d", t.Unix()), nil
}

func fetchAllMessages(api *slack.Client, channelID, startTime, endTime string) ([]slack.Message, error) {
	var allMessages []slack.Message
	cursor := ""

	for {
		params := &slack.GetConversationHistoryParameters{
			ChannelID: channelID,
			Oldest:    startTime,
			Latest:    endTime,
			Cursor:    cursor,
		}

		response, err := api.GetConversationHistory(params)
		if err != nil {
			return nil, err
		}

		allMessages = append(allMessages, response.Messages...)

		if !response.HasMore {
			break
		}

		cursor = response.ResponseMetaData.NextCursor
	}

	return allMessages, nil
}

func saveMessagesToFile(messages []slack.Message, filename string) error {
	jsonData, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}
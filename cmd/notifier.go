package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func SlackNotifier(d *RbData) error {

	announcement := fmt.Sprintf(`{
    "attachments": [
        {
            "title": "%s",
            "pretext": "New Rating for %s!",
            "text": "Timestamp: %s\nUser: %s\nRating: %s\nComment: %s",
			"title_link": "%s",
            "mrkdwn_in": [
                "text",
                "pretext"
            ]
        }
    ]
}`, d.Location, d.Location, d.Time.Format("2006-01-02"), d.User, d.Rating, d.Comment, d.Link)
	var jsonStr = []byte(announcement)
	req, err := http.NewRequest("POST", config.SlackWebHook, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Sending Slack Message failed", err)
	}
	defer resp.Body.Close()

	log.Println("Status: %v Headers: %s", resp.Status)

	return nil
}

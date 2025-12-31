package aiChat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type WebhookPayload struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID        string           `json:"id"`
	Time      int64            `json:"time"`
	Messaging []MessagingEvent `json:"messaging"`
}

type MessagingEvent struct {
	Sender    IDObj     `json:"sender"`
	Recipient IDObj     `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   *Message  `json:"message,omitempty"`
	Postback  *Postback `json:"postback,omitempty"`
}

type IDObj struct {
	ID string `json:"id"`
}

type Message struct {
	MID    string `json:"mid,omitempty"`
	Text   string `json:"text,omitempty"`
	IsEcho bool   `json:"is_echo,omitempty"`
}

type Postback struct {
	Title   string `json:"title,omitempty"`
	Payload string `json:"payload,omitempty"`
}

// --- Send API helper ---

type sendRequest struct {
	MessagingType string `json:"messaging_type,omitempty"`
	Recipient     struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}

func SendText(psid, text string) error {
	pageToken := os.Getenv("FB_PAGE_ACCESS_TOKEN")
	if pageToken == "" {
		return fmt.Errorf("FB_PAGE_ACCESS_TOKEN not set")
	}

	// You can set this to whatever Graph API version you’re using.
	// If unset, we default to a reasonable one.
	apiVer := os.Getenv("FB_GRAPH_API_VERSION")
	if apiVer == "" {
		apiVer = "v20.0"
	}

	reqBody := sendRequest{
		MessagingType: "RESPONSE",
	}
	reqBody.Recipient.ID = psid
	reqBody.Message.Text = text

	b, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/me/messages?access_token=%s", apiVer, pageToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("send api failed: status=%s body=%s", resp.Status, string(respBytes))
	}

	return nil
}

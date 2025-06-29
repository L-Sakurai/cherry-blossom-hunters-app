package notify

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Notifier struct {
    WebhookURL string
}

func NewDiscordNotifier(webhookURL string) *Notifier {
    return &Notifier{WebhookURL: webhookURL}
}

type CustomPayload struct {
    Title     string
    Message   string
    Timestamp string
    Fields    map[string]interface{}
}

type DiscordWebhookPayload struct {
    Embeds []Embed `json:"embeds"`
}

type Embed struct {
    Title       string       `json:"title,omitempty"`
    Timestamp   string       `json:"timestamp,omitempty"`
    Fields      []EmbedField `json:"fields,omitempty"`
    Color       int          `json:"color,omitempty"`
    Description string       `json:"description,omitempty"`
}

type EmbedField struct {
    Name   string `json:"name"`
    Value  string `json:"value"`
    Inline bool   `json:"inline"`
}

func (n *Notifier) Send(p *CustomPayload) error {
    if p.Timestamp == "" {
        p.Timestamp = time.Now().UTC().Format(time.RFC3339)
    }

    fields := []EmbedField{
        {
            Name:   "メッセージ",
            Value:  p.Message,
            Inline: false,
        },
    }

    for k, v := range p.Fields {
        fields = append(fields, EmbedField{
            Name:   k,
            Value:  fmt.Sprintf("%v", v),
            Inline: false,
        })
    }

    payload := DiscordWebhookPayload{
        Embeds: []Embed{
            {
                Title:     p.Title,
                Timestamp: p.Timestamp,
                Fields:    fields,
                Color:     0x3498db,
            },
        },
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }

    req, err := http.NewRequest("POST", n.WebhookURL, bytes.NewBuffer(body))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send webhook request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("webhook returned status %d", resp.StatusCode)
    }

    return nil
}

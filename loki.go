package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/sirupsen/logrus"
)

type LokiHook struct {
    url    string
    labels map[string]string
    client *http.Client
}

type lokiStream struct {
    Stream map[string]string `json:"stream"`
    Values [][2]string       `json:"values"`
}

type lokiPayload struct {
    Streams []lokiStream `json:"streams"`
}

func NewLokiHook(url string, labels map[string]string) *LokiHook {
    return &LokiHook{
        url:    url,
        labels: labels,
        client: &http.Client{Timeout: 5 * time.Second},
    }
}

func (h *LokiHook) Levels() []logrus.Level {
    return logrus.AllLevels
}

func (h *LokiHook) Fire(entry *logrus.Entry) error {
    line, err := entry.String()
    if err != nil {
        return err
    }

    payload := lokiPayload{
        Streams: []lokiStream{
            {
                Stream: h.labels,
                Values: [][2]string{
                    {
                        fmt.Sprintf("%d", entry.Time.UnixNano()),
                        line,
                    },
                },
            },
        },
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    resp, err := h.client.Post(
        h.url,
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}

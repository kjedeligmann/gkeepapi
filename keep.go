package gkeepapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const APIURL = "https://www.googleapis.com/notes/v1/"

type Keep struct {
	Auth
	sessionId string
}

type ResponseBody struct {
	Nodes []Node
}

type Node struct {
	Type       string
	Id         string
	ParentId   string
	Timestamps Timestamps
	Title      string
	Text       string
}

type Timestamps struct {
	Created string
}

type Note struct {
	Created string
	Title   string
	Text    string
}

func (s *Keep) Authenticate(email, gaid, masterToken string) error {
	s.sessionId = generateID(time.Now().UTC())
	err := s.Load(email, gaid, masterToken)
	if err != nil {
		return err
	}
	return nil
}

// List all notes.
func (s *Keep) List() (map[string]Note, error) {
	b, err := s.changes()
	if err != nil {
		return nil, err
	}

	var nodes ResponseBody
	if err := json.Unmarshal(b, &nodes); err != nil {
		return nil, fmt.Errorf("JSON unmarshaling failed: %s", err)
	}

	notes := make(map[string]Note)
	for _, node := range nodes.Nodes {
		switch node.Type {
		case "NOTE":
			if note, ok := notes[node.Id]; ok {
				note.Created = node.Timestamps.Created
				note.Title = node.Title
				notes[node.Id] = note
			} else {
				notes[node.Id] = Note{
					Created: node.Timestamps.Created,
					Title:   node.Title,
				}
			}
		case "LIST_ITEM":
			if note, ok := notes[node.ParentId]; ok {
				note.Text = node.Text
				notes[node.ParentId] = note
			} else {
				notes[node.ParentId] = Note{Text: node.Text}
			}
		}
	}
	return notes, nil
}

func (s *Keep) changes() ([]byte, error) {
	jsonBody := []byte(`
{
    "nodes": [],
    "clientTimestamp": "` + time.Now().UTC().Format(time.RFC3339Nano) + `",
    "requestHeader": {
        "clientSessionId": "` + s.sessionId + `",
        "clientPlatform": "ANDROID",
        "clientVersion": {
            "major": "9",
            "minor": "9",
            "build": "9",
            "revision": "9"
        },
        "capabilities": [
            {
                "type": "NC"
            },
            {
                "type": "PI"
            },
            {
                "type": "LB"
            },
            {
                "type": "AN"
            },
            {
                "type": "SH"
            },
            {
                "type": "DR"
            },
            {
                "type": "TR"
            },
            {
                "type": "IN"
            },
            {
                "type": "SNB"
            },
            {
                "type": "MI"
            },
            {
                "type": "CO"
            }
        ]
    }
}
    `)
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, APIURL+"changes", bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "OAuth "+s.accessToken)
	req.Header.Set("User-Agent", "x-gkeepapi (https://github.com/kjedeligmann/gkeepapi)")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func generateID(t time.Time) string {
	return fmt.Sprintf(
		"s--%d--%d",
		t.Unix()*1000,
		rand.Intn(8999999999)+1000000000,
	)
}

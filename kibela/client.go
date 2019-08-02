package kibela

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/xerrors"
)

// Base API endpoint. Team name is going to replace subdomain.
const apiEndpoint string = "https://%s.kibe.la/api/v1"

// Kibela is struct to call Kibela GraphQL API which includes API token, API endpoint which is
// determined by team name, and HTTP client.
type Kibela struct {
	token    string
	endpoint string
	client   *http.Client
}

// New creates new Kibela client then return its pointer.
func New(token, team string) *Kibela {
	return &Kibela{
		token:    token,
		endpoint: fmt.Sprintf(apiEndpoint, team),
		client:   &http.Client{},
	}
}

// request represents Kibela GraphQL request.
type request struct {
	Query string `json:"query"`
}

// response represents Kibela GraphQL API response. Because response format is determined based on
// the request, the type of Data is `json.RawMessage` to delay determining actual type.
type response struct {
	Data json.RawMessage `json:"data"`
}

// GetNote finds Kibela note of ID then return its pointer
func (k *Kibela) GetNote(noteID string) (*Note, error) {
	resp, err := k.do(&request{Query: getNoteQuery(noteID)})
	if err != nil {
		return nil, err
	}
	var r struct {
		Note *Note `json:"note"`
	}
	if err := json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r.Note, nil
}

// AddLink appends markdown link to the note of ID.
func (k *Kibela) AddLink(noteID, url, title string) (*Note, error) {
	// Get base note
	note, err := k.GetNote(noteID)
	if err != nil {
		return nil, xerrors.Errorf("failed to get note of id=%s. err=%s", noteID, err)
	}

	// Create new content that has new markdown link
	newContent := note.Content
	if !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}
	newContent += fmt.Sprintf("- [%s](%s)", title, url)

	// API call
	resp, err := k.do(&request{Query: updateNoteMutation(note.ID, note.Content, newContent)})
	if err != nil {
		return nil, err
	}
	var r struct {
		UpdateNoteContent struct {
			Note *Note `json:"note"`
		} `json:"updateNoteContent"`
	}
	if err := json.Unmarshal(resp.Data, &r); err != nil {
		return nil, xerrors.Errorf("failed to unmarshaling response JSON. err=%s", err)
	}
	return r.UpdateNoteContent.Note, nil
}

func (k *Kibela) do(reqBody *request) (*response, error) {
	body := bytes.Buffer{}
	if err := json.NewEncoder(&body).Encode(reqBody); err != nil {
		return nil, xerrors.Errorf("failed to create request JSON from %v", reqBody)
	}
	req, err := http.NewRequest(http.MethodPost, k.endpoint, &body)
	if err != nil {
		return nil, xerrors.Errorf("failed to make new HTTP request. err=%s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.token))
	req.Header.Set("Content-Type", "application/json") // TODO: application/x-msgpack
	req.Header.Set("Accept", "application/json")       // TODO: application/x-msgpack
	// TODO: req.Header.Set("User-Agent", "xxx") // UA is recommended to be set

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to call API. err=%s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, xerrors.Errorf("error in HTTP request and failed to parse response message. status=%s, err=%q", resp.Status, err)
		}
		return nil, xerrors.Errorf("error in HTTP request. status=%s, msg=%q", resp.Status, b)
	}

	var data *response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, xerrors.Errorf("failed to parse response JSON. err=%s", err)
	}
	return data, nil
}

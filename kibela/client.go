package kibela

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/xerrors"
)

const apiEndpoint string = "https://%s.kibe.la/api/v1"

type Kibela struct {
	token    string
	endpoint string
	client   *http.Client
}

func New(token, team string) (*Kibela, error) {
	return &Kibela{
		token:    token,
		endpoint: fmt.Sprintf(apiEndpoint, team),
		client:   &http.Client{},
	}, nil
}

type request struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables,omitempty"`
}

type response struct {
	Data json.RawMessage `json:"data"`
}

func (k *Kibela) AddLink(noteID, title, url string) (*Note, error) {
	note, err := k.GetNote(noteID)
	if err != nil {
		return nil, xerrors.New(fmt.Sprintf("Failed to get note of id=%s. err=%s", noteID, err))
	}
	updatedContent := note.Content + "\n" + fmt.Sprintf("- [%s](%s)", title, url)
	resp, err := k.do(&request{
		Query: updateNoteMutation(),
		Variables: struct {
			ID       string     `json:"id"`
			NewNote  *NoteInput `json:"newNote"`
			BaseNote *NoteInput `json:"baseNote"`
		}{
			ID: note.ID,
			NewNote: &NoteInput{
				Title:     note.Title,
				Content:   updatedContent,
				GroupIds:  note.GroupIds,
				Coediting: note.Coediting,
			},
			BaseNote: &NoteInput{
				Title:     note.Title,
				Content:   note.Content,
				GroupIds:  note.GroupIds,
				Coediting: note.Coediting,
			},
		},
	})
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

func (k *Kibela) do(reqBody *request) (*response, error) {
	body := bytes.Buffer{}
	if err := json.NewEncoder(&body).Encode(reqBody); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, k.endpoint, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.token))
	req.Header.Set("Content-Type", "application/json") // TODO: application/x-msgpack
	req.Header.Set("Accept", "application/json")       // TODO: application/x-msgpack
	// TODO: req.Header.Set("User-Agent", "xxx")

	fmt.Println(req.Body)

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error in HTTP request and failed to parse response message. status=%s, err=%q", resp.Status, err)
		}
		return nil, fmt.Errorf("error in HTTP request. status=%s, msg=%q", resp.Status, b)
	}

	var data *response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, xerrors.New("Failed to parse response JSON. err=" + err.Error())
	}
	return data, nil
}

/*
Queries
*/

func getNoteQuery(id string) string {
	return fmt.Sprintf(`{
  note(id: "%s") {
    id
    title
    content
  }
}
`, id)
}

/*
Mutations
*/

type NoteInput struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	GroupIds  []string `json:"groupIds"`
	Coediting bool     `json:"coediting"`
	// folderName string   `jsopn:"folderName,omitempty"`
}

func updateNoteMutation() string {
	return `mutation($id: ID!, $newNote: NoteInput!, $baseNote: NoteInput!) {
  updateNote(input: {
    id: $id
    newNote: $newNote
    baseNote: $baseNote
    draflt: false
  })
  {
    note {
      updatedAt
    }
  }
}`
}

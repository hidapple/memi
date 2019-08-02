package kibela

import "fmt"

type UpdatNoteContentInput struct {
	ID          string `json:"id"`
	BaseContent string `json:"baseContent"`
	NewContent  string `json:"newContent"`
	Touch       bool   `json:"touch"`
}

func updateNoteMutation(id string, baseContent, newContent string) string {
	return fmt.Sprintf(`mutation {
  updateNoteContent(input: {
    id: "%s",
    baseContent: "%s",
    newContent: "%s",
    touch: true })
  {
    note {
      id
      title
      content
      groups {
        id
        name
      }
    }
  }
}`, id, baseContent, newContent)
}

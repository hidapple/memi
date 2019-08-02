package kibela

import "fmt"

func getNoteQuery(id string) string {
	return fmt.Sprintf(`{
  note(id: "%s") {
    id
    title
    content
    coediting
    groups {
      id
      name
    }
  }
}
`, id)
}

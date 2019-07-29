package kibela

type Note struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Groups    []*Group `json:"groups"`
	CoEditing bool     `json:"coediting"`
}

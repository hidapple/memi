package kibela

type Note struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	GroupIds  []string `json:"groupIds"`
	Coediting bool     `json:"coediting"`
}

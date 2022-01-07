package models

type Posts struct {
	ID    int    `json:"id,omitempty"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (r Posts) TableName() string {
	return "posts"
}

package model

type Grade struct {
	BaseModel

	Class   string
	Name    string
	Score   int
	Subject string
}

func (t *Grade) TableName() string {
	return "grade"
}

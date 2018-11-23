package gquiz

type UI interface {
	Println(message string)
	GetInput() string
}

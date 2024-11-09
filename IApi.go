package main
type Todos struct{
	Todos []Todo
}
type Todo struct{
	Name string
	Description string
}
type IApi interface{
	GetTodos() (*Todos,error)
	SetAsDone()error
	PostTodos(Todos) bool
}

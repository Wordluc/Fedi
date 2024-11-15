package Api

type Todos struct {
	Todos []Todo
}
type Todo struct {
	Id          string
	Name        string
	Description string
}
type IApi interface {
	GetTodos() (*Todos, error)
	SetAsDone() error
	PostTodos(Todos) error
	Delete(Todo) error
}

func CreateClient(fileEnvName string) (IApi, error) {
	return CreateNotionClient(fileEnvName)
}

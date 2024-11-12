package Api

type Todos struct {
	Todos []Todo
}
type Todo struct {
	Name        string
	Description string
}
type IApi interface {
	GetTodos() (*Todos, error)
	SetAsDone() error
	PostTodos(Todos) error
}

func CreateClient(fileEnvName string) (IApi, error) {
	return CreateNotionClient(fileEnvName)
}

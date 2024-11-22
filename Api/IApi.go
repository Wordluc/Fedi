package Api

type Todos struct {
	Todos []Todo
}
type Todo struct {
	Id          string
	Name        string
	Description string
	Status      string
}
type IApi interface {
	GetTodos() (*Todos, error)
	SetAsDone(Todo) error
	PostTodos(Todos)(TodoPost, error)
	Delete(Todo) error
}

func CreateClient(fileEnvName string) (IApi, error) {
	return CreateNotionClient(fileEnvName)
}

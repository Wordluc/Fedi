package Api

import (
	"Fedi/Api/internal"
	"Fedi/Api/internal/impl"
)

type IApi interface {
	GetTodos() (*internal.Todos, error)
	SetAsDone() error
	PostTodos(internal.Todos) bool
}

func CreateClient(fileEnvName string) (IApi, error) {
	return impl.CreateNotionClient(fileEnvName)
}

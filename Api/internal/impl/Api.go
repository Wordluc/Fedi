package impl

import (
	"Fedi/Api/internal"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/joho/godotenv"
)
type NotionClient struct {
	Notion_key string
	Notion_database_id  string
}
func CreateNotionClient(fileEnvName string) (*NotionClient,error) {
	env,err:=godotenv.Read(fileEnvName)
	if err!=nil{
		return nil,err
	}
	return &NotionClient{
		Notion_key: env["NOTION_KEY"],
		Notion_database_id:  env["NOTION_DATABASE_ID"],
	},nil
	

}

func (c *NotionClient) getRequest(url string,method string,body io.Reader) (*http.Request, error) {
	var req, err = http.NewRequest(method, url,body)
	if err!=nil{
		return nil,err
	}
	req.Header.Add("Authorization", "Bearer "+c.Notion_key)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	return req, nil
}
func (c *NotionClient) GetTodos() (*internal.Todos,error) {
	var url= "https://api.notion.com/v1/databases/"+c.Notion_database_id+"/query"
	req,err:=c.getRequest(url,"POST",nil)
	if err!=nil{
		return nil,err
	}
	
	resp, err := http.DefaultClient.Do(req)
	if err!=nil{
		return nil,err
	}
	defer resp.Body.Close()
	if resp.StatusCode!=200{
		return nil,errors.New("Notion error "+fmt.Sprint(resp.StatusCode)+":"+resp.Status)
	}
	respBodyString,err := io.ReadAll(resp.Body) 
	if err!=nil{
		return nil,err
	}
	var respose *NotionResponse=&NotionResponse{}
	e:=json.Unmarshal(respBodyString,respose)
	if e!=nil{
		return nil,e
	}
	return mappingResponse(respose),nil
}

func (c *NotionClient) PostTodos(internal.Todos) bool {
	return true
}

func (c *NotionClient) SetAsDone() error {
	return nil
}

func mappingResponse(resp *NotionResponse) *internal.Todos{
	var todos internal.Todos=internal.Todos{
		Todos:make([]internal.Todo,0),
	}
	for _,page:=range resp.Results{
		todo:=internal.Todo{
			Name:page.Properties["Name"].Title[0].Text.Content,
			Description:page.Properties["Description"].Rich_Text[0].Text.Content,
		}
		todos.Todos=append(todos.Todos,todo)
	}
	return &todos
}

package main

import "fmt"

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/joho/godotenv"
)
type NotionClient struct {
	notion_key string
	notion_database_id  string
}
func CreateNotionClient(fileEnvName string) (IApi,error) {
	env,err:=godotenv.Read(fileEnvName)
	if err!=nil{
		return nil,err
	}
	return &NotionClient{
		notion_key: env["NOTION_KEY"],
		notion_database_id:  env["NOTION_DATABASE_ID"],
	},nil
	

}
func (c *NotionClient) GetTodos() (*Todos,error) {
	var url= "https://api.notion.com/v1/databases/"+c.notion_database_id+"/query"
	var req, err = http.NewRequest("POST", url,nil)
	if err!=nil{
		return nil,err
	}
	req.Header.Add("Authorization", "Bearer "+c.notion_key)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	
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
func (c *NotionClient) PostTodos(Todos) bool {
	return true
}
func (c *NotionClient) SetAsDone() error {
	return nil
}
func mappingResponse(resp *NotionResponse) *Todos{
	var todos Todos=Todos{
		Todos:make([]Todo,0),
	}
	for _,page:=range resp.Results{
		todo:=Todo{
			Name:page.Properties["Name"].Title[0].Text.Content,
			Description:page.Properties["Description"].Rich_Text[0].Text.Content,
		}
		todos.Todos=append(todos.Todos,todo)
	}
	return &todos
}

package Api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/joho/godotenv"
)

type NotionClient struct {
	Notion_key         string
	Notion_database_id string
}

func CreateNotionClient(fileEnvName string) (*NotionClient, error) {
	env, err := godotenv.Read(fileEnvName)
	if err != nil {
		return nil, err
	}
	return &NotionClient{
		Notion_key:         env["NOTION_KEY"],
		Notion_database_id: env["NOTION_DATABASE_ID"],
	}, nil

}

func (c *NotionClient) getRequest(url string, method string, body io.Reader) (*http.Request, error) {
	var req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.Notion_key)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	return req, nil
}
func (c *NotionClient) GetTodos() (*Todos, error) {
	var url = "https://api.notion.com/v1/databases/" + c.Notion_database_id + "/query"
	req, err := c.getRequest(url, "POST", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("Notion error " + fmt.Sprint(resp.StatusCode) + ":" + resp.Status)
	}
	respBodyString, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var respose *NotionResponse = &NotionResponse{}
	e := json.Unmarshal(respBodyString, respose)
	if e != nil {
		return nil, e
	}
	return mappingResponse(respose), nil
}

func (c *NotionClient) PostTodos(todos Todos) error {
	url := "https://api.notion.com/v1/pages"

	property := map[string]interface{}{
		"Description": map[string]interface{}{
			"rich_text": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": todos.Todos[0].Description,
					},
				},
			},
		},
		"Name": map[string]interface{}{
			"title": []map[string]interface{}{
				{
					"text": map[string]string{
						"content": todos.Todos[0].Name,
					},
				},
			},
		},
	}

	requestBody := map[string]interface{}{
		"parent": map[string]string{
			"database_id": c.Notion_database_id,
		},
		"properties": property,
	}
	marshal, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	request, err := c.getRequest(url, "POST", bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New("Notion error " + fmt.Sprint(resp.StatusCode) + ":" + string(response))
	}
	return nil
}

func (c *NotionClient) SetAsDone() error {
	return nil
}

func mappingResponse(resp *NotionResponse) *Todos {
	var todos Todos = Todos{
		Todos: make([]Todo, 0),
	}
	for _, page := range resp.Results {
		todo := Todo{
			Name:        page.Properties["Name"].Title[0].Text.Content,
			Description: page.Properties["Description"].Rich_Text[0].Text.Content,
		}
		todos.Todos = append(todos.Todos, todo)
	}
	return &todos
}

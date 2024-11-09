package main

import "time"

type NotionResponse struct {
	Object         string                 
	Results        []Page                 
	NextCursor     string                 
	HasMore        bool                   
	Type           string                 
	PageOrDatabase map[string]interface{} 
	RequestID      string                 
}

type Page struct {
	Object         string              
	ID             string              
	CreatedTime    time.Time           
	LastEditedTime time.Time           
	CreatedBy      User                
	LastEditedBy   User                
	Cover          interface{}         
	Icon           interface{}         
	Parent         Parent              
	Archived       bool                
	InTrash        bool                
	Properties     map[string]Property 
	URL            string              
	PublicURL      string              
}

type User struct {
	Object string 
	ID     string 
}

type Parent struct {
	Type       string 
	DatabaseID string 
}

type Property struct {
	ID       string     
	Type     string     
	Rich_Text []RichText 
	Select   *Select    
	Title    []Title    
}

type RichText struct {
	Type        string      
	Text        Text        
	Annotations Annotations 
	PlainText   string      
	Href        string      
}

type Text struct {
	Content string 
	Link    string 
}

type Annotations struct {
	Bold          bool   
	Italic        bool   
	Strikethrough bool   
	Underline     bool   
	Code          bool   
	Color         string 
}

type Select struct {
	// Add fields as needed
}

type Title struct {
	Type        string      
	Text        Text        
	Annotations Annotations 
	PlainText   string      
	Href        string      
}

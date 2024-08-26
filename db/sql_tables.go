/*
These are structures that can be easily be swapped between JSON, to return to the client side, and SQL, to store them in a databasae.
*/
package db

import (
	"time"
)

// Defines a user in the system
type User struct {
	Id         int    `json:"id" sql_name:"id" sql_props:"INTEGER PRIMARY KEY AUTOINCREMENT"` // Database ID
	Username   string `json:"username" sql_name:"username" sql_props:"TEXT UNIQUE NOT NULL"`
	Email      string `json:"email" sql_name:"email" sql_props:"TEXT NOT NULL"`
	IsCustomer bool   `json:"is_client" sql_name:"is_client" sql_props:"BOOLEAN NOT NULL"` // Is this a customer or a support rep?
	// Password string `json:"password" sql_name:"password" sql_props:"TEXT NOT NULL"`
	// Phone    string `json:"phone" sql_name:"phone" sql_props:"TEXT NOT NULL"`
}

// type IngestFileProfile struct {
// 	Id               int       `json:"id" sql_name:"id" sql_props:"INTEGER PRIMARY KEY AUTOINCREMENT"`
// 	Name             string    `json:"name" sql_name:"name" sql_props:"TEXT NOT NULL"`
// 	Description      string    `json:"description" sql_name:"description" sql_props:"TEXT NOT NULL"`
// 	FirstCreated     time.Time `json:"first_created" sql_name:"first_created" sql_props:"TIMESTAMP NOT NULL"`
// 	LastUpdated      time.Time `json:"last_updated" sql_name:"last_updated" sql_props:"TIMESTAMP NOT NULL"`
// 	FilePath         string    `json:"file_path" sql_name:"file_path" sql_props:"TEXT NOT NULL"`
// 	AiOcrTextSummary string    `json:"ai_ocr_text_summary" sql_name:"ai_ocr_text_summary" sql_props:"TEXT NOT NULL"`
// }

// Defines a support ticket interaction between a customer and a support rep.
// This includes logs, as well as the resolution status of the ticket.
type SupportTicket struct {
	Id             int       `json:"id" sql_name:"id" sql_props:"INTEGER PRIMARY KEY AUTOINCREMENT"`
	CustomerId     string    `json:"name" sql_name:"name" sql_props:"TEXT NOT NULL"`
	SupportStaffId string    `json:"description" sql_name:"description" sql_props:"TEXT NOT NULL"`
	MainLog        string    `json:"main_log" sql_name:"main_log" sql_props:"TEXT NOT NULL"`
	Reason         string    `json:"reason" sql_name:"reason" sql_props:"TEXT NOT NULL"`
	IsResolved     bool      `json:"is_resolved" sql_name:"is_resolved" sql_props:"BOOLEAN NOT NULL"`
	FirstCreated   time.Time `json:"first_created" sql_name:"first_created" sql_props:"TIMESTAMP NOT NULL"`
	LastUpdated    time.Time `json:"last_updated" sql_name:"last_updated" sql_props:"TIMESTAMP NOT NULL"`
	// FilePath         string    `json:"file_path" sql_name:"file_path" sql_props:"TEXT NOT NULL"`
	// AiOcrTextSummary string    `json:"ai_ocr_text_summary" sql_name:"ai_ocr_text_summary" sql_props:"TEXT NOT NULL"`
}

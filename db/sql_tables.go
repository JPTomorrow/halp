/*
These are structures that can be easily be swapped between JSON, to return to the client side, and SQL, to store them in a databasae.
*/
package db

import (
	"time"
)

// Defines a user in the system
type Customer struct {
	Id       int    `json:"id" sql_name:"id" sql_props:"INTEGER PRIMARY KEY AUTOINCREMENT"` // Database ID
	Username string `json:"username" sql_name:"username" sql_props:"TEXT UNIQUE NOT NULL"`
	Email    string `json:"email" sql_name:"email" sql_props:"TEXT NOT NULL"`
	Password string `json:"password" sql_name:"password" sql_props:"TEXT NOT NULL"`
}

type SalesRep struct {
	Id       int    `json:"id" sql_name:"id" sql_props:"INTEGER PRIMARY KEY AUTOINCREMENT"` // Database ID
	Username string `json:"username" sql_name:"username" sql_props:"TEXT UNIQUE NOT NULL"`
	Email    string `json:"email" sql_name:"email" sql_props:"TEXT NOT NULL"`
	Password string `json:"password" sql_name:"password" sql_props:"TEXT NOT NULL"`
}

// Defines a support ticket interaction between a customer and a support rep.
// This includes logs, as well as the resolution status of the ticket.
type SupportTicket struct {
	Id             int       `json:"id" sql_name:"id" sql_props:"INTEGER PRIMARY KEY AUTOINCREMENT"`
	MainLog        string    `json:"main_log" sql_name:"main_log" sql_props:"TEXT NOT NULL"`
	Reason         string    `json:"reason" sql_name:"reason" sql_props:"TEXT NOT NULL"`
	IsResolved     bool      `json:"is_resolved" sql_name:"is_resolved" sql_props:"BOOLEAN NOT NULL"`
	FirstCreated   time.Time `json:"first_created" sql_name:"first_created" sql_props:"TIMESTAMP NOT NULL"`
	LastUpdated    time.Time `json:"last_updated" sql_name:"last_updated" sql_props:"TIMESTAMP NOT NULL"`
	CustomerId     int       `json:"customer_id" sql_name:"customer_id" sql_props:"INTEGER NOT NULL" sql_fk:"customer.id"`
	SupportStaffId int       `json:"support_staff_id" sql_name:"support_staff_id" sql_props:"INTEGER NOT NULL" sql_fk:"sales_rep.id"`
}

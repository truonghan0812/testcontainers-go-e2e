package main

// type User struct {
// 	ID   int64  `json:"id"`
// 	Name string `json:"name"`
// }
type User struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type TicketPost struct {
	UserID int64  `json:"user_id"`
	Movie  string `json:"movie"`
}
type Ticket struct {
	ID    string `json:"id"`
	Movie string `json:"movie"`
	User  User   `json:"user"`
}

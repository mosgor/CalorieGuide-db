package client

type Client struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"user_name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Goal     int    `json:"goal_id,omitempty"`
	Diet     int    `json:"diet_id,omitempty"`
}

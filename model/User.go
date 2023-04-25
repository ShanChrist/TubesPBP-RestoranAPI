package model

type User struct {
	UserID          int    `json:"user_id"`
	UserName        string `json:"username"`
	UserFirstName   string `json:"first_name"`
	UserLastName    string `json:"last_name"`
	UserPhoneNumber string `json:"phone_number"`
	UserEmail       string `json:"useremail"`
	UserPassword    string `json:"userpassword"`
	UserAddress     string `json:"address"`
	UserType        int    `json:"usertype"`
}

type LihatProfile struct {
	UserName        string `json:"username"`
	UserFirstName   string `json:"first_name"`
	UserLastName    string `json:"last_name"`
	UserPhoneNumber string `json:"phone_number"`
	UserEmail       string `json:"useremail"`
	UserAddress     string `json:"address"`
}

type LihatProfileResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    LihatProfile `json:"data"`
}

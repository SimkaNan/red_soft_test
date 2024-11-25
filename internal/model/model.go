package model

type User struct {
	Id          int      `json:"id" db:"id"`
	Name        string   `json:"name" db:"name" binding:"required"`
	Surname     string   `json:"surname" db:"surname" binding:"required"`
	MiddleName  string   `json:"middle_name" db:"middle_name" binding:"required"`
	Age         int      `json:"age" db:"age" binding:"required"`
	Nationality string   `json:"nationality" db:"nation" binding:"required"`
	Gender      string   `json:"gender" db:"gender" binding:"required"`
	Emails      []string `json:"emails"`
}

type Email struct {
	Email string `json:"email" db:"email" binding:"required"`
}

type UserUpdate struct {
	Name        string `json:"name" db:"name"`
	Surname     string `json:"surname" db:"surname"`
	MiddleName  string `json:"middle_name" db:"middle_name"`
	Age         int    `json:"age" db:"age"`
	Nationality string `json:"nation" db:"nationality"`
	Gender      string `json:"gender" db:"gender"`
}

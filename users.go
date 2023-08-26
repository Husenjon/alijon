package inkassback

import "database/sql"

type User struct {
	Id          int            `db:"id"`
	Ism         string         `json:"ism" db:"ism" validate:"required"`
	Familya     string         `json:"familya" db:"familya"`
	OtasiniIsmi string         `json:"otasini_ismi" db:"otasini_ismi"`
	Phone       string         `json:"phone" db:"phone"`
	Username    string         `json:"username" db:"username" validate:"required"`
	Password    string         `json:"password" db:"password" validate:"required"`
	BranchId    int            `json:"branch_id" db:"branch_id" validate:"required" `
	Token       sql.NullString `json:"token" db:"token"`
	Image       string         `json:"image" db:"image"`
	IsActive    bool           `json:"is_active" db:"is_active"`
	CreatedTime string         `json:"created_time" db:"created_time"`
	UpdatedTime sql.NullString `json:"updated_time" db:"updated_time"`
}

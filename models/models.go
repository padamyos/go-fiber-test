package models

import "gorm.io/gorm"

type Person struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type User struct {
	Email string `json:"email,omitempty" validate:"required,email,min=3,max=32"`
	// กำหนดให้ Name รับแต่ตัวอักษรภาษาอังกฤษ ตัวเลข และสัญลักษณ์ _,- เท่านั้น
	// regex ที่ใช้ ^[a-zA-Z0-9_-]+$ คือ ต้องเป็นตัวอักษรภาษาอังกฤษ ตัวเลข และสัญลักษณ์ _,- เท่านั้น
	Name        string `json:"name,omitempty" validate:"required,regexp=^[a-zA-Z0-9_-]+$,min=3,max=32"`
	Password    string `json:"password,omitempty" validate:"required,min=6,max=20"`
	LineId      string `json:"line_id,omitempty" validate:"max=32"`
	PhoneNumber string `json:"phone_number,omitempty" validate:"required,min=10,max=10"`
	// ประเภทธุรกิจ dropdown
	BusinessType string `json:"business_type,omitempty" validate:"required,min=3,max=32"`
	WebsiteName  string `json:"website_name,omitempty" validate:"required,regexp=^[a-zA-Z0-9-]+$,min=2,max=30"`
}

type Dogs struct {
	gorm.Model
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
}

type Company struct {
	gorm.Model
	CompanyName string `json:"companyname"`
	CompanyID   int    `json:"company_id"`
	Assdress    string `json:"address"`
}

// profile
type Profile struct {
	gorm.Model
	Employees int `json:"employees"`
	Name	  string `json:"name"`
	LastNames string `json:"lastnames"`
	BirthDate string `json:"birthdate"`
	Age 	 int `json:"age"`
	Email	 string `json:"email"`
	Tel		 string `json:"tel"`
}
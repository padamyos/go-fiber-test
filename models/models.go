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

// DOGS
type Dogs struct {
	gorm.Model
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
}

type DogsRes struct {
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
	Type  string `json:"type"`
}

type ResultData struct {
	Data       []DogsRes `json:"data"`
	Name       string    `json:"name"`
	Count      int       `json:"count"`
	SumRed     int       `json:"sum_red"`
	SumGreen   int       `json:"sum_green"`
	SumPink    int       `json:"sum_pink"`
	SumNoColor int       `json:"sum_no_color"`
}

// COMPANY
type Company struct {
	gorm.Model
	CompanyName string `json:"company_name"`
	CompanyID   int    `json:"company_id"`
	AssDress    string `json:"address"`
}

// PROFILE
type Profile struct {
	gorm.Model
	Employees string `json:"employees"`
	Name      string `json:"name"`
	LastNames string `json:"last_names"`
	BirthDate string `json:"birthdate"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	Tel       string `json:"tel"`
}

type ProfileRes struct {
	Employees string `json:"employees"`
	Name      string `json:"name"`
	LastNames string `json:"last_names"`
	Age       int    `json:"age"`
	Type      string `json:"type"`
}

type ResultProfile struct {
	Count int `json:"count"`

	Data []ProfileRes `json:"data"`

	GenZ int `json:"gen_z"`

	GenY int `json:"gen_y"`

	GenX int `json:"gen_x"`

	BabyBoomer int `json:"baby_boomer"`

	GiGeneration int `json:"gi_generation"`
}

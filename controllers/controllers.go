package controllers

import (
	"fmt"
	"go-fiber-test/database"
	m "go-fiber-test/models"

	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func Hellotest(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func BodyParser(c *fiber.Ctx) error {
	p := new(m.Person)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	log.Println(p.Name) // john
	log.Println(p.Pass) // doe
	str := p.Name + p.Pass
	return c.JSON(str)
}

func Params(c *fiber.Ctx) error {

	str := "hello ==> " + c.Params("name")
	return c.JSON(str)
}

func Search(c *fiber.Ctx) error {
	a := c.Query("search")
	str := "my search is  " + a
	return c.JSON(str)
}

func Validate(c *fiber.Ctx) error {
	//Connect to database

	user := new(m.User)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	validate := validator.New()
	errors := validate.Struct(user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.Error())
	}
	return c.JSON(user)
}

// 5.1สร้างapi รับค่าตัวเลข ผ่านpath แล้วreturnเป็นค่าfactorialของตัวเลขนั้น
func Factorial(c *fiber.Ctx) error {
	numStr := c.Params("num")
	num, err := strconv.Atoi(numStr) //แปลงเป็น int
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid number",
		})
	}
	factorial := 1
	for i := 1; i <= num; i++ {
		factorial *= i
	}

	return c.JSON(fmt.Sprintf("%d! = %d", num, factorial))
}

// 5.2 สร้างapiขึ้นต้นด้วย api/v3/ (<--ใช้วิธีแบบจัดgroup api)ตามด้วยชื่อเล่นตัวเอง  โดยapiนี้มีการรับ QueryParam ที่ชื่อkeyว่า tax_id นำค่าที่keyเข้าไป(keyได้ทั้งตัวเลขตัวอักษร)แปลงเป็น ascii
func TaxId(c *fiber.Ctx) error {
	taxId := c.Query("tax_id")
	ascii := ""
	for _, char := range taxId {
		ascii += strconv.Itoa(int(char))
		ascii += " "
	}
	return c.JSON(fiber.Map{
		"ascii": ascii,
	})
}

// 5.3 เปลี่ยนชื่อการเรียกใช้ controller ให้สั้นลง controller.TestParams → c.TestParams

// 6.api method POST สมัครสมาชิก ดักฟิลข้อมูลให้ถูกต้อง localhost:3000/api/v1/register และถ้าใส่ข้อมูลไม่ถูกต้องให้โชว์ใส่ข้อมูลผิดพลาด hint : regexp.MatchString
func Register(c *fiber.Ctx) error {
	user := new(m.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// ตรวจสอบข้อมูล
	validate := validator.New()
	validate.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		regex := fl.Param()                            //ดึงค่า param จาก tag
		value := fl.Field().String()                   //ดึงค่าจาก field แปลงเป็น string
		matched, _ := regexp.MatchString(regex, value) //ตรวจสอบความถูกต้องของข้อมูล
		return matched
	})
	// ตรวจความถูกต้องของข้อมูล
	errors := validate.Struct(user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
			"errors":  errors.Error(),
		})
	}
	return c.JSON(user)
}

// CRUD dogs
func GetDogs(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs
	db.Find(&dogs) //delelete = null
	return c.Status(200).JSON(dogs)
}

func GetDog(c *fiber.Ctx) error {
	db := database.DBConn
	search := strings.TrimSpace(c.Query("search"))
	var dog []m.Dogs

	result := db.Find(&dog, "dog_id = ?", search)

	// returns found records count, equals `len(users)
	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	return c.Status(200).JSON(&dog)
}

func AddDog(c *fiber.Ctx) error {
	//twst3
	db := database.DBConn
	var dog m.Dogs

	if err := c.BodyParser(&dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&dog)
	return c.Status(201).JSON(dog)
}

func UpdateDog(c *fiber.Ctx) error {
	db := database.DBConn
	var dog m.Dogs
	id := c.Params("id")

	if err := c.BodyParser(&dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&dog)
	return c.Status(200).JSON(dog)
}

func RemoveDog(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var dog m.Dogs

	result := db.Delete(&dog, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}



// 7.0.2 สร้าง api GET ใน group dogs โชว์ข้อมูลที่ถูกลบไปแล้ว ตารางdogs
func GetDogsDeleted(c *fiber.Ctx) error {
    db := database.DBConn
    var dogs []m.Dogs

    // ดึงเฉพาะข้อมูลที่มี DeletedAt != null
    db.Unscoped().Where("deleted_at IS NOT NULL").Find(&dogs)

    return c.Status(200).JSON(dogs)
}

// 7.1 สร้างapi GETใหม่ แสดงข้อมูลตารางdogโดย where หา dog_id > 50 แต่น้อยกว่า 100  (gorm)
func GetDogsWhere(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	// ดึงข้อมูลที่ dog_id > 50 แต่น้อยกว่า 100
	db.Where("dog_id > ? AND dog_id < ?", 50, 100).Find(&dogs)
	return c.Status(200).JSON(dogs)
}

// 7.2 สร้างapi GETใหม่ แสดงข้อมูลตารางdogโดย where หา dog_id (gorm) พร้อมผลรวมแต่ละสี
func GetDogsJson(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	db.Find(&dogs) //10ตัว

	type DogsRes struct {
		Name  string `json:"name"`
		DogID int    `json:"dog_id"`
		Type  string `json:"type"`
	
	}

	var dataResults []DogsRes
	var sum_red, sum_green, sum_pink, sum_no_color int
	for _, v := range dogs { //1 inet 112 //2 inet1 113
		typeStr := ""
		
		if v.DogID >= 10 && v.DogID <= 50 {
			typeStr = "red"
			sum_red++
		} else if v.DogID >= 100 && v.DogID <= 150 {
			typeStr = "green"
			sum_green++
		} else if v.DogID >= 200 && v.DogID <= 250 {
			typeStr = "pink"
			sum_pink ++
		} else {
			typeStr = "no color"
			sum_no_color++
		}

		d := DogsRes{
			Name:  v.Name,  //inet
			DogID: v.DogID, //112
			Type:  typeStr, //no color
		
		}
		dataResults = append(dataResults, d) 
		
		// sumAmount += v.Amount
	}

	type ResultData struct {
		Data  []DogsRes `json:"data"`
		Name  string    `json:"name"`
		Count int       `json:"count"`
		SumRed int `json:"sum_red"`
		SumGreen int `json:"sum_green"`
		SumPink int `json:"sum_pink"`
		SumNoColor int `json:"sum_no_color"`
	}
	
	r := ResultData{
		Count: len(dogs), //หาผลรวม,
		Data:  dataResults,
		SumRed:     sum_red,
        SumGreen:   sum_green,
        SumPink:    sum_pink,
        SumNoColor: sum_no_color,
		
	}
	return c.Status(200).JSON(r)
}

// CRUD company
func GetCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var company []m.Company

	db.Find(&company) //delelete = null
	return c.Status(200).JSON(company)
}

// add company
func AddCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var company m.Company

	if err := c.BodyParser(&company); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&company)
	return c.Status(201).JSON(company)
}

// update company
func UpdateCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var company m.Company
	id := c.Params("id")

	if err := c.BodyParser(&company); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&company)
	return c.Status(200).JSON(company)
}

// delete company
func RemoveCompany(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var company m.Company

	result := db.Delete(&company, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

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
	db.Where("dog_id > ? ",100).Find(&dogs)
	return c.Status(200).JSON(dogs)
}

// 7.2 สร้างapi GETใหม่ แสดงข้อมูลตารางdogโดย where หา dog_id (gorm) พร้อมผลรวมแต่ละสี
func GetDogsJson(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	db.Find(&dogs) //10ตัว

	var dataResults []m.DogsRes
	var sumRed, sumGreen, sumPink, sumNoColor int
	for _, v := range dogs { 
		typeStr := ""
		
		if v.DogID >= 10 && v.DogID <= 50 {
			typeStr = "red"
			sumRed++
		} else if v.DogID >= 100 && v.DogID <= 150 {
			typeStr = "green"
			sumGreen++
		} else if v.DogID >= 200 && v.DogID <= 250 {
			typeStr = "pink"
			sumPink ++
		} else {
			typeStr = "no color"
			sumNoColor++
		}

		d := m.DogsRes{
			Name:  v.Name,  //inet
			DogID: v.DogID, //112
			Type:  typeStr, //no color
		
		}
		dataResults = append(dataResults, d) 
		
		// sumAmount += v.Amount
	}
	r := m.ResultData{
		Count: len(dogs), //หาผลรวม,
		Name: "Dogs",
		Data:  dataResults,
		SumRed:     sumRed,
        SumGreen:   sumGreen,
        SumPink:    sumPink,
        SumNoColor: sumNoColor,
		
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

// CRUD profile
// Get profile
func GetProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var profile []m.Profile

	db.Find(&profile) //delelete = null
	return c.Status(200).JSON(profile)
}
// add profile
func AddProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var profile m.Profile

	if err := c.BodyParser(&profile); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&profile)
	return c.Status(201).JSON(profile)
}
// update profile
func UpdateProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var profile m.Profile
	id := c.Params("id")

	if err := c.BodyParser(&profile); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&profile)
	return c.Status(200).JSON(profile)
}
// delete profile
func RemoveProfile(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var profile m.Profile

	result := db.Delete(&profile, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

// API GET ข้อมูลผู้ใช้และโชว์จำนวนของแต่ละประเภทกลุ่มอายุ 
func GetProfileGroup(c *fiber.Ctx) error {
	db := database.DBConn
	var profile []m.Profile

	db.Find(&profile) //delelete = null
	var dataResults []m.ProfileRes
	var GenZ , GenY , GenX , Baby_Boomer , GI_Generation int

	for _, v := range profile { 
		typeStr := ""
	
		
		if v.Age <= 24 {
			typeStr = "GenZ"
			GenZ++
		} else if v.Age >= 24 && v.Age <= 41 {
			typeStr = "GenY"
			GenY++
		} else if v.Age >= 42 && v.Age <= 56 {
			typeStr = "GenX"
			GenX++
		} else if v.Age >= 57 && v.Age <= 75 {
			typeStr = "Baby_Boomer"
			Baby_Boomer++
		} else if v.Age > 75 {
			typeStr = "GI_Generation"
			GI_Generation++
		} else {
			typeStr = "no age"
		}

		d := m.ProfileRes{
			Employees:  v.Employees,
			Name:  v.Name,
			LastNames:  v.LastNames,
			Age:  v.Age,
			Type:  typeStr,
		}
		dataResults = append(dataResults, d)
}

	r := m.ResultProfile{
		Count: len(profile), //หาผลรวม,
		Data:  dataResults,
		GenZ:     GenZ,
		GenY:   GenY,
		GenX:    GenX,
		BabyBoomer: Baby_Boomer,
		GiGeneration: GI_Generation,
	}
	return c.Status(200).JSON(r)
}

// สร้าง API search ข้อมูลโปรไฟล์ผู้ใช้ โดยที่สามารถserchได้3ตัวคือ employee_id, name ,lastname ภายในคีย์searchตัวเดียว  xxx/search
// func GetProfileSearch(c *fiber.Ctx) error {
// 	db := database.DBConn

// 	search := strings.TrimSpace(c.Query("search"))
// 	var profile []m.Profile
	
// 	result := db.Where("employees = ? OR name = ? OR last_names = ?", search, search, search).Find(&profile)

// 	return c.Status(200).JSON(result)
// }

func SearchProfile(c *fiber.Ctx) error {
    db := database.DBConn
    // search := c.Query("search")
	search := strings.TrimSpace(c.Query("search"))

    var profiles []m.Profile
	db.Where("employees LIKE ? OR name LIKE ? OR last_names LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").Find(&profiles)


    return c.Status(200).JSON(profiles)
}
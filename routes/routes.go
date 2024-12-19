package routes

import (
	c "go-fiber-test/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func Routes(app *fiber.App) {

	v1 := app.Group("/api/v1")
	v3 := app.Group("/api/v3/ta")

	// Provide a minimal config
	v1.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"gofiber": "21022566",
		},
	}))

	v1.Get("/", c.Hellotest)

	v1.Post("/", c.BodyParser)

	v1.Get("/user/:name", c.Params)

	v1.Post("/inet", c.Search)

	v1.Post("/valid", c.Validate)

	// ข้อที่ 5.1
	v1.Get("/fact/:num", c.Factorial)

	// ข้อที่ 5.2
	v3.Post("/taxId", c.TaxId)

	// ข้อที่ 6
	v1.Post("/register", c.Register)

	//CRUD dogs
	dog := v1.Group("/dog")
	dog.Get("", c.GetDogs)
	dog.Get("/filter", c.GetDog)
	dog.Get("/json", c.GetDogsJson)
	dog.Post("/", c.AddDog)
	dog.Put("/:id", c.UpdateDog)
	dog.Delete("/:id", c.RemoveDog)
	// -ข้อที่  7.0.2
	// dog.Get("/json/delete", c.GetDogsDeleted)
	dog.Get("/json/delete", c.GetDogsDeleted)
	// ข้อที่ 7.1
	dog.Get("/where", c.GetDogsWhere)

	//CRUD company
	company := v1.Group("/company")
	company.Get("", c.GetCompany)
	company.Post("/", c.AddCompany)
	company.Put("/:id", c.UpdateCompany)
	company.Delete("/:id", c.RemoveCompany)
}

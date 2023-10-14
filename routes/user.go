package routes

import (
	"errors"
	"fiber-api/database"
	"fiber-api/models"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// User serializer
type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func CreateResponseUser(userModel models.User) User {
	return User{ID: userModel.ID, FirstName: userModel.FirstName, LastName: userModel.LastName}
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	database.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)

	return c.Status(http.StatusCreated).JSON(responseUser)
}

func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}

	database.Database.Db.Find(&users)
	responseUsers := []User{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(http.StatusOK).JSON(responseUsers)
}

func findUser(id int, user *models.User) error {
	database.Database.Db.Find(&user, "id= ?", id)
	if user.ID == 0 {
		return errors.New("user does not exist")
	}
	return nil
}

func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	// when we pass the point, it will create two-way binding
	var user models.User
	if err := findUser(id, &user); err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}

	responseUser := CreateResponseUser(user)
	return c.Status(http.StatusOK).JSON(responseUser)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	// when we pass the point, it will create two-way binding
	var user models.User
	if err := findUser(id, &user); err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}

	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	var updateUser UpdateUser
	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	user.FirstName = updateUser.FirstName
	user.LastName = updateUser.LastName

	database.Database.Db.Save(&user)
	responseUser := CreateResponseUser(user)
	return c.Status(fiber.StatusAccepted).JSON(responseUser)
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	// when we pass the point, it will create two-way binding
	var user models.User
	if err := findUser(id, &user); err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&user).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}

	return c.Status(http.StatusAccepted).SendString("Successfully deleted user")
}

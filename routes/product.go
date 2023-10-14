package routes

import (
	"errors"
	"fiber-api/database"
	"fiber-api/models"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type Product struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

func CreateResponseProduct(productModel models.Product) Product {
	return Product{ID: productModel.ID, Name: productModel.Name, SerialNumber: productModel.SerialNumber}
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	database.Database.Db.Create(&product)
	responseProduct := CreateResponseProduct(product)
	return c.Status(http.StatusCreated).JSON(responseProduct)
}

func GetProducts(c *fiber.Ctx) error {
	products := []models.Product{}
	database.Database.Db.Find(&products)
	responseProducts := []Product{}
	for _, product := range products {
		responseProduct := CreateResponseProduct(product)
		responseProducts = append(responseProducts, responseProduct)
	}
	return c.Status(http.StatusOK).JSON(responseProducts)
}

func findProduct(id int, product *models.Product) error {
	database.Database.Db.Find(&product, "id= ?", id)
	if product.ID == 0 {
		return errors.New("product does not exist")
	}
	return nil
}

func GetProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	// when we pass the point, it will create two-way binding
	var product models.Product
	if err := findProduct(id, &product); err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}

	responseProduct := CreateResponseProduct(product)
	return c.Status(http.StatusOK).JSON(responseProduct)
}

func UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	// when we pass the point, it will create two-way binding
	var product models.Product
	if err := findProduct(id, &product); err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}

	type UpdateProduct struct {
		Name         string `json:"name"`
		SerialNumber string `json:"serial_number"`
	}

	var updateProduct UpdateProduct
	if err := c.BodyParser(&updateProduct); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	product.Name = updateProduct.Name
	product.SerialNumber = updateProduct.SerialNumber

	database.Database.Db.Save(&product)
	responseProduct := CreateResponseProduct(product)
	return c.Status(fiber.StatusAccepted).JSON(responseProduct)
}

func DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	// when we pass the point, it will create two-way binding
	var product models.Product
	if err := findProduct(id, &product); err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&product).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(err.Error())
	}

	return c.Status(http.StatusAccepted).SendString("Successfully deleted product")
}

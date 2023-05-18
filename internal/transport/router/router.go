package router

import (
	"encoding/json"
	"wb-first-lvl/internal/database/queries"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

const tampPath = "internal/transport/tamplate"

func Router(repo *queries.OrderRepo) {
	webApp := fiber.New()

	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.Render("internal/transport/tamplate/index.html", fiber.Map{})
	})

	webApp.Post("/order", func(c *fiber.Ctx) error {
		orderUID := c.FormValue("orderUID")

		ord, err := repo.GetExistingOrder(orderUID)
		if err != nil {
			return c.Render(tampPath+"/orderNotFound.html", fiber.Map{
				"uid": orderUID,
			})
		}

		jsonBytes, err := json.Marshal(&ord)
		if err != nil {
			logrus.Error(err)
		}

		return c.Render(tampPath+"/order.html", fiber.Map{
			"uid":   "orderUID",
			"order": string(jsonBytes),
		})
	})

	webApp.Get("/orderNotFound", func(c *fiber.Ctx) error {
		return c.Render(tampPath+"/orderNotFound.html", fiber.Map{})
	})

	port := "3080"
	logrus.WithFields(logrus.Fields{
		"port": port,
	}).Info("Starting a web-server on port")

	logrus.Fatal(webApp.Listen(":" + port))
}

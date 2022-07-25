package api

import (
	"fmt"
	"tabungan-api/app"
	"tabungan-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TabunganRESTAPI struct {
	server *fiber.App
	host   string
	port   int
	app    app.TabunganAppInterface
	log    *logrus.Logger
}

func (t *TabunganRESTAPI) registrasiNasabah(c *fiber.Ctx) (err error) {
	var request models.RequestRegistrasiNasabah
	response := make(map[string]interface{})
	err = c.BodyParser(&request)
	if err != nil {
		t.log.WithField("error", err.Error()).Error("parse request body to JSON error")
		response["remark"] = "Failed to parse request body"
		c.Status(400)
		return c.JSON(response)
	}
	rekening, err := t.app.RegistrasiNasabah(request)
	if err != nil {
		response["remark"] = err.Error()
		c.Status(400)
		return c.JSON(response)
	}
	response["data"] = rekening
	return c.JSON(response)
}

func (t *TabunganRESTAPI) Start() {
	addr := fmt.Sprintf("%s:%d", t.host, t.port)
	t.server.Listen(addr)
}

func NewRESTAPI(host string, port int, app app.TabunganAppInterface, logger *logrus.Logger) *TabunganRESTAPI {
	server := fiber.New()
	api := &TabunganRESTAPI{
		server: server,
		host:   host,
		port:   port,
		app:    app,
		log:    logger,
	}
	api.server.Post("/registrasi", api.registrasiNasabah)
	return api
}

package api

import (
	"fmt"
	"net/http"
	"strconv"
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
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	rekening, err := t.app.RegistrasiNasabah(request)
	if err != nil {
		response["remark"] = err.Error()
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	response["data"] = rekening
	return c.JSON(response)
}

func (t *TabunganRESTAPI) uploadFile(c *fiber.Ctx) (err error) {
	response := make(map[string]interface{})
	nik := c.Get("Authorization", "")
	if nik == "" {
		err = fmt.Errorf("missing NIK in authorization header")
		t.log.Warn(err.Error())
		response["remark"] = err.Error()
		c.Status(http.StatusUnauthorized)
		return c.JSON(response)
	}
	photo, err := c.FormFile("photo")
	if err != nil {
		t.log.WithField("error", err.Error()).Error("parse photo in multiform error")
		response["remark"] = "Failed to read photo file in multiform"
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	file, err := photo.Open()
	if err != nil {
		t.log.WithField("error", err.Error()).Error("parse photo in multiform error")
		response["remark"] = "Failed to read photo file in multiform"
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	err = t.app.SavePhoto(file, photo.Filename, nik)
	if err != nil {
		t.log.WithField("error", err.Error()).Error("save photo error")
		response["remark"] = "Failed to save photo"
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}

	doc, err := c.FormFile("doc")
	if err != nil {
		t.log.WithField("error", err.Error()).Error("parse doc in multiform error")
		response["remark"] = "Failed to read doc file in multiform"
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	file, err = doc.Open()
	if err != nil {
		t.log.WithField("error", err.Error()).Error("parse doc in multiform error")
		response["remark"] = "Failed to read doc file in multiform"
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	err = t.app.SaveDoc(file, doc.Filename, nik)
	if err != nil {
		t.log.WithField("error", err.Error()).Error("save doc error")
		response["remark"] = "Failed to save doc"
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	return c.SendStatus(http.StatusOK)
}

func (t *TabunganRESTAPI) getNasabah(c *fiber.Ctx) (err error) {
	response := make(map[string]interface{})
	nik := c.Get("Authorization", "")
	if nik == "" {
		err = fmt.Errorf("missing NIK in authorization header")
		t.log.Warn(err.Error())
		response["remark"] = err.Error()
		c.Status(http.StatusUnauthorized)
		return c.JSON(response)
	}
	nasabah, err := t.app.GetNasabah(nik)
	if err != nil {
		response["remark"] = err.Error()
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	response["data"] = nasabah
	return c.JSON(response)
}

func (t *TabunganRESTAPI) getDaftarRekening(c *fiber.Ctx) (err error) {
	response := make(map[string]interface{})
	nik := c.Get("Authorization", "")
	if nik == "" {
		err = fmt.Errorf("missing NIK in authorization header")
		t.log.Warn(err.Error())
		response["remark"] = err.Error()
		c.Status(http.StatusUnauthorized)
		return c.JSON(response)
	}
	rekening, err := t.app.GetDaftarRekening(nik)
	if err != nil {
		response["remark"] = err.Error()
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	response["data"] = rekening
	return c.JSON(response)
}

func (t *TabunganRESTAPI) getRekening(c *fiber.Ctx) (err error) {
	response := make(map[string]interface{})
	nik := c.Get("Authorization", "")
	if nik == "" {
		err = fmt.Errorf("missing NIK in authorization header")
		t.log.Warn(err.Error())
		response["remark"] = err.Error()
		c.Status(http.StatusUnauthorized)
		return c.JSON(response)
	}
	noRekening := c.Params("rekening", "")
	if noRekening == "" {
		err = fmt.Errorf("missing no-rekening in path parameter")
		t.log.Warn(err.Error())
		response["remark"] = err.Error()
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	rekening, err := t.app.GetRekening(nik, noRekening)
	if err != nil {
		response["remark"] = err.Error()
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	response["data"] = rekening
	return c.JSON(response)
}

func (t *TabunganRESTAPI) tarikDana(c *fiber.Ctx) (err error) {
	var request models.RequestTarikSetorDana
	response := make(map[string]interface{})
	nik := c.Get("Authorization", "")
	if nik == "" {
		err = fmt.Errorf("missing NIK in authorization header")
		t.log.Warn(err.Error())
		response["remark"] = err.Error()
		c.Status(http.StatusUnauthorized)
		return c.JSON(response)
	}
	err = c.BodyParser(&request)
	if err != nil {
		t.log.WithField("error", err.Error()).Error("parse request body to JSON error")
		response["remark"] = "Failed to parse request body"
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	saldoAkhir, err := t.app.TarikDana(nik, request.NoRekening, request.Nominal)
	if err != nil {
		response["remark"] = err.Error()
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	response["saldo_akhir"] = saldoAkhir
	return c.JSON(response)
}

func (t *TabunganRESTAPI) setorDana(c *fiber.Ctx) (err error) {
	var request models.RequestTarikSetorDana
	response := make(map[string]interface{})
	nik := c.Get("Authorization", "")
	if nik == "" {
		err = fmt.Errorf("missing NIK in authorization header")
		t.log.Warn(err.Error())
		response["remark"] = err.Error()
		c.Status(http.StatusUnauthorized)
		return c.JSON(response)
	}
	err = c.BodyParser(&request)
	if err != nil {
		t.log.WithField("error", err.Error()).Error("parse request body to JSON error")
		response["remark"] = "Failed to parse request body"
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	saldoAkhir, err := t.app.SetorDana(nik, request.NoRekening, request.Nominal)
	if err != nil {
		response["remark"] = err.Error()
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	response["saldo_akhir"] = saldoAkhir
	return c.JSON(response)
}

func (t *TabunganRESTAPI) updateNasabah(c *fiber.Ctx) (err error) {
	var request models.RequestUpdateNasabah
	response := make(map[string]interface{})
	nik := c.Get("Authorization", "")
	if nik == "" {
		err = fmt.Errorf("missing NIK in authorization header")
		t.log.Warn(err.Error())
		response["remark"] = err.Error()
		c.Status(http.StatusUnauthorized)
		return c.JSON(response)
	}
	err = c.BodyParser(&request)
	if err != nil {
		t.log.WithField("error", err.Error()).Error("parse request body to JSON error")
		response["remark"] = "Failed to parse request body"
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	err = t.app.UpdateNasabah(nik, request)
	if err != nil {
		response["remark"] = err.Error()
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	return c.SendStatus(http.StatusOK)
}

func (t *TabunganRESTAPI) getMutasi(c *fiber.Ctx) (err error) {
	response := make(map[string]interface{})
	noRekening := c.Params("rekening", "")
	if noRekening == "" {
		err = fmt.Errorf("missing no-rekening in path parameter")
		t.log.Warn(err.Error())
		response["remark"] = err.Error()
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	page := c.Query("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"no_rekening": noRekening,
			"page":        page,
			"error":       err.Error(),
		}).Error("parsing page query parameter to int failed")
		response["remark"] = "parsing page query parameter to int failed"
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	show := c.Query("show", "1")
	showInt, err := strconv.Atoi(show)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"no_rekening": noRekening,
			"show":        show,
			"error":       err.Error(),
		}).Error("parsing show query parameter to int failed")
		response["remark"] = "parsing show query parameter to int failed"
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	mutasi, err := t.app.GetMutasi(noRekening, pageInt, showInt)
	if err != nil {
		response["remark"] = err.Error()
		c.Status(http.StatusBadRequest)
		return c.JSON(response)
	}
	response["data"] = mutasi
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
	api.server.Post("/file", api.uploadFile)
	api.server.Get("/nasabah", api.getNasabah)
	api.server.Put("/nasabah", api.updateNasabah)
	api.server.Get("/rekening/list", api.getDaftarRekening)
	api.server.Get("/rekening/:rekening", api.getRekening)
	api.server.Post("/tarik", api.tarikDana)
	api.server.Post("/setor", api.setorDana)
	api.server.Get("/mutasi/:rekening", api.getMutasi)
	return api
}

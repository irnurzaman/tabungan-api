package app

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"tabungan-api/models"
	"tabungan-api/repository"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

type TabunganAppInterface interface {
	RegistrasiNasabah(request models.RequestRegistrasiNasabah) (rekening models.Rekening, err error)
	UpdateNasabah(nik string, request models.RequestUpdateNasabah) (err error)
	PembukaanRekening(nik string) (rekening models.Rekening, err error)
	GetNasabah(nik string) (nasabah models.Nasabah, err error)
	GetDaftarRekening(nik string) (rekening []string, err error)
	GetRekening(nik, noRekening string) (rekening models.Rekening, err error)
	GetMutasi(noRekening string) (mutasi []models.Mutasi, err error)
	TarikDana(nik, noRekening string, nominal float64) (saldoAkhir float64, err error)
	SetorDana(nik, noRekening string, nominal float64) (saldoAkhir float64, err error)
	SavePhoto(file io.Reader, filename, nik string) (err error)
	SaveDoc(file io.Reader, filename, nik string) (err error)
}

type TabunganApp struct {
	repo     repository.TabunganRepoInterface
	log      *logrus.Logger
	photoDir string
	docDir   string
}

func (t *TabunganApp) RegistrasiNasabah(request models.RequestRegistrasiNasabah) (rekening models.Rekening, err error) {
	var nasabah models.Nasabah
	copier.Copy(&nasabah, request)
	err = t.repo.InsertNasabah(nasabah)
	if err != nil {
		err = fmt.Errorf("registrasi nasabah gagal")
		t.log.WithFields(logrus.Fields{
			"nik":             nasabah.NIK,
			"nama":            nasabah.Nama,
			"alamat_ktp":      nasabah.AlamatKTP,
			"alamat_domisili": nasabah.AlamatDomisili,
			"jenis_kelamin":   nasabah.JenisKelamin,
			"tanggal_lahir":   nasabah.TanggalLahir,
		}).Warn(err.Error())
		return
	}
	rekening, err = t.PembukaanRekening(nasabah.NIK)
	if err != nil {
		err = fmt.Errorf("registrasi nasabah gagal")
		t.log.WithFields(logrus.Fields{
			"nik":             nasabah.NIK,
			"nama":            nasabah.Nama,
			"alamat_ktp":      nasabah.AlamatKTP,
			"alamat_domisili": nasabah.AlamatDomisili,
			"jenis_kelamin":   nasabah.JenisKelamin,
			"tanggal_lahir":   nasabah.TanggalLahir,
		}).Warn(err.Error())
	}
	return
}

func (t *TabunganApp) UpdateNasabah(nik string, request models.RequestUpdateNasabah) (err error) {
	request.NIK = nik
	err = t.repo.UpdateNasabah(request)
	if err != nil {
		err = fmt.Errorf("update nasabah error")
		t.log.WithFields(logrus.Fields{
			"nik":             nik,
			"nama":            request.Nama,
			"alamat_ktp":      request.AlamatKTP,
			"alamat_domisili": request.AlamatDomisili,
		}).Warn(err.Error())
	}
	return
}

func (t *TabunganApp) PembukaanRekening(nik string) (rekening models.Rekening, err error) {
	rekening.NIK = nik
	rekening.NoRekening = genNoRekening()
	rekening.Saldo = 0.0
	err = t.repo.InsertRekening(rekening)
	if err != nil {
		err = fmt.Errorf("pembukaan rekening gagal")
		t.log.WithFields(logrus.Fields{
			"nik":         nik,
			"no_rekening": rekening.NoRekening,
			"saldo":       rekening.Saldo,
		}).Warn(err.Error())
	}
	return
}

func (t *TabunganApp) GetNasabah(nik string) (nasabah models.Nasabah, err error) {
	nasabah, err = t.repo.GetNasabah(nik)
	if err != nil {
		err = fmt.Errorf("query data nasabah gagal")
		t.log.WithFields(logrus.Fields{
			"nik": nik,
		}).Warn(err.Error())
	}
	return
}

func (t *TabunganApp) GetDaftarRekening(nik string) (rekening []string, err error) {
	rekening, err = t.repo.GetDaftarRekening(nik)
	if err != nil {
		err = fmt.Errorf("query daftar rekening gagal")
		t.log.WithFields(logrus.Fields{
			"nik": nik,
		}).Warn(err.Error())
	}
	return
}

func (t *TabunganApp) GetRekening(nik, noRekening string) (rekening models.Rekening, err error) {
	rekening, err = t.repo.GetRekening(nik, noRekening)
	if err != nil {
		err = fmt.Errorf("query data rekening gagal")
		t.log.WithFields(logrus.Fields{
			"nik":         nik,
			"no_rekening": noRekening,
		}).Warn(err.Error())
	}
	return
}

func (t *TabunganApp) GetMutasi(noRekening string) (mutasi []models.Mutasi, err error) {
	panic("not implemented") // TODO: Implement
}

func (t *TabunganApp) TarikDana(nik, noRekening string, nominal float64) (saldoAkhir float64, err error) {
	rekening, err := t.GetRekening(nik, noRekening)
	if err != nil {
		return
	}
	if nominal > rekening.Saldo {
		err = fmt.Errorf("saldo tidak mencukupi")
		t.log.WithFields(logrus.Fields{
			"no_rekening": noRekening,
			"saldo":       rekening.Saldo,
			"nominal":     nominal,
		}).Warn("tarik dana gagal")
		return
	}
	saldoAkhir = rekening.Saldo - nominal
	err = t.repo.UpdateSaldo(noRekening, -nominal)
	if err != nil {
		err = fmt.Errorf("tarik dana rekening error")
		t.log.WithFields(logrus.Fields{
			"no_rekening": noRekening,
			"saldo":       rekening.Saldo,
			"nominal":     nominal,
		})
		return
	}
	return
}

func (t *TabunganApp) SetorDana(nik, noRekening string, nominal float64) (saldoAkhir float64, err error) {
	rekening, err := t.GetRekening(nik, noRekening)
	if err != nil {
		return
	}
	saldoAkhir = rekening.Saldo + nominal
	err = t.repo.UpdateSaldo(noRekening, nominal)
	if err != nil {
		err = fmt.Errorf("tarik dana rekening error")
		t.log.WithFields(logrus.Fields{
			"no_rekening": noRekening,
			"saldo":       rekening.Saldo,
			"nominal":     nominal,
		})
		return
	}
	return
}

func (t *TabunganApp) SavePhoto(file io.Reader, filename, nik string) (err error) {
	id, err := t.saveFile(file, t.photoDir, filename)
	if err != nil {
		err = fmt.Errorf("failed to save photo")
		t.log.WithField("nik", nik).Warn(err.Error())
		return
	}
	err = t.repo.SaveFoto(nik, id)
	if err != nil {
		err = fmt.Errorf("failed to update photoID in database")
		t.log.WithFields(logrus.Fields{
			"nik":     nik,
			"photoID": id,
		}).Warn(err.Error())
	}
	return
}

func (t *TabunganApp) SaveDoc(file io.Reader, filename, nik string) (err error) {
	id, err := t.saveFile(file, t.docDir, filename)
	if err != nil {
		err = fmt.Errorf("failed to save document")
		t.log.WithField("nik", nik).Warn(err.Error())
		return
	}
	err = t.repo.SaveDokumen(nik, id)
	if err != nil {
		err = fmt.Errorf("failed to update documentID in database")
		t.log.WithFields(logrus.Fields{
			"nik":        nik,
			"documentID": id,
		}).Warn(err.Error())
	}
	return
}

func (t *TabunganApp) saveFile(file io.Reader, folder, filename string) (id string, err error) {
	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"folder": folder,
			"error":  err.Error(),
		}).Error("failed to create directory if not exists")
		return
	}
	ext := filepath.Ext(filename)
	id = fmt.Sprintf("%s%s", genID(), ext)
	path := fmt.Sprintf("%s/%s", folder, id)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error()
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"error": err.Error(),
		})
	}
	return
}

func genNoRekening() (noRekening string) {
	noRekening = strconv.Itoa(10000000 + rand.Intn(89999999))
	return
}

func genID() string {
	return uuid.NewString()
}

func NewTabunganApp(photoDir, docDir string, repo repository.TabunganRepoInterface, log *logrus.Logger) (app *TabunganApp) {
	return &TabunganApp{
		repo:     repo,
		log:      log,
		photoDir: photoDir,
		docDir:   docDir,
	}
}

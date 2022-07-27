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
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TabunganAppInterface interface {
	RegistrasiNasabah(request models.RequestRegistrasiNasabah) (rekening models.Rekening, err error)
	UpdateNasabah(nik string, request models.RequestUpdateNasabah) (err error)
	PembukaanRekening(tx *sqlx.Tx, nik string) (rekening models.Rekening, err error)
	GetNasabah(nik string) (nasabah models.Nasabah, err error)
	GetDaftarRekening(nik string) (rekening []string, err error)
	GetRekening(nik, noRekening string) (rekening models.Rekening, err error)
	GetMutasi(noRekening string, page, show int) (mutasi []models.Mutasi, err error)
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
	tx, err := t.repo.StartTransaction()
	if err != nil {
		err = fmt.Errorf("registrasi nasabah error")
		t.log.WithFields(logrus.Fields{
			"nik":             request.NIK,
			"nama":            nasabah.Nama,
			"alamat_ktp":      nasabah.AlamatKTP,
			"alamat_domisili": nasabah.AlamatDomisili,
			"jenis_kelamin":   nasabah.JenisKelamin,
			"tanggal_lahir":   nasabah.TanggalLahir,
		}).Warn(err.Error())
		tx.Rollback()
		return
	}
	err = t.repo.InsertNasabah(tx, nasabah)
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
		tx.Rollback()
		return
	}
	rekening, err = t.PembukaanRekening(tx, nasabah.NIK)
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
		tx.Rollback()
	}
	tx.Commit()
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

func (t *TabunganApp) PembukaanRekening(tx *sqlx.Tx, nik string) (rekening models.Rekening, err error) {
	rekening.NIK = nik
	rekening.NoRekening = genNoRekening()
	rekening.Saldo = 0.0
	err = t.repo.InsertRekening(tx, rekening)
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

func (t *TabunganApp) GetMutasi(noRekening string, page, show int) (mutasi []models.Mutasi, err error) {
	offset := (page - 1) * show
	mutasi, err = t.repo.GetMutasi(noRekening, show, offset)
	if err != nil {
		err = fmt.Errorf("query data mutasi gagal")
		t.log.WithFields(logrus.Fields{
			"no_rekening": noRekening,
			"page":        page,
			"show":        show,
		}).Warn(err.Error())
	}
	return
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
		}).Warn("tarik dana gagal")
		return
	}
	err = t.insertMutasi(noRekening, "D", nominal, rekening.Saldo, saldoAkhir)
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
	err = t.insertMutasi(noRekening, "C", nominal, rekening.Saldo, saldoAkhir)
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

func (t *TabunganApp) insertMutasi(noRekening, jenisMutasi string, nominal, saldoAwal, saldoAkhir float64) (err error) {
	mutasi := models.Mutasi{
		TransaksiID: genID(),
		Waktu:       time.Now().String(),
		NoRekening:  noRekening,
		JenisMutasi: jenisMutasi,
		Nominal:     nominal,
		SaldoAwal:   saldoAwal,
		SaldoAkhir:  saldoAkhir,
	}
	err = t.repo.InsertMutasi(mutasi)
	if err != nil {
		err = fmt.Errorf("pencatatan transaksi gagal")
		t.log.WithFields(logrus.Fields{
			"no_rekening":  noRekening,
			"jenis_mutasi": jenisMutasi,
			"nominal":      nominal,
			"saldo_awal":   saldoAwal,
			"saldo_akhir":  saldoAkhir,
		}).Warn("pencatatan transaksi gagal")
		return
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
			"path":  path,
		}).Error("failed to open file")
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"path":  path,
		}).Error("failed to write file")
	}
	return
}

func genNoRekening() (noRekening string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	noRekening = strconv.Itoa(10000000 + r.Intn(89999999))
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

package app

import (
	"fmt"
	"math/rand"
	"strconv"
	"tabungan-api/models"
	"tabungan-api/repository"

	"github.com/sirupsen/logrus"
)

type TabunganAppInterface interface {
	RegistrasiNasabah(nasabah models.Nasabah) (err error)
	PembukaanRekening(nik string) (rekening models.Rekening, err error)
	GetNasabah(nik string) (nasabah models.Nasabah, err error)
	GetDaftarRekening(nik string) (rekening []models.Rekening, err error)
	GetMutasi(noRekening string) (mutasi []models.Mutasi, err error)
	TarikDana(noRekening string, nominal float64) (err error)
	SetorDana(noRekening string, nominal float64) (err error)
}

type TabunganApp struct {
	repo repository.TabunganRepoInterface
	log *logrus.Logger
}

func (t *TabunganApp) RegistrasiNasabah(nasabah models.Nasabah) (err error) {
	err = t.repo.InsertNasabah(nasabah)
	if err != nil {
		err = fmt.Errorf("registrasi nasabah gagal")
		t.log.WithFields(logrus.Fields{
			"nik": nasabah.NIK,
			"nama": nasabah.Nama,
			"alamat_ktp": nasabah.AlamatKTP,
			"alamat_domisili": nasabah.AlamatDomisili,
			"jenis_kelamin": nasabah.JenisKelamin,
			"tanggal_lahir": nasabah.TanggalLahir,
		}).Warn(err.Error())
		return
	}
	rekening, err := t.PembukaanRekening(nasabah.NIK)
	if err != nil {
		err = fmt.Errorf("pembukaan rekening gagal")
		t.log.WithFields(logrus.Fields{
			"nik": nasabah.NIK,
			"no_rekening": rekening.NoRekening,
			"saldo": rekening.Saldo,
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
			"nik": nik,
			"no_rekening": rekening.NoRekening,
			"saldo": rekening.Saldo,
		}).Warn(err.Error())
	}
	return
}

func (t *TabunganApp) GetNasabah(nik string) (nasabah models.Nasabah, err error) {
	panic("not implemented") // TODO: Implement
}

func (t *TabunganApp) GetDaftarRekening(nik string) (rekening []models.Rekening, err error) {
	panic("not implemented") // TODO: Implement
}

func (t *TabunganApp) GetMutasi(noRekening string) (mutasi []models.Mutasi, err error) {
	panic("not implemented") // TODO: Implement
}

func (t *TabunganApp) TarikDana(noRekening string, nominal float64) (err error) {
	panic("not implemented") // TODO: Implement
}

func (t *TabunganApp) SetorDana(noRekening string, nominal float64) (err error) {
	panic("not implemented") // TODO: Implement
}

func genNoRekening() (noRekening string) {
	noRekening = strconv.Itoa(10000000 + rand.Intn(89999999))
	return
}
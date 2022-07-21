package app

import "tabungan-api/models"

type TabunganApp interface {
	RegistrasiNasabah(nasabah models.Nasabah) (err error)
	PembukaanRekening(nik string) (rekening models.Rekening, err error)
	GetNasabah(nik string) (nasabah models.Nasabah, err error)
	GetDaftarRekening(nik string) (rekening []models.Rekening, err error)
	GetMutasi(noRekening string) (mutasi []models.Mutasi, err error)
	TarikDana(noRekening string, nominal float64) (err error)
	SetorDana(noRekening string, nominal float64) (err error)
}

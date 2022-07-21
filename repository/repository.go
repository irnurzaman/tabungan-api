package repository

import "tabungan-api/models"

type TabunganRepo interface {
	InsertNasabah(nasabah models.Nasabah) (err error)
	GetNasabah(nik string) (nasabah models.Nasabah, err error)
	UpdateNasabah(nasabah models.Nasabah) (err error)
	InsertRekening(rekening models.Rekening) (err error)
	GetDaftarRekening(nik string) (rekening []models.Rekening, err error)
	GetRekening(nik string) (rekening models.Rekening, err error)
	InsertMutasi(mutasi models.Mutasi) (err error)
	GetMutasi(noRekening string, page int, show int) (mutasi []models.Mutasi)
	TarikDana(noRekening string, nominal float64) (err error)
	SetorDana(noRekening string, nominal float64) (err error)
}

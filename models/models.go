package models

import "time"

type Nasabah struct {
	NIK            string    `json:"nik" db:"nik"`
	Nama           string    `json:"nama" db:"nama"`
	AlamatKTP      string    `json:"alamat_ktp" db:"alamat_ktp"`
	AlamatDomisili string    `json:"alamat_domisili" db:"alamat_domisili"`
	JenisKelamin   string    `json:"jenis_kelamin" db:"jenis_kelamin"`
	TanggalLahir   time.Time `json:"tanggal_lahir" db:"tanggal_lahir"`
	PhotoID        string    `json:"photo_id" db:"photo_id"`
	DokumenID      string    `json:"dokumen_id" db:"dokumen_id"`
}

type Rekening struct {
	NIK        string  `json:"nik" db:"nik"`
	NoRekening string  `json:"no_rekening" db:"no_rekening"`
	Saldo      float64 `json:"saldo" db:"saldo"`
}

type Mutasi struct {
	TransaksiID string    `json:"transaksi_id" db:"transaksi_id"`
	Waktu       time.Time `json:"waktu" db:"waktu"`
	JenisMutasi string    `json:"jenis_mutasi" db:"jenis_mutasi"`
	NoRekening  string    `json:"no_rekening" db:"no_rekening"`
	Nominal     float64   `json:"nominal" db:"nominal"`
	SaldoAwal   float64   `json:"saldo_awal"`
	SaldoAkhir  float64   `json:"saldo_akhir"`
}

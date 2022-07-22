package repository

import (
	"fmt"
	"tabungan-api/models"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
)

type TabunganRepoInterface interface {
	StartTransaction() (tx *sqlx.Tx, err error)
	InsertNasabah(nasabah models.Nasabah) (err error)
	GetNasabah(nik string) (nasabah models.Nasabah, err error)
	UpdateNasabah(nasabah models.Nasabah) (err error)
	SaveFoto(nik string, fotoID string) (err error)
	SaveDokumen(nik string, dokumenID string) (err error)
	InsertRekening(rekening models.Rekening) (err error)
	GetDaftarRekening(nik string) (rekening []models.Rekening, err error)
	GetRekening(noRekening string) (rekening models.Rekening, err error)
	InsertMutasi(mutasi models.Mutasi) (err error)
	GetMutasi(noRekening string, page int, show int) (mutasi []models.Mutasi)
	TarikDana(noRekening string, nominal float64) (err error)
	SetorDana(noRekening string, nominal float64) (err error)
}

type TabunganRepo struct {
	db  *sqlx.DB
	log *logrus.Logger
}

func (t *TabunganRepo) initDatabase() {
	SQL := `CREATE TABLE IF NOT EXISTS nasabah (
		nik text PRIMARY KEY,
		nama text,
		alamat_ktp text,
		alamat_domisili text,
		jenis_kelamin text,
		tanggal_lahir text,
		foto_id text,
		dokumen_id text);`
	t.db.MustExec(SQL)

	SQL = `CREATE TABLE IF NOT EXISTS rekening (
		nik text,
		no_rekening text PRIMARY KEY,
		saldo real);`
	t.db.MustExec(SQL)

	SQL = `CREATE TABLE IF NOT EXISTS mutasi (
		transaksi_id text PRIMARY KEY,
		waktu text,
		jenis_mutasi text,
		no_rekening text,
		nominal real,
		saldo_awal real,
		saldo_akhir text);`
	t.db.MustExec(SQL)
}

func (t *TabunganRepo) StartTransaction() (tx *sqlx.Tx, err error) {
	tx, err = t.db.Beginx()
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("start transaction error")
	}
	return
}

func (t *TabunganRepo) InsertNasabah(tx *sqlx.Tx, nasabah models.Nasabah) (err error) {
	SQL := "INSERT INTO nasabah VALUES (:nik, :nama, :alamat_ktp, :alamat_domisili, :jenis_kelamin, :tanggal_lahir, :foto_id, :dokumen_id)"
	_, err = tx.NamedExec(SQL, nasabah)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"nik":             nasabah.NIK,
			"nama":            nasabah.Nama,
			"alamat_ktp":      nasabah.AlamatKTP,
			"alamat_domisili": nasabah.AlamatDomisili,
			"jenis_kelamin":   nasabah.JenisKelamin,
			"tanggal_lahir":   nasabah.TanggalLahir,
			"error":           err.Error(),
		}).Error("insert data nasabah error")
	}
	return
}

func (t *TabunganRepo) GetNasabah(nik string) (nasabah models.Nasabah, err error) {
	SQL := "SELECT * FROM nasabah WHERE nik = $1"
	err = t.db.Get(&nasabah, SQL, nik)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"nik":     nik,
			"error":   err.Error(),
		}).Error("get nasabah error")
	}
	return
}

func (t *TabunganRepo) UpdateNasabah(nasabah models.Nasabah) (err error) {
	SQL := "UPDATE nasabah SET nama = :nama, alamat_ktp = :alamat_ktp, alamat_domisili = :alamat_domisili WHERE nik = :nik"
	_, err = t.db.NamedExec(SQL, nasabah)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"nik":     nasabah.NIK,
			"nama": nasabah.Nama,
			"alamat_ktp": nasabah.AlamatKTP,
			"alamat_domisili": nasabah.AlamatDomisili,
			"error":   err.Error(),
		}).Error("update data nasabah error")
	}
	return
}

func (t *TabunganRepo) SaveFoto(nik string, fotoID string) (err error) {
	SQL := "UPDATE nasabah SET foto_id = $1 WHERE nik = $2"
	_, err = t.db.Exec(SQL, fotoID, nik)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"nik":     nik,
			"foto_id": fotoID,
			"error":   err.Error(),
		}).Error("update nasabah foto_id error")
	}
	return
}

func (t *TabunganRepo) SaveDokumen(nik string, dokumenID string) (err error) {
	SQL := "UPDATE nasabah SET dokumen_id = $1 WHERE nik = $2"
	_, err = t.db.Exec(SQL, dokumenID, nik)
	if err != nil {
		err = fmt.Errorf("update nasabah dokumen_id error")
		t.log.WithFields(logrus.Fields{
			"nik":        nik,
			"dokumen_id": dokumenID,
		}).Error(err.Error())
	}
	return
}

func (t *TabunganRepo) InsertRekening(rekening models.Rekening) (err error) {
	panic("not implemented") // TODO: Implement
}

func (t *TabunganRepo) GetDaftarRekening(nik string) (rekening []models.Rekening, err error) {
	panic("not implemented") // TODO: Implement
}

func (t *TabunganRepo) GetRekening(nik string) (rekening models.Rekening, err error) {
	panic("not implemented") // TODO: Implement
}

func (t *TabunganRepo) InsertMutasi(mutasi models.Mutasi) (err error) {
	panic("not implemented") // TODO: Implement
}

func (t *TabunganRepo) GetMutasi(noRekening string, page int, show int) (mutasi []models.Mutasi) {
	panic("not implemented") // TODO: Implement
}

func (t *TabunganRepo) TarikDana(noRekening string, nominal float64) (err error) {
	panic("not implemented") // TODO: Implement
}

func (t *TabunganRepo) SetorDana(noRekening string, nominal float64) (err error) {
	panic("not implemented") // TODO: Implement
}

func InitDatabase(database string, logger *logrus.Logger) (repo *TabunganRepo) {
	db, err := sqlx.Connect("sqlite3", database)
	if err != nil {
		panic(err)
	}

	repo = &TabunganRepo{
		db: db,
		log: logger,
	}
	repo.initDatabase()
	return
}
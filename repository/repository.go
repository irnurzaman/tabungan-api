package repository

import (
	"fmt"
	"tabungan-api/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type TabunganRepoInterface interface {
	StartTransaction() (tx *sqlx.Tx, err error)
	InsertNasabah(tx *sqlx.Tx, nasabah models.Nasabah) (err error)
	GetNasabah(nik string) (nasabah models.Nasabah, err error)
	UpdateNasabah(nasabah models.RequestUpdateNasabah) (err error)
	SaveFoto(nik string, fotoID string) (err error)
	SaveDokumen(nik string, dokumenID string) (err error)
	InsertRekening(tx *sqlx.Tx, rekening models.Rekening) (err error)
	GetDaftarRekening(nik string) (rekening []string, err error)
	GetRekening(nik, noRekening string) (rekening models.Rekening, err error)
	InsertMutasi(tx *sqlx.Tx, mutasi models.Mutasi) (err error)
	GetMutasi(noRekening string, limit, offset int) (mutasi []models.Mutasi, err error)
	UpdateSaldo(tx *sqlx.Tx, noRekening string, nominal float64) (err error)
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
		t.log.WithField("error", err.Error()).Error("begin transaction error")
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
			"nik":   nik,
			"error": err.Error(),
		}).Error("get nasabah error")
	}
	return
}

func (t *TabunganRepo) UpdateNasabah(nasabah models.RequestUpdateNasabah) (err error) {
	SQL := "UPDATE nasabah SET nama = :nama, alamat_ktp = :alamat_ktp, alamat_domisili = :alamat_domisili WHERE nik = :nik"
	_, err = t.db.NamedExec(SQL, nasabah)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"nik":             nasabah.NIK,
			"nama":            nasabah.Nama,
			"alamat_ktp":      nasabah.AlamatKTP,
			"alamat_domisili": nasabah.AlamatDomisili,
			"error":           err.Error(),
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

func (t *TabunganRepo) InsertRekening(tx *sqlx.Tx, rekening models.Rekening) (err error) {
	SQL := "INSERT INTO rekening VALUES (:nik, :no_rekening, :saldo)"
	_, err = tx.NamedExec(SQL, rekening)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"nik":         rekening.NIK,
			"no_rekening": rekening.NoRekening,
			"saldo":       rekening.Saldo,
			"error":       err.Error(),
		}).Error("insert rekening error")
	}
	return
}

func (t *TabunganRepo) GetDaftarRekening(nik string) (rekening []string, err error) {
	SQL := "SELECT no_rekening FROM rekening WHERE nik = $1"
	err = t.db.Select(&rekening, SQL, nik)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"nik":   nik,
			"error": err.Error(),
		}).Error("get rekening error")
	}
	return
}

func (t *TabunganRepo) GetRekening(nik, noRekening string) (rekening models.Rekening, err error) {
	SQL := "SELECT * FROM rekening WHERE nik = $1 AND no_rekening = $2"
	err = t.db.Get(&rekening, SQL, nik, noRekening)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"nik":         nik,
			"no_rekening": noRekening,
			"error":       err.Error(),
		}).Error("get rekening error")
	}
	return
}

func (t *TabunganRepo) InsertMutasi(tx *sqlx.Tx, mutasi models.Mutasi) (err error) {
	SQL := "INSERT INTO mutasi VALUES (:transaksi_id, :waktu, :jenis_mutasi, :no_rekening, :nominal, :saldo_awal, :saldo_akhir)"
	_, err = t.db.NamedExec(SQL, mutasi)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"transaksi_id": mutasi.TransaksiID,
			"waktu":        mutasi.Waktu,
			"jenis_mutasi": mutasi.JenisMutasi,
			"no_rekening":  mutasi.NoRekening,
			"nominal":      mutasi.Nominal,
			"saldo_awal":   mutasi.SaldoAwal,
			"saldo_akhir":  mutasi.SaldoAkhir,
			"error":        err.Error(),
		}).Error("insert mutasi error")
	}
	return
}

func (t *TabunganRepo) GetMutasi(noRekening string, limit, offset int) (mutasi []models.Mutasi, err error) {
	SQL := "SELECT * FROM mutasi WHERE no_rekening = $1 LIMIT $2 OFFSET $3"
	err = t.db.Select(&mutasi, SQL, noRekening, limit, offset)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"no_rekening": noRekening,
			"limit":       limit,
			"offset":      offset,
			"error":       err.Error(),
		}).Error("query mutasi error")
	}
	return
}

func (t *TabunganRepo) UpdateSaldo(tx *sqlx.Tx, noRekening string, nominal float64) (err error) {
	SQL := "UPDATE rekening SET saldo = saldo + $1 WHERE no_rekening = $2"
	_, err = t.db.Exec(SQL, nominal, noRekening)
	if err != nil {
		t.log.WithFields(logrus.Fields{
			"no_rekening": noRekening,
			"error":       err.Error(),
		}).Error("tarik dana rekening error")
	}
	return
}

func InitDatabase(database string, logger *logrus.Logger) (repo *TabunganRepo) {
	db, err := sqlx.Connect("sqlite3", database)
	if err != nil {
		panic(err)
	}

	repo = &TabunganRepo{
		db:  db,
		log: logger,
	}
	repo.initDatabase()
	return
}

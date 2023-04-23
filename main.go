package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Mahasiswa struct {
	Id                 int    `json:"id"`
	Nama               string `json:"nama"`
	Usia               int    `json:"usia"`
	Gender             int    `json:"gender"`
	Tanggal_Registrasi string `json:"tanggal_registrasi"`
}

type Jurusan struct {
	Id           int    `json:"id"`
	Nama_Jurusan string `json:"nama_jurusan"`
}

type Hobi struct {
	Id        int    `json:"id"`
	Nama_Hobi string `json:"nama_hobi"`
}

var db *sql.DB

func main() {
	db, _ = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/db_jobhun")
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/mahasiswa", getMahasiswas).Methods("GET")
	router.HandleFunc("/mahasiswa/{id}", getMahasiswa).Methods("GET")
	router.HandleFunc("/mahasiswa", createMahasiswa).Methods("POST")
	router.HandleFunc("/mahasiswa/{id}", updateMahasiswa).Methods("PUT")
	router.HandleFunc("/mahasiswa/{id}", deleteMahasiswa).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8091", router))
}

func getMahasiswas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var mahasiswas []Mahasiswa

	result, err := db.Query("SELECT * from Mahasiswa")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {

		var mahasiswa Mahasiswa
		err := result.Scan(&mahasiswa.Id, &mahasiswa.Nama, &mahasiswa.Usia, &mahasiswa.Gender, &mahasiswa.Tanggal_Registrasi)
		if err != nil {
			panic(err.Error())
		}
		mahasiswas = append(mahasiswas, mahasiswa)
	}

	json.NewEncoder(w).Encode(mahasiswas)
}

func getMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT * FROM Mahasiswa WHERE Id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var mahasiswa Mahasiswa

	for result.Next() {
		err := result.Scan(&mahasiswa.Id, &mahasiswa.Nama, &mahasiswa.Usia, &mahasiswa.Gender, &mahasiswa.Tanggal_Registrasi)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(mahasiswa)
}

func createMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var mahasiswa Mahasiswa
	_ = json.NewDecoder(r.Body).Decode(&mahasiswa)

	insertQuery := "INSERT INTO Mahasiswa (Nama, Usia, Gender, Tanggal_Registrasi) VALUES(?,?,?,?)"

	stmt, err := db.Prepare(insertQuery)
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(mahasiswa.Nama, mahasiswa.Usia, mahasiswa.Gender, mahasiswa.Tanggal_Registrasi)
	if err != nil {
		panic(err.Error())
	}

	defer stmt.Close()

	json.NewEncoder(w).Encode(mahasiswa)
}

func updateMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var mahasiswa Mahasiswa
	_ = json.NewDecoder(r.Body).Decode(&mahasiswa)

	updateQuery := "UPDATE Mahasiswa SET Nama=?, Usia=?, Gender=?, Tanggal_Registrasi=? WHERE Id=?"

	stmt, err := db.Prepare(updateQuery)
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(mahasiswa.Nama, mahasiswa.Usia, mahasiswa.Gender, mahasiswa.Tanggal_Registrasi, params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer stmt.Close()

	json.NewEncoder(w).Encode(mahasiswa)
}

func deleteMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	deleteQuery := "DELETE FROM Mahasiswa WHERE Id = ?"

	stmt, err := db.Prepare(deleteQuery)
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer stmt.Close()

	json.NewEncoder(w).Encode("Mahasiswa telah dihapus")
}

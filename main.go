package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// Barang struct (Model) ...
type Barang struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
	Jenis string `json:"jenis"`
	Satuan string `json:"satuan"`
	Harga string `json:"harga"`
}

// Get all orders

func getBarangs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var barangList []Barang

	sql := `SELECT
				kode,
				IFNULL(nama,''),
				IFNULL(jenis,'') jenis,
				IFNULL(satuan,'') satuan,
				IFNULL(harga,'') harga
			FROM barang`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var barang Barang
		err := result.Scan(&barang.Kode, &barang.Nama, &barang.Jenis,
			&barang.Satuan, &barang.Harga)

		if err != nil {
			panic(err.Error())
		}
		barangList = append(barangList, barang)
	}

	json.NewEncoder(w).Encode(barangList)
}

func createBarang(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		kode := r.FormValue("kode")
		nama := r.FormValue("nama")
		jenis := r.FormValue("jenis")
		satuan := r.FormValue("satuan")
		harga := r.FormValue("harga")

		stmt, err := db.Prepare("INSERT INTO barang (kode,nama,jenis,satuan,harga) VALUES (?,?,?,?,?)")

		_, err = stmt.Exec(kode, nama, jenis, satuan, harga)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}

func getBarang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var barangList []Barang
	params := mux.Vars(r)

	sql := `SELECT
				kode,
				IFNULL(nama,''),
				IFNULL(jenis,'') jenis,
				IFNULL(satuan,'') satuan,
				IFNULL(harga,'') harga
			FROM barang WHERE kode = ?`

	result, err := db.Query(sql, params["id"])

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var barang Barang

	for result.Next() {

		err := result.Scan(&barang.Kode, &barang.Nama, &barang.Jenis,
			&barang.Satuan, &barang.Harga)

		if err != nil {
			panic(err.Error())
		}

		barangList = append(barangList, barang)
	}

	json.NewEncoder(w).Encode(barang)
}

func updateBarang(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newnama := r.FormValue("nama")

		stmt, err := db.Prepare("UPDATE barang SET nama = ? WHERE kode = ?")

		_, err = stmt.Exec(newnama, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Barang with kode = %s was updated", params["id"])
	}
}

func deleteBarang(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM barang WHERE kode = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Barang with kode = %s was deleted", params["id"]) //
}

func getPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var barangList []Barang

	kode := r.FormValue("kode")
	nama := r.FormValue("nama")

	sql := `SELECT
				kode,
				IFNULL(nama,''),
				IFNULL(jenis,'') jenis,
				IFNULL(satuan,'') satuan,
				IFNULL(harga,'') harga
			FROM barang WHERE kode = ? AND nama = ?`

	result, err := db.Query(sql, kode, nama)

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var barang Barang

	for result.Next() {

		err := result.Scan(&barang.Kode, &barang.Nama, &barang.Jenis,
			&barang.Satuan, &barang.Harga)

		if err != nil {
			panic(err.Error())
		}

		barangList = append(barangList, barang)
	}

	json.NewEncoder(w).Encode(barang)

}

// Main function
func main() {

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/barang")
	if err != nil {
		panic(err.Error())
	}
 
	defer db.Close()

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/barang", getBarangs).Methods("GET")
	r.HandleFunc("/barang/{id}", getBarang).Methods("GET")
	r.HandleFunc("/barang", createBarang).Methods("POST")
	r.HandleFunc("/barang/{id}", updateBarang).Methods("PUT")
	r.HandleFunc("/barang/{id}", deleteBarang).Methods("DELETE")

	//New
	r.HandleFunc("/getbarang", getPost).Methods("POST")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}

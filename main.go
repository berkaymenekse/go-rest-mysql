package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Il struct {
	ilno uint8
	isim string
}

type Ilce struct {
	ilce_no uint8
	isim    string
	il_no   uint8
}

type Page struct {
	title string
}

func createConnection() bool {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		AllowNativePasswords: true,
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "go_test",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return true
}

func getIller() ([]Il, error) {

	var iller []Il
	rows, err := db.Query("SELECT * FROM iller")
	if err != nil {
		return nil, fmt.Errorf("iller gelirken problem oldu")
	}
	defer rows.Close()
	for rows.Next() {
		var il Il
		if err := rows.Scan(&il.ilno, &il.isim); err != nil {
			return nil, fmt.Errorf("Degerleri parse ederken bi problem oldu")
		}
		iller = append(iller, il)
	}
	return iller, nil
}

func getIl(il_no uint8) (Il, error) {
	var il Il
	row := db.QueryRow("SELECT * FROM iller WHERE  il_no = ?", il_no)

	if err := row.Scan(&il.ilno, &il.isim); err != nil {
		if err == sql.ErrNoRows {
			return il, fmt.Errorf("boyle bir il yok %d", il_no)
		}
		return il, fmt.Errorf("getIl %d: %v", il_no, err)
	}
	return il, nil
}

func getIlceler(il_no uint8) ([]Ilce, error) {

	var ilceler []Ilce
	rows, err := db.Query("SELECT * FROM ilceler WHERE il_no = ?", il_no)
	if err != nil {
		return nil, fmt.Errorf("ilceler gelirken problem oldu")
	}
	defer rows.Close()
	for rows.Next() {
		var ilce Ilce
		if err := rows.Scan(&ilce.ilce_no, &ilce.isim, &ilce.il_no); err != nil {
			return nil, fmt.Errorf("Degerleri parse ederken bi problem oldu")
		}
		ilceler = append(ilceler, ilce)
	}
	return ilceler, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})
	if createConnection() {
		iller, err := getIller()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("iller ", iller)
		il, err := getIl(59)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Il", il)
		ilceler, err := getIlceler(59)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Il", ilceler)
	}

	http.ListenAndServe(":8097", nil)
}

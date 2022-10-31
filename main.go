package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/microsoft/go-mssqldb"
)

var db *sql.DB
var server = "aksserverdemo.database.windows.net"
var port = 1433
var user = "vamshi"
var password = "A1*krishna"
var database = "aksdatabase"

type studentinfo struct {
	Sid    string `json::"sid,omitempty"`
	Name   string `json::"name,omitempty"`
	Course string `json::"course,omitempty"`
}

// MSSQL DB configuration
func GetMySQLDB() *sql.DB {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)

	if err != nil {
		log.Fatal(err)
	}
	return db
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	db := GetMySQLDB()
	defer db.Close()
	ss := []studentinfo{}
	s := studentinfo{}
	rows, err := db.Query("select * from student")
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		for rows.Next() {
			rows.Scan(&s.Sid, &s.Name, &s.Course)
			ss = append(ss, s)
		}
		json.NewEncoder(w).Encode(ss)

	}

}

func addStudents(w http.ResponseWriter, r *http.Request) {
	db := GetMySQLDB()
	defer db.Close()
	s := studentinfo{}
	json.NewDecoder(r.Body).Decode(&s)
	sid, _ := strconv.Atoi(s.Sid)
	tsql := fmt.Sprintf("insert into student (sid, name, course) values('%d','%s','%s')", sid, s.Name, s.Course)
	result, err := db.Exec(tsql)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.LastInsertId()
		if err != nil {
			json.NewEncoder(w).Encode(s)
		} else {
			json.NewEncoder(w).Encode(s)

		}
	}

}

func updateStudents(w http.ResponseWriter, r *http.Request) {
	db := GetMySQLDB()
	defer db.Close()
	s := studentinfo{}
	json.NewDecoder(r.Body).Decode(&s)
	vars := mux.Vars(r)
	sid, _ := strconv.Atoi(vars["sid"])
	tsql := fmt.Sprintf("update student set name='%s', course='%s' where sid='%d'", s.Name, s.Course, sid)
	result, err := db.Exec(tsql)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.RowsAffected()
		if err != nil {
			json.NewEncoder(w).Encode("{ error: record not updated }")
		} else {
			json.NewEncoder(w).Encode(s)

		}
	}
}

func deleteStudents(w http.ResponseWriter, r *http.Request) {
	db := GetMySQLDB()
	defer db.Close()
	vars := mux.Vars(r)
	sid, _ := strconv.Atoi(vars["sid"])
	tsql := fmt.Sprintf("delete from student where sid='%d';", sid)
	result, err := db.Exec(tsql)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.RowsAffected()
		if err != nil {
			json.NewEncoder(w).Encode("{ error: record not Deleted }")
		} else {
			json.NewEncoder(w).Encode("{ result: record Sucessfully Deleted }")

		}
	}

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students", addStudents).Methods("POST")
	r.HandleFunc("/students/{sid}", updateStudents).Methods("PUT")
	r.HandleFunc("/students/{sid}", deleteStudents).Methods("DELETE")
	fmt.Println("server started")
	http.ListenAndServe(":8000", r)
}

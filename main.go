package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"

    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    Id    string
    Name  string
    Age   int
    Email string
}

func dbConn() (db *sql.DB) {
    db, err := sql.Open("mysql", "root:password@/go_api")
    if err != nil {
        panic(err.Error())
    }
    return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM User ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
    }
    emp := User{}
    res := []User{}
    for selDB.Next() {
        var Id string
        var Name string
        var Age int
        var Email string
        err = selDB.Scan(&Id, &Name, &Age, &Email)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = Id
        emp.Name = Name
        emp.Age = Age
        emp.Email = Email
        res = append(res, emp)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("Id")
    selDB, err := db.Query("SELECT * FROM User WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := User{}
    for selDB.Next() {
    var Id string
    var Name string
    var Age int
    var Email string
    err = selDB.Scan(&Id, &Name, &Age, &Email)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = Id
        emp.Name = Name
        emp.Age = Age
        emp.Email = Email
    }
    tmpl.ExecuteTemplate(w, "Show", emp)
    defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("Id")
    selDB, err := db.Query("SELECT * FROM User WHERE Id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := User{}
    for selDB.Next() {
    var Id string
    var Name string
    var Age int
    var Email string
    err = selDB.Scan(&Id, &Name, &Age, &Email)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = Id
        emp.Name = Name
        emp.Age = Age
        emp.Email = Email
    }
    tmpl.ExecuteTemplate(w, "Edit", emp)
    defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        Id := r.FormValue("Id")
        Name := r.FormValue("Name")
        Age := r.FormValue("Age")
        Email := r.FormValue("Email")
        insForm, err := db.Prepare("INSERT INTO User(Id, Name, Age, Email) VALUES(?,?,?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(Id, Name, Age, Email)
        log.Println("INSERT: Id: " + Id + " | Name: " + Name + " | Age: " + Age + " | Email: " + Email)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        Name := r.FormValue("Name")
        Age := r.FormValue("Age")
        Email := r.FormValue("Email")
        Id := r.FormValue("uid")
        insForm, err := db.Prepare("UPDATE User SET Name=?, Age=?, Email=? WHERE Id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(Name, Age, Email, Id)
        log.Println("UPDATE:  Name: " + Name + " | Age: " + Age + " | Email: " + Email)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    emp := r.URL.Query().Get("Id")
    delForm, err := db.Prepare("DELETE FROM User WHERE Id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(emp)
    log.Println("DELETE")
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func main() {
    log.Println("Server started on: http://localhost:8080")
    http.HandleFunc("/", Index)
    http.HandleFunc("/show", Show)
    http.HandleFunc("/new", New)
    http.HandleFunc("/edit", Edit)
    http.HandleFunc("/insert", Insert)
    http.HandleFunc("/update", Update)
    http.HandleFunc("/delete", Delete)
    http.ListenAndServe(":8080", nil)
}

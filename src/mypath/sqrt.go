package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

func main(){
    db,err := sql.Open("mysql","wustan:websocket@tcp(127.0.0.1:3306)/websocket")
    if err != nil{
        panic(err.Error())
    }
    fmt.Println("conn suc")
    defer db.Close()

}

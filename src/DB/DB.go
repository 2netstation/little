package main
import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

//func getTask (date string, time string, day string, status int) (bool) {
    //db, err := sql.Open("mysql", "wustan:websocket@(127.0.0.1:3306)/websocket")
    //if err!= nil {
        //panic(err.Error())
    //}
    //defer db.Close()
    //for i := 0; i < 10; i++{
        //var (
            //time_from = "2016-03-27 12:00:00"


        //)

        //result, err := db.Exec(
            //"INSERT INTO task (time_from, time_to, day, trigger,callback, status, create_time) VALUES ($1, $2, $3, $4, $5, $6, $7)",
            //""
        //)
    //}
    //defer rows.Close()

//}






func main(){
    db,err := sql.Open("mysql","wustan:websocket@(127.0.0.1:3306)/websocket")

    //need use propur error handler
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()
    fmt.Println("db open!")

    err = db.Ping()
    if err != nil{
        panic(err.Error())
    }
    fmt.Println("db connect !")

    var(
        id int
        process_id int
        process_name string
        process_type int
        last_time string
        add_time string


    )

    rows,err := db.Query("select * from websocket limit  1")
    if err != nil{
        panic(err.Error())
    }
    defer rows.Close()

    for rows.Next(){
        err := rows.Scan(&id,&process_id,&process_name,&process_type,&last_time,&add_time)
        if err != nil{
            panic(err.Error())
        }
        fmt.Println(id,process_id,process_name,process_type,last_time,add_time)

    }
    err = rows.Err()
    if err != nil{
        panic(err.Error())
    }

    temp,err := db.Prepare("select process_name from websocket where id = ? limit 1")
    if err != nil{
         panic(err.Error())

    }
    defer temp.Close()
    nrows,err :=  temp.Query(1)
    if err != nil{
        panic(err.Error())
    }
    defer rows.Close()

    var(
        pro_name string
    )
    for nrows.Next(){
        err := nrows.Scan(&pro_name)
        if err != nil{
             panic(err.Error())
        }
        fmt.Println(pro_name)
    }

    var pro_id int

    err = db.QueryRow("select process_id from websocket where id = ? limit 1",1).Scan(&pro_id)
    if err != nil{
         panic(err.Error())
    }
    fmt.Println(pro_id)

    var add_t string

    stmt,err := db.Prepare("select add_time from websocket where id = ?")
    if err != nil{
        panic(err.Error())

    }
    err = stmt.QueryRow(1).Scan(&add_t)
    fmt.Println(add_t)

    //exec (insert,delete,update)
    exec_tmp,err := db.Prepare("insert into websocket (process_name,process_id,type)values(?,?,?)")
    if err != nil{
        panic(err.Error())
    }
    var(
        newpro_name string = "load_page"
        newpro_id int = 1000
        newpro_type int = 1
    )
    res,err := exec_tmp.Exec(newpro_name,newpro_id,newpro_type)
    if err != nil{
        panic(err.Error())
    }
    lastid,err := res.LastInsertId()
    affect,err := res.RowsAffected()
    fmt.Printf("lastid = %d,afected rows = %d",lastid,affect)
}




package main
import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
    "time"
)

func checkErr(err error)(bool){
    if err != nil {
        panic (err.Error())
    }
    return true
}


func createData ()  (bool) {
    db, err := sql.Open("mysql", "wustan:websocket@(127.0.0.1:3306)/websocket")
    if err!= nil {
        panic(err.Error())
    }
    defer db.Close()
    for i := 1; i < 10; i++ {
        var (
            enc = strconv.Itoa(i)
            time_from = "2016-03-27 1"+ enc + ":00:00"
            time_to   = "2016-03-29 1"+ enc + ":00:00"
            day = strconv.Itoa(i%7)
            trigger  = "1" + enc + ":00:" + enc + "0"
            callback = "callback_" + enc
            status   = i%2
            tempTime = time.Now()
            timelayout  = "2006-01-02 15:04:05"
            create_time = tempTime.Format(timelayout)
        )

        stmt, err := db.Prepare( "INSERT INTO task(time_from, time_to, day, trigger, callback, status, create_time) VALUES ($1, $2, $3, $4, $5, $6, $7) ")
        if err != nil {
            panic(err.Error())
        }
        res, err := stmt.Exec(
            time_from,
            time_to,
            day,
            trigger,
            callback,
            status,
            create_time,
        )
        if err != nil {
            panic(err.Error())
        }
        id, err := res.LastInsertId()
        if err != nil {
            panic(err.Error())
        }
        fmt.Println(id)

        //_, err := db.Exec(
            //"INSERT INTO task (time_from, time_to, day, trigger,callback, status, create_time) VALUES ($1, $2, $3, $4, $5, $6, $7)",
            //time_from,
            //time_to,
            //day,
            //trigger,
            //callback,
            //status,
            //create_time,
        //)
        //if err != nil {
            //panic(err.Error())
        //}
    }
    return true
}


func main(){
    ret := createData()
    if true == ret {
        fmt.Println("create data success !")
    }

}

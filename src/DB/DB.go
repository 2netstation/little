package main
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
)

type myTask struct {
    daytime string
    trigger_time string
    callback string
}

func createTask(db *sql.DB, ){
    db, err := sql.Open("mysql", "wustan:websocket@(127.0.0.1:3306)/websocket")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()
    pre,err := db.Prepare("INSERT INTO task (time_from, time_to, daytime, trigger_time, callback, status,create_time) VALUES (?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        panic(err.Error())
    }
    res, err := pre.Exec(
        "2016-04-10 00:00:00",
        "2016-05-10 00:00:00",
        "monday",
        "01:01:01",
        "callBackTest",
        1,
        "2016-05-10 00:00:00",
    )
    if err != nil {
        panic(err.Error())
    }
    lastid,err := res.LastInsertId()
    fmt.Println(lastid)

}

func getTask(db *sql.DB, begin string,end string) []myTask {
    defer db.Close()
    rows, err := db.Query("SELECT `daytime`, `trigger_time`, `callback` FROM task WHERE time_from >= ? AND time_to <= ?", begin, end)
    if err != nil {
        panic(err.Error())
    }
    var ret []myTask
    for rows.Next() {
        var temp myTask
        if err := rows.Scan(&temp.daytime, &temp.trigger_time, &temp.callback); err != nil {
            panic(err.Error())
        }
        ret = append(ret, temp)
    }
    return ret
}

func main(){
    begin   := "2015-02-02 00:00:00"
    end     := "2017-02-02 00:00:00"
    db, err := sql.Open("mysql","wustan:websocket@(127.0.0.1:3306)/websocket")
    if err != nil {
        panic(err.Error())
    }
    ret := getTask(db, begin, end)
    fmt.Println(ret[1].daytime)
}

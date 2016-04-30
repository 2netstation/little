package  DB

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    /*"fmt"*/
)

//加载任务结构
type MyTask struct {
    Time_from string
    Time_to string
    Daytime string
    Second_interval int
    Minute_interval int
    Hour_interval int
    Trigger_time string
    Callback string
}

//创建任务结构
type Task struct {
    time_from string
    time_to   string
    daytime   string
    second_interval int
    minute_interval int
    hour_interval int
    trigger_time string
    callback  string
    status    int
    create_time string
}

//创建任务
func CreateTask(db *sql.DB, task Task) int64 {
    pre,err := db.Prepare("INSERT INTO task (time_from, time_to, daytime, second_interval, minute_interval, hour_interval, trigger_time, callback, status,create_time) VALUES (?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        panic(err.Error())
    }
    res, err := pre.Exec(
        task.time_from,
        task.time_to,
        task.daytime,
        task.second_interval,
        task.minute_interval,
        task.hour_interval,
        task.trigger_time,
        task.callback,
        task.status,
        task.create_time,
    )
    if err != nil {
        panic(err.Error())
    }
    lastid,err := res.LastInsertId()
    return lastid
}

//获取任务
func GetTask(db *sql.DB, begin, end, daytime string) []MyTask {
    rows, err := db.Query("SELECT `time_from`, `time_to`, `daytime`, `trigger_time`, `callback`, `second_interval`, `minute_interval`, `hour_interval` FROM task WHERE time_from >= ? AND time_to <= ? AND daytime = ?", begin, end, daytime)
    if err != nil {
        panic(err.Error())
    }
    var ret []MyTask
    for rows.Next() {
        var temp MyTask
        if err := rows.Scan(&temp.Time_from, &temp.Time_to, &temp.Daytime, &temp.Trigger_time, &temp.Callback, &temp.Second_interval, &temp.Minute_interval, &temp.Hour_interval); err != nil {
            panic(err.Error())
        }
        ret = append(ret, temp)
    }
    return ret
}

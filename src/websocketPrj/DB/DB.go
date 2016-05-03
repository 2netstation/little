package DB

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

//Query return task struct
type MyTask struct {
	Time_from       string
	Time_to         string
	Daytime         string
	Second_interval int
	Minute_interval int
	Hour_interval   int
	Trigger_time    string
	Callback        string
}

//Create task struct
type Task struct {
	Time_from       string
	Time_to         string
	Daytime         string
	Second_interval int
	Minute_interval int
	Hour_interval   int
	Trigger_time    string
	Callback        string
	Status          int
	Create_time     string
	Module          string
}

func MakeSqlConn() *sql.DB {
	db, err := sql.Open("mysql", "wustan:websocket@(127.0.0.1:3306)/websocket")
	if err != nil {
		panic(err.Error())
	}
	return db
}

//Update task status
func SetTaskStatus(db *sql.DB, task MyTask, st int) bool {
	pre, err := db.Prepare("UPDATE task SET status = ? WHERE callback = ? and status = 0")
	if err != nil {
		panic(err.Error())
	}
	_, err = pre.Exec(
		st,
		task.Callback,
	)
	if err != nil {
		panic(err.Error())
	}
	return true
}

//Create a new task
func CreateTask(db *sql.DB, task Task) int64 {
	pre, err := db.Prepare("INSERT INTO task (time_from, time_to, daytime, second_interval, minute_interval, hour_interval, trigger_time, callback, status, create_time, module) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	res, err := pre.Exec(
		task.Time_from,
		task.Time_to,
		task.Daytime,
		task.Second_interval,
		task.Minute_interval,
		task.Hour_interval,
		task.Trigger_time,
		task.Callback,
		task.Status,
		task.Create_time,
		task.Module,
	)
	if err != nil {
		panic(err.Error())
	}
	lastid, err := res.LastInsertId()
	return lastid
}

//Get task
func GetTask(db *sql.DB, end, module string) ([]MyTask, int) {
	rows, err := db.Query("SELECT `time_from`, `time_to`, `daytime`, `trigger_time`, `callback`, `second_interval`, `minute_interval`, `hour_interval` FROM task WHERE time_to >= ? AND status = 0 AND module = ?", end, module)
	var temp []MyTask

	if err != nil {
		if err == sql.ErrNoRows {
			// there were no rows, but otherwise no error occurred
			return temp, 0
		} else {
			log.Fatal(err)
		}
	}
	var ret []MyTask
	for rows.Next() {
		var temp MyTask
		if err := rows.Scan(&temp.Time_from, &temp.Time_to, &temp.Daytime, &temp.Trigger_time, &temp.Callback, &temp.Second_interval, &temp.Minute_interval, &temp.Hour_interval); err != nil {
			panic(err.Error())
		}
		ret = append(ret, temp)
	}
	return ret, 1
}

package master

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
	"websocketPrj/DB"
	"websocketPrj/service"

	"code.google.com/p/go.net/websocket"
)

type taskTemp struct {
	Time_from       string
	Time_to         string
	Daytime         string
	Second_interval string
	Minute_interval string
	Hour_interval   string
	Trigger_time    string
	Callback        string
	Status          string
	Create_time     string
	Module          string
}

func worker(wg *sync.WaitGroup, cs chan string, msg string) {
	defer wg.Done()
	cs <- msg
}

func monitorWorker(wg *sync.WaitGroup, cs chan string) {
	wg.Wait()
}

//Decode receive msg, get task
func decodMsg(msg string) DB.Task {

	var temp taskTemp
	var task DB.Task

	json.Unmarshal([]byte(msg), &temp)

	task.Time_from = temp.Time_from
	task.Time_to = temp.Time_to
	task.Daytime = temp.Daytime
	task.Second_interval, _ = strconv.Atoi(temp.Second_interval)
	task.Minute_interval, _ = strconv.Atoi(temp.Minute_interval)
	task.Hour_interval, _ = strconv.Atoi(temp.Hour_interval)
	task.Trigger_time = temp.Trigger_time
	task.Callback = temp.Callback
	task.Status, _ = strconv.Atoi(temp.Status)
	task.Create_time = temp.Create_time
	task.Module = temp.Module
	return task
}

//Check new task effect
func checkNewTask(task DB.Task) bool {
	if 0 == len(task.Time_from) || 0 == len(task.Time_to) || 0 == len(task.Daytime) || 0 == len(task.Callback) {
		return false
	}
	return true
}

func taskTickerRun(ticker *time.Ticker, ws *websocket.Conn, db *sql.DB) {
	for _ = range ticker.C {
		now := time.Now()

		//daytime := now.Weekday().String()
		today := now.Format("2006-01-02")
		//begin := today + " 00:00:00"
		end := today + " 23:59:59"
		tasks, get := DB.GetTask(db, end, "all")
		if get == 1 {
			for _, task := range tasks {
				cb := service.CallbackEvent{task.Callback, ws}
				service.TaskManager(task, &cb)
				DB.SetTaskStatus(db, task, 1)
			}
		}

	}
}

func initTask(db *sql.DB, ws *websocket.Conn, module string) {
	now := time.Now()
	end := now.Format("2006-01-02 15:04:05")
	tasks, get := DB.GetTask(db, end, module)
	if get == 1 {
		for _, task := range tasks {
			//Print loaded task
			fmt.Println(task)
			cb := service.CallbackEvent{task.Callback, ws}
			service.TaskManager(task, &cb)
		}
	}

}

func initCurrentTask(db *sql.DB, ws *websocket.Conn, module string) {
	now := time.Now()
	today := now.Format("2006-01-02")
	end := today + " 23:59:59"
	tasks, get := DB.GetTask(db, end, module)
	if get == 1 {
		for _, task := range tasks {
			cb := service.CallbackEvent{task.Callback, ws}
			service.TaskManager(task, &cb)
			//DB.SetTaskStatus(db, task, 1)
		}
	}

}

func MasterHandler(ws *websocket.Conn) {
	var err error
	var recvMsg string
	//Get sql connection
	db := DB.MakeSqlConn()
	defer db.Close()

	wg := &sync.WaitGroup{}
	ch := make(chan string)

	current := time.Now()
	todayEnd := current.Format("2006-01-02") + " 23:59:59"
	t := service.FormatTime(todayEnd)

	dur := t.Sub(current) + time.Second
	timer := time.NewTimer(dur)
	//Init exists tasks
	initTask(db, ws, "all")
	for {
		if err = websocket.Message.Receive(ws, &recvMsg); err != nil {
			fmt.Print("Cant't receive!\n")
		}

		wg.Add(1)
		go worker(wg, ch, recvMsg)
		go monitorWorker(wg, ch)

		select {
		case rawTask := <-ch:
			//Get a new task
			task := decodMsg(rawTask)
			//Check new task
			if false == checkNewTask(task) {
			} else {
				//Save new task to database
				taskId := DB.CreateTask(db, task)
				//Init new tasks
				msg := "Create task succed! task id = " + strconv.FormatInt(taskId, 10)
				if err := websocket.Message.Send(ws, msg); err != nil {
					fmt.Println("can't send! \n")
				}
				initTask(db, ws, task.Module)
			}
		case <-timer.C:
			taskTicker := time.NewTicker(time.Second * 86400)
			go taskTickerRun(taskTicker, ws, db)
		}
	}

}

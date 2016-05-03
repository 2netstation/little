package service

import (
	"fmt"
	"time"
	"websocketPrj/DB"

	"code.google.com/p/go.net/websocket"
)

const REGULAR_STOP_SIGNAL int = 1
const REGULAR_START_SIGNAL int = 0

type CallbackEvent struct {
	Name string
	Ws   *websocket.Conn
}

//Set an regular clock event
func SetRegularClock(
	mytask DB.MyTask,
	cb *CallbackEvent,
	interval int) string {

	dur := time.Duration(interval) * time.Second

	//current := FormatTime("2016-05-01 12:50:00")
	current := time.Now()
	start := FormatTime(mytask.Trigger_time)
	stop := FormatTime(mytask.Time_to)

	start_dur := start.Sub(current)
	stop_dur := stop.Sub(current)

	//get start and stop timer
	s := time.NewTimer(start_dur)
	e := time.NewTimer(stop_dur)
	ticker := time.NewTicker(dur)

	for sig := REGULAR_START_SIGNAL; sig != REGULAR_STOP_SIGNAL; {
		select {
		case <-s.C:
			fmt.Println("Start :=======> " + mytask.Callback)
			go tickerRun(ticker, cb)
		case <-e.C:
			sig = REGULAR_STOP_SIGNAL
			ticker.Stop()
			fmt.Println("Stop :=======> " + mytask.Callback)
			db := DB.MakeSqlConn()
			defer db.Close()
			DB.SetTaskStatus(db, mytask, 1)
			break
		}
	}
	return mytask.Callback
}

//Set an Irregular clock event
func SetSingleClock(mytask DB.MyTask, cb *CallbackEvent) string {
	current := time.Now()

	trigger := FormatTime(mytask.Trigger_time)
	trigger_dur := trigger.Sub(current)
	t := time.NewTimer(trigger_dur)
	for sig := REGULAR_START_SIGNAL; sig != REGULAR_STOP_SIGNAL; {
		select {
		case <-t.C:
			fmt.Println("Do :=======> " + mytask.Callback)
			go runCallback(cb)
			sig = REGULAR_STOP_SIGNAL
			db := DB.MakeSqlConn()
			defer db.Close()
			DB.SetTaskStatus(db, mytask, 1)
			break
		}
	}
	return mytask.Callback
}

//Task manager
func TaskManager(mytask DB.MyTask, cb *CallbackEvent) bool {
	current := time.Now()
	weekday := current.Weekday().String()

	interval := mytask.Second_interval
	interval += mytask.Minute_interval * 60
	interval += mytask.Hour_interval * 60 * 60
	//Set regular clock
	if 0 == interval && weekday == mytask.Daytime {
		go SetSingleClock(mytask, cb)
		return true
	}
	//Set single clock
	if 0 != interval {
		go SetRegularClock(mytask, cb, interval)
		return true
	}
	return true
}

//Formate input time string to right CST time
func FormatTime(source_time string) time.Time {
	time_format := "2006-01-02 15:04:05 (MST)"
	source_time += " (CST)"
	res_time, _ := time.Parse(time_format, source_time)
	return res_time
}

//Run a ticker
func tickerRun(ticker *time.Ticker, cb *CallbackEvent) {
	for _ = range ticker.C {
		if err := websocket.Message.Send(cb.Ws, cb.Name); err != nil {
			fmt.Println("can't send! \n")
		}
	}
}

//Run call back
func runCallback(cb *CallbackEvent) {
	if err := websocket.Message.Send(cb.Ws, cb.Name); err != nil {
		fmt.Println("can't send! \n")
	}
}

package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
	"websocketPrj/DB"

	"code.google.com/p/go.net/websocket"
)

func worker(wg *sync.WaitGroup, cs chan string, msg string) {
	defer wg.Done()
	cs <- "worker: " + msg
}

func monitorWorker(wg *sync.WaitGroup, cs chan string) {
	wg.Wait()
	close(cs)
}

func goPrint(msg *string) string {
	*msg += "goPrint"
	return *msg
}

func print(msg string) {
	go goPrint(&msg)
	fmt.Println(msg)
}

func recursiveCall(product int, num int, ch chan int) {
	product += num

	if num == 1 {
		ch <- product
		return
	}
	fmt.Println(product)
	go recursiveCall(product, num-1, ch)
}

func testDB(db *sql.DB) int64 {
	task := DB.Task{
		"2012-02-02 12:00:00",
		"2012-02-02 12:00:00",
		"Tuesday",
		1,
		0,
		0,
		"2012-02-02 12:00:00",
		"method1",
		0,
		"2012-02-02 12:00:00",
		"all",
	}

	id := DB.CreateTask(db, task)
	return id
}

func testString(s string) int {
	return len(s)
}

func main() {

	//db, err := sql.Open("mysql", "wustan:websocket@(127.0.0.1:3306)/websocket")
	//if err != nil {
	//panic(err.Error())
	//}
	//id := testDB(db)
	//fmt.Println(id)

	//fmt.Println(testString(""))

	time := time.Now()
	fmt.Println(time.Format("2006-01-02 15:04:05"))
}

const (
	Sunday time.Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func Echo(ws *websocket.Conn) {
	var err error
	var recvMsg string
	/*ticker := time.NewTicker(time.Second * 1)*/

	wg := &sync.WaitGroup{}
	ch := make(chan string)
	for {

		if err = websocket.Message.Receive(ws, &recvMsg); err != nil {
			fmt.Print("Cant't receive!\n")
		}
		wg.Add(1)
		go worker(wg, ch, recvMsg)
		go monitorWorker(wg, ch)

		select {
		case task := <-ch:
			msg := "Send to browser: " + task
			fmt.Print("Received back from client: " + task + "\n")
			t := time.Now()
			time := t.Format("2006-01-02 15:04:05")
			msg += time + "\n"
			if err := websocket.Message.Send(ws, msg); err != nil {
				fmt.Println("can't send! \n")
			}

		}

	}
}

//规则定时
func regularClock(ws *websocket.Conn) {
	ticker := time.NewTicker(time.Second * 10)
	var reply string

	if err := websocket.Message.Receive(ws, &reply); err != nil {
		fmt.Print("Cant't receive!\n")
	}
	fmt.Println("Received unix time : " + "\n")
	fmt.Println(reply)

	timelayout := "2006-01-02 15:04:05"
	repTime, _ := time.Parse(timelayout, reply)
	//Unix 时间戳转为时间
	setTime := time.Unix(repTime.Unix(), 0)
	fmt.Println("Received from client : " + "\n")
	fmt.Println(setTime)
	fmt.Println("Start send message!")
	for _ = range ticker.C {
		t := time.Now()
		time := t.Format("2006-01-02 15:04:05")
		msg := "Send time: " + time + "\n"
		if err := websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("can't send! \n")
		}
	}
}

//设置任务
/*func SetTask(mytask task.MyTask) bool {*/
/*daytime := mytask.Daytime*/
/*trigger := mytask.Trigger_time*/
/*callback := mytask.Callback*/

/*now := time.Now()*/
/*if Wednesday == now.Weekday() {*/
/*fmt.Println(true)*/
/*}*/
/*return true*/
/*}*/

//确定性一次定时
func singleClock(ws *websocket.Conn) {
	var setString string
	if err := websocket.Message.Receive(ws, &setString); err != nil {
		fmt.Println("Can't receive set time")
	}
	timelayout := "2006-01-02 15:04:05 (MST)"
	setString += " (CST)"
	time_now := time.Now()
	set_time, _ := time.Parse(timelayout, setString)
	sec := set_time.Sub(time_now)
	fmt.Println(time_now)
	fmt.Println(set_time)
	ticker := time.NewTicker(sec)
	sendmsg := "Time now: " + time_now.Format(timelayout)
	if err := websocket.Message.Send(ws, sendmsg); err != nil {
		fmt.Println("can't send")
	}
	for {
		select {
		case <-ticker.C:
			time := time.Now().Format("2006-01-02 15:04:05")
			msg := "Send time: " + time
			if err := websocket.Message.Send(ws, msg); err != nil {
				fmt.Println("Can't send !")
			}
		}
	}
}

//测试读取数据
func readFrom(ws *websocket.Conn) {
	var reply string
	if err := websocket.Message.Receive(ws, &reply); err != nil {
		fmt.Println("can't read ")
	}
	fmt.Println(reply)
}

/* test service/regularClock
func main() {
	task := DB.MyTask{
		"2016-05-01 12:50:05",
		"2016-05-01 12:50:10",
		"Satyrday",
		1,
		0,
		0,
		"2016-05-1 09:14:00",
		"This callback method",
	}
	callback := CallbackEvent{
		"my test callback",
	}
	cbPoint := &callback
	ret := SetRegularClock(task, cbPoint)
	fmt.Println(ret)
}
*/

/*
func main() {
	task := DB.MyTask{
		"2016-05-01 21:15:10",
		"2016-05-01 21:16:00",
		"Sunday",
		0,
		0,
		0,
		"2016-05-01 23:42:00",
		"This callback method",
	}
	callback := CallbackEvent{
		"my test callback",
	}
	cbPoint := &callback
	//ret := SetRegularClock(task, cbPoint)
	ret := SetSingleClock(task, cbPoint)
	fmt.Println(ret)
}*/

package main

import (
    "code.google.com/p/go.net/websocket"
    "fmt"
    "log"
    "net/http"
    "time"
    "database/sql"
    task "DB"
    "reflect"
)

const (
   Sunday time.Weekday = iota
   Monday
   Tuesday
   Wednesday
   Thursday
   Friday
   Saturday
)

func Echo(ws *websocket.Conn){
    var err error
    var reply string
    ticker := time.NewTicker(time.Second * 10)
    for _ = range ticker.C {
        if err = websocket.Message.Receive(ws,&reply);err != nil{
            fmt.Print("Cant't receive!\n")
        }
        msg := "Send to browser: " + reply
        fmt.Print("Received back from client: " + reply+"\n")
        if err := websocket.Message.Send(ws, msg); err != nil {
            fmt.Println("Can't send!\n")
        }
        t := time.Now()
        time := t.Format("2006-01-02 15:04:05")
        msg  = time + "\n"
        if err := websocket.Message.Send(ws,msg); err != nil {
           fmt.Println("can't send! \n")
       }
    }
}

//规则定时
func regularClock(ws *websocket.Conn){
    ticker := time.NewTicker(time.Second * 10)
    var reply string

    if err := websocket.Message.Receive(ws, &reply); err != nil {
        fmt.Print("Cant't receive!\n")
    }
    fmt.Println("Received unix time : " + "\n")
    fmt.Println(reply)

    timelayout := "2006-01-02 15:04:05"
    repTime,_  := time.Parse(timelayout, reply)
    //Unix 时间戳转为时间
    setTime := time.Unix(repTime.Unix(), 0)
    fmt.Println("Received from client : " + "\n")
    fmt.Println(setTime)
    fmt.Println("Start send message!")
    for _ = range ticker.C {
        t := time.Now()
        time := t.Format("2006-01-02 15:04:05")
        msg  := "Send time: "+time+"\n"
        if err := websocket.Message.Send(ws, msg); err != nil {
            fmt.Println("can't send! \n")
        }
    }
}

//设置任务
func SetTask(mytask task.MyTask) bool {
    daytime  := mytask.Daytime
    trigger  := mytask.Trigger_time
    callback := mytask.Callback

    now := time.Now()
    if Wednesday == now.Weekday() {
        fmt.Println(true)
    }
    return true
}

//确定性一次定时
func singleClock(ws *websocket.Conn){
    var setString string
    if err := websocket.Message.Receive(ws, &setString); err != nil {
        fmt.Println("Can't receive set time")
    }
    timelayout := "2006-01-02 15:04:05 (MST)"
    setString  += " (CST)"
    time_now   := time.Now()
    set_time,_ := time.Parse(timelayout, setString)
    sec := set_time.Sub(time_now)
    fmt.Println(time_now)
    fmt.Println(set_time)
    ticker  := time.NewTicker(sec)
    sendmsg := "Time now: "+ time_now.Format(timelayout)
    if err  := websocket.Message.Send(ws, sendmsg); err != nil {
        fmt.Println("can't send")
    }
    for {
        select {
        case <-ticker.C:
            time   := time.Now().Format("2006-01-02 15:04:05")
            msg    := "Send time: " + time
            if err := websocket.Message.Send(ws, msg); err != nil {
                fmt.Println("Can't send !")
            }
        }
    }
}

//测试读取数据
func readFrom(ws *websocket.Conn){
    var reply string
    if err := websocket.Message.Receive(ws, &reply); err!=nil{
        fmt.Println("can't read ")
    }
    fmt.Println( reply )
}

func main(){
    begin   := "2015-02-02 00:00:00"
    end     := "2017-02-02 00:00:00"
    daytime := Monday
    db, err := sql.Open("mysql","wustan:websocket@(127.0.0.1:3306)/websocket")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()
    ret := task.GetTask(db, begin, end, daytime)
    fmt.Println(reflect.TypeOf(ret[0]))
    fmt.Println(ret[0].Daytime)
    SetTask(ret[0])

    fmt.Println("start ")
    http.Handle("/",http.FileServer(http.Dir(".")))
    /*http.Handle("/socket",websocket.Handler(singleClock))*/
    /*http.Handle("/socket",websocket.Handler(regularClock))*/
    /*http.Handle("/socket",websocket.Handler(Echo))*/
    if err:=http.ListenAndServe(":1234",nil);err!=nil{
         log.Fatal("ListenAndServe",err)
    }
    fmt.Println("end")
}

package main

import (
    "code.google.com/p/go.net/websocket"
    "fmt"
    "log"
    "net/http"
    "time"
)

type clockModule struct {
    day      [7]string
    dateFrom string
    dateTo   string
    trigger  string//触发时间
    status   int

}

type message struct {
    script     string
    callBack   string
    statusCode int
    errno      int
    errmsg     string
}



type myTimer struct {
    script    string
    callBack  string
    interval  int
    beginTime time.Time
    endTime   time.Time
    trigger   time.Time
}


//func Echo(ws *websocket.Conn){
    //var err error
    //var reply string
    //ticker := time.NewTicker(time.Second * 10)
    //for _ = range ticker.C{
       // if err = websocket.Message.Receive(ws,&reply);err != nil{
       //     fmt.Print("Cant't receive!\n")
       // }
       // msg:= "Send to browser"+ reply
       // fmt.Print("Received back from client: "+ reply+"\n")
       // if err:= websocket.Message.Send(ws,msg);err!=nil{
       //     fmt.Println("Can't send!\n")
       // }
        //t := time.Now()
        //time := t.Format("2006-01-02 15:04:05")
        //msg:="Send to browser :"+time+"\n"
        //if err:=websocket.Message.Send(ws,msg);err!=nil{
            //fmt.Println("can't send! \n")
       //}
    //}
//}

//规则定时
func regularClock(ws *websocket.Conn){
    ticker := time.NewTicker(time.Second * 10)
    var reply string

    if err := websocket.Message.Receive(ws, &reply);err!=nil{
        fmt.Print("Cant't receive!\n")
    }
    fmt.Println("Received unix time : "  + "\n")
    fmt.Println(reply)

    timelayout := "2006-01-02 15:04:05"
    repTime,_ := time.Parse(timelayout,reply)
    //Unix 时间戳转为时间
    setTime := time.Unix(repTime.Unix(),0)
    fmt.Println("Received from client : " + "\n")
    fmt.Println(setTime)
    fmt.Println("Start send message!")
    for _ = range ticker.C{
        t := time.Now()
        time := t.Format("2006-01-02 15:04:05")
        msg := "Send time: "+time+"\n"
        if err:=websocket.Message.Send(ws, msg);err!=nil{
            fmt.Println("can't send! \n")
        }
    }
}

//确定性一次定时
func singleClock(ws *websocket.Conn){
    var setString string
    if err := websocket.Message.Receive(ws, &setString);err != nil{
        fmt.Println("Can't receive set time")
    }
    timelayout := "2006-01-02 15:04:05 (MST)"
    setString  += " (CST)"
    //time_now,_ := time.Parse(timelayout, "2016-03-20 10:59:49")
    time_now   := time.Now()
    set_time,_ := time.Parse(timelayout, setString)
    sec := set_time.Sub(time_now)
    fmt.Println(time_now)
    fmt.Println(set_time)
    ticker  := time.NewTicker(sec)
    sendmsg := "Time now: "+ time_now.Format(timelayout)
    if err  := websocket.Message.Send(ws, sendmsg);err!=nil{
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
     fmt.Println("start ")
     http.Handle("/",http.FileServer(http.Dir(".")))
     http.Handle("/socket",websocket.Handler(singleClock))

     if err:=http.ListenAndServe(":1234",nil);err!=nil{
          log.Fatal("ListenAndServe",err)
     }
     fmt.Println("end")
}

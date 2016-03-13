package main

import (
    "code.google.com/p/go.net/websocket"
    "fmt"
    "log"
    "net/http"

)

func Echo(ws *websocket.Conn){
    var err error
    var reply string
    if err = websocket.Message.Receive(ws,&reply);err != nil{
        fmt.Print("Cant't receive!\n")
    }
    msg:= "Send to browser"+ reply
    fmt.Print("Received back from client: "+ reply+"\n")
    if err:= websocket.Message.Send(ws,msg);err!=nil{
        fmt.Println("Can't send!\n")
    }
}


func main(){
     fmt.Println("start ")
     http.Handle("/",http.FileServer(http.Dir(".")))
     http.Handle("/socket",websocket.Handler(Echo))

     if err:=http.ListenAndServe(":1234",nil);err!=nil{
          log.Fatal("ListenAndServe",err)

     }
     fmt.Println("end")
}

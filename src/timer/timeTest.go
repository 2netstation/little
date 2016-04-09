package main

import "fmt"
import "time"

func main(){
    //ticker := time.NewTicker(time.Second *10)

    //for _ = range ticker.C{
        //fmt.Printf("ticked at %v\n",time.Now())
    //}
    var t time.Time
    t = time.Now()
    var hour int
    var min int
    var sec int
    hour,min,sec = t.Clock()
    unix_t:= time.Now().Unix()
    unix_nano := time.Now().UnixNano()
    fmt.Println(unix_t, unix_nano, hour, min, sec)
    fmt.Println(time.Unix(1458448008, 0))


}

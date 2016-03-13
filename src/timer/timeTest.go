package main

import "fmt"
import "time"

func main(){
    ticker := time.NewTicker(time.Second *10)

    for _ = range ticker.C{
        fmt.Printf("ticked at %v\n",time.Now())
    }

}

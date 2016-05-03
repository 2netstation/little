package main

import "fmt"
import "time"

func FormatTime(source_time string) time.Time {
	time_format := "2006-01-02 15:04:05 (MST)"
	source_time += " (CST)"
	res_time, _ := time.Parse(time_format, source_time)
	return res_time
}

func pr(ticker *time.Ticker) {
	for _ = range ticker.C {
		fmt.Println(time.Now())
	}
}

func regular() string {
	current := FormatTime("2016-05-01 12:50:00")
	start := FormatTime("2016-05-01 12:50:05")
	end := FormatTime("2016-05-01 12:50:10")

	start_dur := start.Sub(current)
	end_dur := end.Sub(current)
	s := time.NewTimer(start_dur)
	e := time.NewTimer(end_dur)

	//<-s.C
	//fmt.Println("s timer")
	//<-e.C
	//fmt.Println("e timer")
	ticker := time.NewTicker(time.Second * 1)

	for sig := 0; sig != 1; {
		select {
		case <-s.C:
			fmt.Println("s timer !")
			go pr(ticker)
		case <-e.C:
			sig = 1
			ticker.Stop()
			fmt.Println("e timer!")
			break
		}
	}
	return "end func !"
}

func main() {
	//ticker := time.NewTicker(time.Second *10)

	//for _ = range ticker.C{
	//fmt.Printf("ticked at %v\n",time.Now())
	//}
	var t time.Time
	t = time.Now()
	var hour int
	var min int
	var sec int
	wek := t.Weekday()
	hour, min, sec = t.Clock()
	unix_t := time.Now().Unix()
	fmt.Println(unix_t, hour, min, sec)
	fmt.Println(time.Unix(1458448008, 0))
	fmt.Println(wek)
	//ret := regular()
	//fmt.Println(ret)
}

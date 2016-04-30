package task

import (
    "DB"
    "time"
)

//设置规则性定时
func SetRegularClock(mytask DB.MyTask)(*time.Ticker, time.Duration, time.Duration) {
    //time_format := "2006-01-02 15:04:05 (MST)"
    set_string  := "201604-023 15:04:05"
    set_string   += " (CST)"

    interval := mytask.Second_interval
    interval += mytask.Minute_interval * 60
    interval += mytask.Hour_interval * 60 * 60
    dur := time.Second * time.Duration(interval)
    ticker   := time.NewTicker(dur)

    current    := time.Now()
    start      := mytask.Time_from
    start_unix := FormatTime2Unix(start + " (CST)")
    start_duration := start_unix.Sub(current)

    stop       := mytask.Time_to
    stop_unix  := FormatTime2Unix(stop + " (CST)")
    stop_duration := stop_unix.Sub(current)

    return ticker, start_duration, stop_duration
}

func SetIrregularClock(mytask DB.MyTask)(*time.Ticker, time.Duration, time.Duration) {
    set_string := "201604-023 15:04:05"
    set_string   += " (CST)"


}


//停止打点器
func stopTask(ticker *time.Ticker, stop_duration time.Duration) {
    stop_timer := time.NewTimer(time.Second * stop_duration)
    <-stop_timer.C
    ticker.Stop()
}


//func (t *time.Ticker) StopClock (ticker *time.Ticker, duration time.Duration) {
    //current := time.Now()
//}


func FormatTime2Unix(source_time string) time.Time{
    time_format := "2006-01-02 15:04:05 (MST)"
    res_time,_ := time.Parse(time_format, source_time)
    return res_time
}

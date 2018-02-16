package main

import (
	"fmt"
	"os/exec"
	"time"
)

var (
	currentIteration = 1
)

const (
	workTime             = 25
	breakTime            = 10
	longBreakTime        = 30
	workTimeDisplay      = "display notification \"Starting work time\" with title \"Pomodoro Notification\""
	breakTimeDisplay     = "display notification \"Starting break time\" with title \"Pomodoro Notification\""
	longBreakTimeDisplay = "display notification \"Starting long break time\" with title \"Pomodoro Notification\""
	pomodoroCycle        = 5
)

func pomodoroWorkStart(workChan chan bool) {
	informStartWork()
	time.Sleep(time.Minute * workTime)
	workChan <- true
}

func pomodoroBreakTimeStart(breakChan chan bool) {
	informBreakTime()
	time.Sleep(time.Minute * breakTime)
	breakChan <- true
}

func pomodoroLongBreakTimeStart(longBreakChan chan bool) {
	informLongBreakTime()
	time.Sleep(time.Minute * longBreakTime)
	longBreakChan <- true
}

func informStartWork() {
	exec.Command("say", "Starting Working").Output()
	exec.Command("osascript", "-e", workTimeDisplay).Output()
}

func informBreakTime() {
	exec.Command("say", "Start Break").Output()
	exec.Command("osascript", "-e", breakTimeDisplay).Output()
}

func informLongBreakTime() {
	exec.Command("say", "Start Longer Break").Output()
	exec.Command("osascript", "-e", breakTimeDisplay).Output()
}

func askAnotherSession() string {
	fmt.Println("Want another pomodoro session (Y/N)")
	var input string
	fmt.Scanln(&input)
	return input
}

func pomodoroService(workChan chan bool, breakChan chan bool, longBreakChan chan bool, done chan bool) {
	for {
		select {

		case endWorkTime := <-workChan:
			_ = endWorkTime
			if currentIteration >= pomodoroCycle {
				go pomodoroLongBreakTimeStart(longBreakChan)
				currentIteration = 1
			} else {
				currentIteration++
				go pomodoroBreakTimeStart(breakChan)
			}

		case endSmallBreak := <-breakChan:
			_ = endSmallBreak
			go pomodoroWorkStart(workChan)

		case endLongBreak := <-longBreakChan:
			_ = endLongBreak
			c := askAnotherSession()
			for c != "Y" && c != "N" {
				c = askAnotherSession()
			}
			if c == "Y" {
				go pomodoroWorkStart(workChan)
			} else {
				done <- true
			}
		}
	}
}

func main() {
	fmt.Println("Starting Pomodoro")

	workChan := make(chan bool)
	breakChan := make(chan bool)
	longBreakChan := make(chan bool)
	done := make(chan bool)

	go pomodoroWorkStart(workChan)
	go pomodoroService(workChan, breakChan, longBreakChan, done)

	<-done

}

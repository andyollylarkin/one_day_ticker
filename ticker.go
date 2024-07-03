package onedayticker

import "time"

type OneDayTicker struct {
	stop chan struct{}
}

func NewOneDayTicker() *OneDayTicker {
	return &OneDayTicker{
		stop: make(chan struct{}, 0),
	}
}

func (t *OneDayTicker) Stop() {
	close(t.stop)
}

func (t *OneDayTicker) Ticker(tickHours, tickMinutes int, excludeDays []time.Weekday) <-chan time.Time {
	ticker := time.NewTicker(time.Second * 3)

	outTicker := make(chan time.Time)

	var tickDay int
	var todayAlreadyTick bool

	go func() {
		for {
			now := time.Now()
			nowDay := now.Day()

			select {
			case <-t.stop:
				return
			case <-ticker.C:
				if nowDay > tickDay {
					todayAlreadyTick = false
				}
			}

			if todayAlreadyTick {
				continue
			}

			if inExcludeDays(now.Weekday(), excludeDays) {
				continue
			}

			if now.Hour() < tickHours || now.Minute() < tickMinutes {
				continue
			}

			outTicker <- now
			todayAlreadyTick = true
			tickDay = nowDay
		}
	}()

	return outTicker
}

func inExcludeDays(day time.Weekday, excludeDays []time.Weekday) bool {
	for _, d := range excludeDays {
		if d == day {
			return true
		}
	}

	return false
}

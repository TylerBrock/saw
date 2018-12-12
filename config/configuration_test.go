package config

import (
  "testing"
  "time"
)

var testTimeNow = time.Date(2018, 12, 01, 15, 30, 0, 0, time.UTC)

// test the relative parsing functionality of getTime
func TestGetTimeRelative(t *testing.T) {
  relativeTest1, err := getTime("-2h", testTimeNow)

  if err != nil || relativeTest1.Hour() != 13 {
    t.Errorf("Expected relativeTest1 to set time back by two hours. Error: %s", err)
  }

  relativeTest2, err := getTime("-1h15m", testTimeNow)
  if err != nil || !(relativeTest2.Hour() == 14 && relativeTest2.Minute() == 15) {
    t.Errorf("Expected relativeTest2 to set time back by 1 hour and 15 minutes. Error: %s", err)
  }
}

// test the absolute RFC3339 parsing functionality of getTime
func TestGetTimeAbsoluteRFC3339(t *testing.T) {
  absoluteTime1, err := getTime("2006-01-02T15:04:05+07:00", testTimeNow)

  if err != nil && absoluteTime1.Hour() == 15 && absoluteTime1.Minute() == 4 {
    t.Errorf("Expected to parse absolute time. Error: %s", err)
  }
}

func TestGetTimeAbsoluteSimpleDate(t *testing.T) {
  absoluteTime1, err := getTime("2018-06-26", testTimeNow)

  if err != nil && absoluteTime1.Year() == 2018 && absoluteTime1.Month() == 6 && absoluteTime1.Day() == 26 && absoluteTime1.Hour() == 0 && absoluteTime1.Minute() == 0 && absoluteTime1.Second() == 0  {
    t.Errorf("Expected to parse absolute time from simple string. Error: %s", err)
  }
}

func TestGetTimeAbsoluteSimpleDateAndTime(t *testing.T) {
  absoluteTime1, err := getTime("2018-06-26 12:43:30", testTimeNow)

  if err != nil && absoluteTime1.Year() == 2018 && absoluteTime1.Month() == 6 && absoluteTime1.Day() == 26 && absoluteTime1.Hour() == 12 && absoluteTime1.Minute() == 43 && absoluteTime1.Second() == 30  {
    t.Errorf("Expected to parse absolute time from simple string. Error: %s", err)
  }
}

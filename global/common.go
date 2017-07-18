package global

import (
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func RequestID() string {
	base := strconv.FormatInt(time.Now().UnixNano(), 10)

	return strconv.Itoa(Conf.Api.ServerID) + base
}

func LogExtra(c echo.Context) (time.Time, string, bool) {
	ts := time.Now()

	rid := RequestID()
	var debugOn bool
	// 获取debug选项
	if c.FormValue("log_debug") == "on" {
		debugOn = true
	} else {
		debugOn = false
	}

	return ts, rid, debugOn
}

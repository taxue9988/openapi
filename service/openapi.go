package service

import "github.com/labstack/echo"

func Start() {
	initConfig()
	initLogger(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug)
	initMysql()

	loadApis()
	e := echo.New()
	e.Any("/*", apiRoute)
	e.Logger.Fatal(e.Start(":1323"))
}

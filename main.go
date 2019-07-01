package main

import (
	"dump/cache"
	"dump/route"
	"fmt"
	"net/http"
	"time"
	"util/think"
	"util/thinkFile"
)

func main() {
	/***********************************************************/
	port := ":9094"
	projectName := "sql_dump"
	fmt.Println(time.Now().String()[:19], projectName, port, "运行目录:", thinkFile.GetAbsPathWith("./"))
	defer fmt.Println(time.Now().String()[:19], projectName, port, "退出")
	/***********************************************************/
	//thinkLog.SetLogFileTask()
	cache.InitCache()
	/***********************************************************/
	http.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir("view"))))
	http.HandleFunc("/", route.GetHttpFunc)
	err := http.ListenAndServe(port, nil)
	think.IsNil(err)
}

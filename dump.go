package main

import (
	"bufio"
	"dump/body"
	"dump/input"
	"dump/root"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"util/think"
	"util/thinkFile"
	"util/thinkLog"
)

func main() {
	/***********************************************************/
	port := ":9094"
	projectName := "sql_dump"
	fmt.Println(time.Now().String()[:19], projectName, port, "运行目录:", thinkFile.GetAbsPathWith("./"))
	defer fmt.Println(time.Now().String()[:19], projectName, port, "退出")
	/***********************************************************/
	thinkLog.SetLogFileTask()
	/***********************************************************/
	// 选取要备份的服务器配置
	//messageMap := input.GetMessageFromCommand()
	//path := "./config/localhost.json"
	path := "./config/server.json"
	messageMap := input.GetMessageFromFile(path)
	// 确认
	fmt.Println("input Y to continue:")
	input := bufio.NewScanner(os.Stdin) //初始化一个扫表对象
	for i := 0; i < 1; i++ {
		input.Scan() //扫描输入内容
		str := input.Text()
		flag := strings.Contains(str, "Y") || strings.Contains(str, "y")
		if !flag {
			return
		}
	}
	// 定时任务
	body.MysqlDumpTask(messageMap["nextTime"][0], messageMap["userName"][0], messageMap["password"][0], messageMap["host"][0], messageMap["databases"])
	/***********************************************************/
	root.InitRouter()
	http.HandleFunc("/", root.GetHttpFunc)
	err := http.ListenAndServe(port, nil)
	think.IsNil(err)
}

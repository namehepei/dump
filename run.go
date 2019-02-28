package main

import (
	"bufio"
	"dump/body"
	"dump/input"
	"dump/util"
	"fmt"
	"os"
	"strings"
	"time"
	"util/thinkFile"
)

func main() {
	/***********************************************************/
	projectName := "sql_run"
	fmt.Println(time.Now().String()[:19], projectName, "运行目录:", thinkFile.GetAbsPathWith("./"))
	defer fmt.Println(time.Now().String()[:19], projectName, "退出")
	/***********************************************************/
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
	// 运行SQL写入对应数据库中
	body.RunSql(messageMap["userName"][0], messageMap["password"][0], messageMap["host"][0], messageMap["databases"])
	/***********************************************************/
}

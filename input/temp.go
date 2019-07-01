package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ToDo 命令启动备份方式
func dump() {
	/***********************************************************/
	// 选取要备份的服务器配置
	//messageMap := input.GetMessageFromCommand()
	//path := "./config/localhost.json"
	path := "./config/server.json"
	messageMap := GetMessageFromFile(path)
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
	_ = messageMap
}

func run() {
	//messageMap := input.GetMessageFromCommand()
	path := "./config/localhost.json"
	//path := "./config/server.json"
	messageMap := GetMessageFromFile(path)
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
	_ = messageMap
	//serv.RunSql(messageMap["databases"])

}

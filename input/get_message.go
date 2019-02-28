package input

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"util/think"
)

func GetMessageFromFile(filePath string) map[string][]string {
	file, err := ioutil.ReadFile(filePath)
	think.IsNil(err)

	fileMap := make(map[string]interface{})
	json.Unmarshal(file, &fileMap)
	fmt.Println(fileMap)
	messageMap := make(map[string][]string)
	messageMap["userName"] = []string{fileMap["username"].(string)}
	messageMap["password"] = []string{fileMap["password"].(string)}
	messageMap["host"] = []string{fileMap["host"].(string)}
	messageMap["nextTime"] = []string{fileMap["timer"].(string)}
	for i := 0; i < len(fileMap["databases"].([]interface{})); i++ {
		messageMap["databases"] = append(messageMap["databases"], fileMap["databases"].([]interface{})[i].(string))
	}
	fmt.Println("**********************Check Please**********************")
	for k, v := range messageMap {
		fmt.Printf("%v : %v\n", k, v)
	}
	return messageMap
}

func GetMessageFromCommand() map[string][]string {
	notice := []string{"username", "password", "host", "timer", "databases"}
	message := make([][]string, 0)
	input := bufio.NewScanner(os.Stdin) //初始化一个扫表对象
	for i := 0; i < len(notice); i++ {
		fmt.Print(notice[i] + ":")
		input.Scan() //扫描输入内容
		message = append(message, []string{input.Text()})
	}
	var userName = message[0][0]
	var password = message[1][0]
	var host = message[2][0]
	var nextTime = message[3][0]
	var databases = strings.Split(message[4][0], ",")
	fmt.Println(len(databases))
	// 默认值
	if nextTime == "" {
		panic("no nextTime")
	}
	if userName == "" {
		panic("no userName")
	}
	if password == "" {
		panic("no password")
	}
	if host == "" {
		panic("no host")
	}
	if databases[0] == "" {
		panic("no databases")
	}
	messageMap := make(map[string][]string)
	messageMap["userName"] = []string{userName}
	messageMap["password"] = []string{password}
	messageMap["host"] = []string{host}
	messageMap["nextTime"] = []string{nextTime}
	messageMap["databases"] = databases
	fmt.Println("**********************Check Please**********************")
	for k, v := range messageMap {
		fmt.Printf("%v : %v\n", k, v)
	}
	return messageMap
}

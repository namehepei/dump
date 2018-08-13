package main

import (
	"bufio"
	"dump/body"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"util/think"
	"util/thinkFile"
	"util/thinkHttp"
	"io/ioutil"
		"encoding/json"
	)

func main() {
	/***********************************************************/
	port := ":9094"
	projectName := "dump"
	fmt.Println(time.Now().String()[:19], projectName, port, "运行目录:", thinkFile.GetAbsPathWith("./"))
	defer fmt.Println(time.Now().String()[:19], projectName, port, "退出")
	/***********************************************************/
	//thinkLog.SetLogFileTask()
	/***********************************************************/
	messageMap := getMessageFromCommand()
	body.MysqlDumpTask(messageMap["nextTime"][0], messageMap["userName"][0], messageMap["password"][0], messageMap["host"][0], messageMap["databases"])
	//body.RunSql("","","",nil)
	/***********************************************************/
	initRouter()
	http.HandleFunc("/", getHttpFunc)
	err := http.ListenAndServe(port, nil)
	think.Check(err)
}

func getMessageFromFile() map[string][]string {
	file, err := ioutil.ReadFile("./default.json")
	think.Check(err)

	fileMap := make(map[string]interface{})
	json.Unmarshal(file,&fileMap)
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

func getMessageFromCommand() map[string][]string {
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
	//nextTime := "* * * 1 0 0"
	//userName := "root"
	//password := "mysql"
	//host := "localhost"
	//databases := []string{"test"}
	if nextTime == "" {
		nextTime = "* * * 1 0 0"
	}
	if userName == "" {
		userName = "root"
	}
	if password == "" {
		password = "mysql"
	}
	if host == "" {
		host = "localhost"
	}
	if databases[0] == "" {
		databases = []string{"lima", "danger_game"}
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

func initRouter() {
	thinkHttp.AddHttpFunc("/dump/ok", isAlive)
}

func isAlive(w http.ResponseWriter, r *http.Request) {
	think.GetResponseJsonOK(w, r.Host+",is OK")
}

func getHttpFunc(w http.ResponseWriter, r *http.Request) {
	defer think.DeferRecover(w)
	header := r.Header.Get("superman")
	if header != "superman" {
		think.GetResponseJsonFail(w, 403, "no access", 403)
		return
	}
	url := r.RequestURI
	thinkHttp.RequestMap[url].Function(w, r)
}

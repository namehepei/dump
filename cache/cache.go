package cache

import (
	"io/ioutil"
	"net/http"
	"util/think"
	"util/thinkFile"
	"util/thinkHttp"
	"util/thinkJson"
)

const filePath = "./cache/"
const fileName = "user.json"

var usernameList = make([]string, 0)
var passwordList = make([]string, 0)
var hostList = make([]string, 0)

// 初始化缓冲
func InitCache() {
	// 读取
	bs, err := ioutil.ReadFile(filePath + fileName)
	think.IsNil(err)
	jo := thinkJson.MustGetJsonObject(bs)
	usernameList = jo.MustGetStringList("username")
	passwordList = jo.MustGetStringList("password")
	hostList = jo.MustGetStringList("host")
}

func insertUnique(list []string, value string) []string {
	// 方法1.遍历
	for i := 0; i < len(list); i++ {
		if list[i] == value {
			return list
		}
	}
	list = append(list, value)
	// 方法2.map
	return list
}

func AddCache(host, username, password string) {
	// 非重复写入
	usernameList = insertUnique(usernameList, username)
	passwordList = insertUnique(passwordList, password)
	hostList = insertUnique(hostList, host)

	// 写入json文件
	f := thinkFile.OpenFile(filePath, fileName, 770)
	data := make(map[string][]string)
	data["username"] = usernameList
	data["password"] = passwordList
	data["host"] = hostList
	f.Write(thinkJson.MustMarshal(data))
}

func CacheUser(w http.ResponseWriter, r *http.Request) {
	data := make(map[string][]string)
	data["username"] = usernameList
	data["password"] = passwordList
	data["host"] = hostList
	thinkHttp.WriteJsonOk(w, data)
}

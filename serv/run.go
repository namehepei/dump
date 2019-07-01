package serv

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"util/database"
	"util/think"
	"util/thinkFile"
	"util/thinkHttp"
	"util/thinkJson"
)

func RunPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./view/run/*.html")
	think.IsNil(err)
	tmpl.ExecuteTemplate(w, "root.html", nil)
}

// 还原当天的
func RunDatabases(w http.ResponseWriter, r *http.Request) {
	body := thinkHttp.GetRequestBody(r)
	obj := thinkJson.MustGetJsonObject(body)
	databases := obj.MustGetStringList("databases")

	for i := 0; i < len(databases); i++ {
		RunSqlInDatabase("", databases[i])
	}

	thinkHttp.WriteJsonOk(w, "ok")
}

// 获取备份目录下的一级文件夹名称
func ListFile(w http.ResponseWriter, r *http.Request) {
	body := thinkHttp.GetRequestBody(r)
	obj, err := thinkJson.GetJsonObject(body)

	filePath := dumpPath + "/" + obj.MustGetString("filePath")
	allDir, err := ioutil.ReadDir(filePath)
	think.IsNil(err)
	data := make([]string, 0)
	for i := 0; i < len(allDir); i++ {
		data = append(data, allDir[i].Name())
	}

	thinkHttp.WriteJsonOk(w, data)
}

func RunTables(w http.ResponseWriter, r *http.Request) {
	body := thinkHttp.GetRequestBody(r)
	obj := thinkJson.MustGetJsonObject(body)

	databaseRun := obj.MustGetString("database")
	filePath := dumpPath + "/" + obj.MustGetString("filePath")

	RunSqlInDatabase(filePath, databaseRun)

	thinkHttp.WriteJsonOk(w, "ok")
}

func RunSqlInDatabase(filePath, databaseName string) {
	allFile := make([]string, 0)
	// 还原当日的数据库
	// 目录结构: ./dump/20180702/database/table
	//filePath := dumpPath + databases[i] + "/20180728"
	if filePath == "" {
		filePath = dumpPath + time.Now().Format("20060102") + "/" + databaseName
	}
	filePath = thinkFile.GetAbsPathWith(filePath)

	// 重指向默认数据库
	dsn := userName + ":" + password + "@tcp(" + host + ")/" + databaseName
	db = database.SetConn(dsn)
	// 遍历filePath下的所有文件以及目录,ls .sql 文件
	filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".sql") {
			allFile = append(allFile, path)
			return nil
		} else {
			return nil
		}
		return nil
	})
	for j := 0; j < len(allFile); j++ {
		fileFullName := allFile[j]
		fmt.Println(allFile[j])
		runSqlSingle(fileFullName)
	}
}

func runSqlSingle(fileFullName string) {
	file, err := os.Open(fileFullName)
	defer file.Close()
	think.IsNil(err)
	info, err := file.Stat()
	think.IsNil(err)
	temp := make([]byte, info.Size())
	_, err = file.Read(temp)
	think.IsNil(err)
	tempSep := bytes.Split(temp, []byte{';', '\n'})
	for i := 0; i < len(tempSep); i++ {
		sqlString := string(tempSep[i])
		if len(sqlString) == 0 {
			break
		}
		//fmt.Println(sqlString)
		result, err := database.Exec(sqlString)
		think.IsNil(err)
		affect, err := result.RowsAffected()
		think.IsNil(err)
		fmt.Println("affect", affect)
	}
}

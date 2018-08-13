package body

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"util/database"
	"util/think"
	"util/thinkFile"
)

func RunSql(userName, password, host string, databases []string) {
	sourceName := userName + ":" + password + "@tcp(" + host + ":3306)/test"
	db := database.SetConn(sourceName)
	defer db.Close()

	for i := 0; i < len(databases); i++ {
		allFile := make([]string, 0)
		// 目录结构: ./dump/20180702/database/table
		//filePath := dumpPath + databases[i] + "/20180728"
		filePath := dumpPath + time.Now().Format("20060102") + "/" + databases[i]
		filePath = thinkFile.GetAbsPathWith(filePath)
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
}

func runSqlSingle(fileFullName string) {
	file, err := os.Open(fileFullName)
	defer file.Close()
	think.Check(err)
	info, err := file.Stat()
	think.Check(err)
	temp := make([]byte, info.Size())
	_, err = file.Read(temp)
	think.Check(err)
	tempSep := bytes.Split(temp, []byte{';', '\n'})
	for i := 0; i < len(tempSep); i++ {
		sqlString := string(tempSep[i])
		if len(sqlString) == 0 {
			break
		}
		//fmt.Println(sqlString)
		result, err := database.Idb.Exec(sqlString)
		if err != nil {
			fmt.Println(len(sqlString))
			fmt.Println(sqlString[:200])
			fmt.Println(sqlString[len(sqlString)-200:])
			panic(err)
		}
		think.Check(err)
		affect, err := result.RowsAffected()
		think.Check(err)
		fmt.Println("affect", affect)
	}
}

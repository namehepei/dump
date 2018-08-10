package body

import (
	"bytes"
	"os"
	"runtime/debug"
	"strings"
	"time"
	"util/database"
	"util/think"
	"util/thinkFile"
	"util/thinkLog"
	"util/thinkTimer"
)

var dumpPath = "./dumps/"

// 因为MysqlDump()使用包database.Idb,故会关闭主函数的db,不可以在主函数中存在db的情况下使用
// 默认值
//nextTime := "* * * 1 0 0"
//userName := "root"
//password := "mysql"
//host := "localhost"
//databases := []string{"test"}
func MysqlDumpTask(nextTime, userName, password, host string, databases []string) {
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
	if databases == nil || len(databases) == 0 {
		databases = []string{"test"}
	}
	f := func() {
		for i := 0; i < len(databases); i++ {
			MysqlDump(userName, password, host, databases[i])
		}
	}

	go thinkTimer.TaskNow(nextTime, f)
}

func MysqlDump(userName, password, host, DatabaseName string) {
	// mysql
	sourceName := userName + ":" + password + "@tcp(" + host + ":3306)/" + DatabaseName
	thinkLog.DebugLog.Println("db", sourceName)
	db := database.SetConn(sourceName)
	defer db.Close()

	// 目录(\转义字符)
	filePath := dumpPath + time.Now().Format("20060102") + "/" + DatabaseName
	filePath = thinkFile.GetAbsPathWith(filePath)
	fileFullNameSlice := make([]string, 0)
	tables := getTables(DatabaseName)
	for i := 0; i < len(tables); i++ {
		TableName := tables[i]["TABLE_NAME"]
		//fmt.Println("TableName",TableName)
		tableFullName := DatabaseName + "." + TableName
		InsertString := getInsertValues(TableName)
		DDL := getCreateTable(tableFullName)

		var buffer bytes.Buffer
		buffer.WriteString("USE " + DatabaseName + ";\n")
		buffer.WriteString("DROP TABLE IF EXISTS `" + TableName + "`;\n")
		buffer.WriteString(DDL + ";\n")
		buffer.WriteString("LOCK TABLES `" + TableName + "` WRITE;\n")
		buffer.WriteString(InsertString)
		buffer.WriteString("UNLOCK TABLES;\n")

		comment := buffer.String()
		//fmt.Println(comment)
		fileName := TableName + ".sql"
		fileFullName := createSql(filePath, fileName, comment)
		fileFullNameSlice = append(fileFullNameSlice, fileFullName)
	}
	thinkLog.DebugLog.Println("dump "+DatabaseName+"finish", filePath)
	//for i := 0; i < len(fileFullNameSlice); i++ {
	//	fmt.Println(fileFullNameSlice[i])
	//}
	lsDumpFile(filePath, fileFullNameSlice)
}

//
func lsDumpFile(filePath string, fileFullNameSlice []string) {
	fileName := "dir.txt"
	file := thinkFile.OpenFile(filePath, fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	defer file.Close()
	for i := 0; i < len(fileFullNameSlice); i++ {
		file.WriteString(fileFullNameSlice[i] + "\n")
	}
	thinkLog.DebugLog.Println("dump ls finish", filePath+fileName)
}

// XXX.sql
func createSql(filePath, fileName, comment string) string {
	thinkFile.CreatePath(filePath)
	file := thinkFile.OpenFile(filePath, fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	defer file.Close()
	file.WriteString(comment)
	return file.Name()
}

// show tables
func getTables(databaseName string) []map[string]string {
	sqlString := "SELECT distinct TABLE_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ?"
	_, rows := database.SelectMap(nil, sqlString, databaseName)

	return rows
}

// DDL
func getCreateTable(tableName string) string {
	sqlString := "SHOW CREATE TABLE " + tableName
	_, rows := database.SelectMap(nil, sqlString)
	//fmt.Println(rows)
	return rows[0]["Create Table"]
}

// INSERT VALUES
// 实现根据 VALUES (,,,) 的大小自动切割数据
func getInsertValues(tableName string) string {
	// 获取全部的数据rows [][]string
	sqlString := "SELECT * FROM " + tableName
	_, rows := database.SelectList(nil, sqlString)
	// 无记录
	if len(rows) == 0 {
		return ""
	}
	// max_allowed_packet 略大于512 KB
	const oneMB = 1024 * 512
	panicMsg := "[notDeffer]packet for query is too large.Try adjusting the 'max_allowed_packet' variable on the server."
	resultSlice := make([]bytes.Buffer, 0)
	// [][]string 检验(NULL,'),长度,并转化为 []bytes.Buffer(string)
	for i := 0; i < len(rows); i++ {
		insertSingle := rows[i]
		var buffer bytes.Buffer
		buffer.WriteByte('(')
		for j := 0; j < len(insertSingle); j++ {
			temp := rows[i][j]
			if temp == "" {
				temp = "NULL"
			} else {
				// 在temp中的查看是否含有需要转义(')的字符
				temp = "'" + strings.Replace(temp, "'", "\\'", -1) + "'"
			}
			buffer.WriteString(temp)
			if j == len(insertSingle)-1 {
				buffer.WriteByte(')')
			} else {
				buffer.WriteByte(',')
			}
		}
		if buffer.Len() >= oneMB {
			thinkLog.ErrorLog.Println(string(debug.Stack()))
			thinkLog.ErrorLog.Println(panicMsg)
			panic(panicMsg)
		}
		resultSlice = append(resultSlice, buffer)
	}
	// 切割数据(每条INSERT语句大小不超过 1MB)
	var buffer bytes.Buffer
	var estimatedSize = resultSlice[0].Len()
	var i = 0
	for i < len(resultSlice) {
		//oldSize := buffer.Len()
		buffer.WriteString("INSERT INTO `" + tableName + "` VALUES ")
		for true {
			bufferInner := resultSlice[i]

			buffer.WriteByte('\n')
			_, err := buffer.Write(bufferInner.Bytes())
			think.Check(err)

			i++
			var addSize = 0
			if i != len(resultSlice) {
				addSize = resultSlice[i].Len()
			} else {
				buffer.WriteByte(';')
				buffer.WriteByte('\n')
				break
			}
			estimatedSize += addSize
			if estimatedSize < oneMB {
				buffer.WriteByte(',')
			} else {
				estimatedSize = addSize
				buffer.WriteByte(';')
				buffer.WriteByte('\n')
				break
			}
		}
		//newSize := buffer.Len()
		//fmt.Println(tableName,"size",newSize - oldSize)
	}
	return buffer.String()
}

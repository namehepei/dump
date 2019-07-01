package serv

import (
	"bytes"
	"database/sql"
	"dump/cache"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"text/template"
	"time"
	"util/database"
	"util/think"
	"util/thinkFile"
	"util/thinkHttp"
	"util/thinkJson"
	"util/thinkLog"
)

var db *sql.DB

const dumpPath = "./dumps/"

var userName string
var password string
var host string

func DumpPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./view/dump/*.html")
	think.IsNil(err)
	tmpl.ExecuteTemplate(w, "root.html", nil)
}

// 建立数据库连接
func SetDB(w http.ResponseWriter, r *http.Request) {
	body := thinkHttp.GetRequestBody(r)
	obj := thinkJson.MustGetJsonObject(body)
	userName = obj.MustGetString("username")
	password = obj.MustGetString("password")
	host = obj.MustGetString("host")
	dsn := userName + ":" + password + "@tcp(" + host + ")/test"
	fmt.Println(dsn)
	db = database.SetConn(dsn)
	cache.AddCache(host, userName, password)
	//defer db.Close()
	thinkHttp.WriteJsonOk(w, db.Stats())
}

// 复制几个库的全部表
func DumpDatabases(w http.ResponseWriter, r *http.Request) {
	body := thinkHttp.GetRequestBody(r)
	obj := thinkJson.MustGetJsonObject(body)
	databases := obj.MustGetStringList("databases")

	for i := 0; i < len(databases); i++ {
		MysqlDump(databases[i], database.GetTables(databases[i]))
	}

	thinkHttp.WriteJsonOk(w, "ok")
}

// 复制一个库的几个表
func DumpTables(w http.ResponseWriter, r *http.Request) {
	body := thinkHttp.GetRequestBody(r)
	obj := thinkJson.MustGetJsonObject(body)
	databaseName := obj.MustGetString("database")
	tables := obj.MustGetStringList("tables")

	MysqlDump(databaseName, tables)

	thinkHttp.WriteJsonOk(w, "ok")
}

func MysqlDump(databaseName string, tables []string) {
	// 目录(\转义字符)
	filePath := dumpPath + time.Now().Format("20060102") + "/" + databaseName
	filePath = thinkFile.GetAbsPathWith(filePath)
	fileFullNameSlice := make([]string, 0)

	for i := 0; i < len(tables); i++ {
		tableName := tables[i]
		//fmt.Println("TableName",TableName)
		InsertString := getInsertValues(databaseName, tableName)
		DDL := database.GetDDL(databaseName + "." + tableName)

		var buffer bytes.Buffer

		buffer.WriteString("DROP TABLE IF EXISTS `" + tableName + "`;\n")
		buffer.WriteString(DDL + ";\n")
		buffer.WriteString("LOCK TABLES `" + tableName + "` WRITE;\n")
		buffer.WriteString(InsertString)
		buffer.WriteString("UNLOCK TABLES;\n")

		comment := buffer.String()
		//fmt.Println(comment)
		fileName := tableName + ".sql"
		fileFullName := createSql(filePath, fileName, comment)
		fileFullNameSlice = append(fileFullNameSlice, fileFullName)
	}
	thinkLog.DebugLog.Println("dump "+databaseName+" finish", filePath)
	//for i := 0; i < len(fileFullNameSlice); i++ {
	//	fmt.Println(fileFullNameSlice[i])
	//}
	lsDumpFile(filePath, fileFullNameSlice)
}

func DatabasesInMysql(w http.ResponseWriter, r *http.Request) {
	list := database.GetDatabases()
	thinkHttp.WriteJsonOk(w, list)
}

func TablesInMysql(w http.ResponseWriter, r *http.Request) {
	body := thinkHttp.GetRequestBody(r)
	obj := thinkJson.MustGetJsonObject(body)
	databaseName := obj.MustGetString("database")
	list := database.GetTables(databaseName)
	thinkHttp.WriteJsonOk(w, list)
}

func DdlInMysql(w http.ResponseWriter, r *http.Request) {
	body := thinkHttp.GetRequestBody(r)
	obj := thinkJson.MustGetJsonObject(body)
	databaseName := obj.MustGetString("database")
	tables := obj.MustGetStringList("tables")

	s := ""
	for i := 0; i < len(tables); i++ {
		table := tables[i]
		s += database.GetDDL(databaseName + "." + table)
		s += ";\n\n"
	}

	thinkHttp.WriteJsonOk(w, s)
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

// INSERT VALUES
// 实现根据 VALUES (,,,) 的大小自动切割数据
func getInsertValues(databaseName, tableName string) string {
	// 获取全部的数据rows [][]string
	sqlString := "SELECT * FROM " + databaseName + "." + tableName
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
			think.IsNil(err)

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

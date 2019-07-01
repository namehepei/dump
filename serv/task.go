package serv

import (
	"net/http"
	"text/template"
	"util/database"
	"util/think"
	"util/thinkHttp"
	"util/thinkJson"
	"util/thinkTimer"
)

var ts = new(thinkTimer.ThinkTasks)

func init() {
	ts.Init()
	//ts.AddTask("* * * 0 0 0", )
}

func TaskList(w http.ResponseWriter, r *http.Request) {
	m := ts.GetTasks()
	s := make([]*thinkTimer.ThinkTask, 0)
	for _, value := range m {
		s = append(s, value)
	}
	thinkHttp.WriteJsonOk(w, s)
}

func TaskStart(w http.ResponseWriter, r *http.Request) {
	body := thinkHttp.GetRequestBody(r)
	obj := thinkJson.MustGetJsonObject(body)
	nextTime := obj.MustGetString("nextTime")
	kind := obj.MustGetString("kind")

	var t *thinkTimer.ThinkTask
	switch kind {
	case "some":
		databases := obj.MustGetStringList("databases")
		f := func() {
			for i := 0; i < len(databases); i++ {
				MysqlDump(databases[i], database.GetTables(databases[i]))
			}
		}
		t = ts.AddTaskNow(nextTime, obj, f)
	case "one":
		databaseName := obj.MustGetString("database")
		tables := obj.MustGetStringList("tables")
		f := func() {
			MysqlDump(databaseName, tables)
		}
		t = ts.AddTaskNow(nextTime, obj, f)
	}
	thinkHttp.WriteJsonOk(w, t.TaskGuid)
}

func TaskStop(w http.ResponseWriter, r *http.Request) {
	body := thinkHttp.GetRequestBody(r)
	obj := thinkJson.MustGetJsonObject(body)
	taskGuid := obj.MustGetString("taskGuid")

	ts.StopTask(taskGuid)
}

func TaskPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./view/task/*.html")
	think.IsNil(err)
	tmpl.ExecuteTemplate(w, "root.html", nil)
}

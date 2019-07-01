package route

import (
	"dump/cache"
	"dump/serv"
	"net/http"
	"text/template"
	"util/think"
	"util/thinkHttp"
)

func init() {
	thinkHttp.AddHttpFunc("/", IndexPage)
	thinkHttp.AddHttpFunc("/dump", serv.DumpPage)
	thinkHttp.AddHttpFunc("/run", serv.RunPage)
	thinkHttp.AddHttpFunc("/task", serv.TaskPage)

	thinkHttp.AddHttpFunc("/ok", OkHandle)
	thinkHttp.AddHttpFunc("/cache/user", cache.CacheUser)
	thinkHttp.AddHttpFunc("/conn", serv.SetDB)

	thinkHttp.AddHttpFunc("/database", serv.DatabasesInMysql)
	thinkHttp.AddHttpFunc("/table", serv.TablesInMysql)
	thinkHttp.AddHttpFunc("/ddl", serv.DdlInMysql)

	thinkHttp.AddHttpFunc("/dump/some", serv.DumpDatabases)
	thinkHttp.AddHttpFunc("/dump/one", serv.DumpTables)

	thinkHttp.AddHttpFunc("/run/dir", serv.ListFile)
	thinkHttp.AddHttpFunc("/run/one", serv.RunTables)

	thinkHttp.AddHttpFunc("/task/list", serv.TaskList)
	thinkHttp.AddHttpFunc("/task/start", serv.TaskStart)
	thinkHttp.AddHttpFunc("/task/stop", serv.TaskStop)
}

func OkHandle(w http.ResponseWriter, r *http.Request) {
	thinkHttp.WriteJsonOk(w, r.Host+"sql_dump")
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./view/index/*.html")
	think.IsNil(err)
	tmpl.ExecuteTemplate(w, "root.html", nil)
}

func GetHttpFunc(w http.ResponseWriter, r *http.Request) {
	defer thinkHttp.DeferRecoverHttp(w)
	if check(r) {
		thinkHttp.WriteStatus(w, 403)
	}
	f, ok := thinkHttp.FindHttpFunc(r.RequestURI)
	if !ok {
		thinkHttp.WriteStatus(w, 404)
		return
	}
	f(w, r)
}

func check(r *http.Request) bool {
	if r.Host == "localhost" {
		return true
	}
	if header := r.Header.Get("superman"); header == "superman" {
		return true
	}
	return false
}

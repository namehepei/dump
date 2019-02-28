package root

import (
	"net/http"
	"util/think"
	"util/thinkHttp"
)

func InitRouter() {
	thinkHttp.AddHttpFunc("/dump/ok", isAlive)
}

func isAlive(w http.ResponseWriter, r *http.Request) {
	think.GetResponseJsonOK(w, r.Host+",is OK")
}

func GetHttpFunc(w http.ResponseWriter, r *http.Request) {
	defer think.DeferRecover(w)
	header := r.Header.Get("superman")
	if header != "superman" {
		think.GetResponseJsonFail(w, 403, "no access", 403)
		return
	}
	url := r.RequestURI
	thinkHttp.RequestMap[url].Function(w, r)
}

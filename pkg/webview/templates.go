package webview

import "html/template"

var (
	tmplLogin   *template.Template
	tmplSignup  *template.Template
	tmplProfile *template.Template
)

func initTemplates() {
	tmplLogin = template.Must(template.ParseFiles("web/templates/login/index.html"))
	tmplSignup = template.Must(template.ParseFiles("web/templates/signup/index.html"))
	tmplProfile = template.Must(template.ParseFiles("web/templates/profile/index.html"))
}

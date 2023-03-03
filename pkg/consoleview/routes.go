package consoleview

func (app *CMDApp) initRoutes() {
	app.Route("auth", app.auth)
	app.Route("signup", app.signup)
	app.Route("login", app.login)
	app.Route("profile", app.profile)
	app.Route("editProfile", app.editProfile)
}

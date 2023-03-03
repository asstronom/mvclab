package webview

func (s *Server) initRoutes() {
	s.router.HandleFunc("/", s.home).Methods("GET")
	s.router.HandleFunc("/signup", s.signUp).Methods("GET")
	s.router.HandleFunc("/signup", s.validateSignUp).Methods("POST")
	s.router.HandleFunc("/login", s.login).Methods("GET")
	s.router.HandleFunc("/login", s.validateLogin).Methods("POST")
	s.router.HandleFunc("/profile", s.isAuthorized(s.myProfile)).Methods("GET")
	s.router.HandleFunc("/profile/{id}", s.isAuthorized(s.profile)).Methods("GET")
	s.router.HandleFunc("/profile/{id}", s.isAuthorized(s.editProfile)).Methods("POST")
}

package hdle

import (
	"fmt"
	"net/http"
	"time"

	auth "forum/Authentication"
	db "forum/Database"
	Rt "forum/Routes"
	sec "forum/Security"
)

func Handlers(tabb db.Db) {
	tab := tabb
	staticHandler := http.FileServer(http.Dir("templates"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// Initialisation des paramètres de rate limiting
	windowSize := time.Minute // Fenêtre de temps d'une minute
	maxRequests := 10         // Nombre maximum de requêtes autorisées
	maxLoginRequests := 10    // Nombre maximum de tentatives de connexions autorisées
	checkloginTimeOut := false
	checkotherTimeOut := false

	// Serveur HTTP
	http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Votre code existant pour le routeur HTTP
		w.Header().Set("Strict-Transport-Security", "max-age=3336000; includeSubDomains")

		switch r.URL.Path {
		case "/": //default page
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}

			Rt.Index(w, r, tab)

		case "/create": //create account page
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}
			Rt.CreateAccountPage(w, r, tab)

		case "/auth/google/login": // googleAuth login page
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}

			Rt.HandleGoogleLogin(w, r, tab)

		case "/auth/google/callback": //googleAuth response url
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}

			Rt.HandleCallback(w, r, tab)

		case "/auth/github/login": // githubAuth login page
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}

			Rt.HandleGitHubLogin(w, r, tab)

		case "/auth/github/callback": //githubAuth response url
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}

			Rt.HandleGitHubCallback(w, r, tab)

		case "/login": //login page
			if !sec.LoginLimiterMiddleware(r, windowSize, maxLoginRequests) {
				auth.Snippets(w, 429)
				checkloginTimeOut = true
				return
			}
			if checkloginTimeOut {
				time.Sleep(3 * time.Second)
				checkloginTimeOut = false
			}

			Rt.LoginPage(w, r, tab)

		case "/logout": //logout page
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}
			Rt.LogOutHandler(w, r, tab)

		case "/home": //home page
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}

			Rt.HomeHandler(w, r, tab)

		case "/myprofil/posts": //filtered created post page

			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}

			Rt.Profil(w, r, tab)

		case "/myprofil/favorites": //filtered liked post page
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}
			Rt.Profil_fav(w, r, tab)

		case "/myprofil/comments": //filtered commented post page
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}
			Rt.Profil_comment(w, r, tab)

		case "/filter": //filtered post by categorie page for registered
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}
			Rt.Filter(w, r, tab)

		case "/index": //filtered post by categorie page for non-registered
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
				auth.Snippets(w, 429)
				checkotherTimeOut = true
				return
			}
			if checkotherTimeOut {
				time.Sleep(3 * time.Second)
				checkotherTimeOut = false
			}
			Rt.Indexfilter(w, r, tab)

		default: // page does not exist
			auth.Snippets(w, http.StatusNotFound)
		}
	}))

	fmt.Println("📡----------------------------------------------------📡")
	fmt.Println("|                                                    |")
	fmt.Println("| 🌐 Server has started at \033[32mhttps://localhost:8080\033[0m 🟢  |")
	fmt.Println("|                                                    |")
	fmt.Println("📡----------------------------------------------------📡")

	if errr := http.ListenAndServeTLS(":8080", "Security/server.crt", "Security/server.key", nil); errr != nil {
		fmt.Printf("Erreur de serveur HTTPS : %s\n", errr)
	}

}

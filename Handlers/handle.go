package hdle

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"

	auth "forum/Authentification"
	db "forum/Database"
	Rt "forum/Routes"
)

// newLimiterMiddleware crée un middleware rate limiting
func newLimiterMiddleware(r *http.Request, limiter *rate.Limiter) bool {
	return limiter.Allow()
}

// Handlers configure les gestionnaires pour différents endpoints
func Handlers() {
	// Configurer le gestionnaire pour les fichiers statiques
	staticHandler := http.FileServer(http.Dir("templates"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// Initialiser la base de données
	tab, err := db.Init_db()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Initialisation des rate limiter pour certains endpoints
	limiter := rate.NewLimiter(30, 10)          // Limiteur global
	loginLimiter := rate.NewLimiter(2, 1)     // Limiteur spécifique à /login

	// Utilisation du middleware rate limiting pour chaque endpoint
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		switch r.URL.Path {
		case "/":
			if !newLimiterMiddleware(r, limiter) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			Rt.Index(w, r, tab)
		case "/create":
			if !newLimiterMiddleware(r, limiter) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			Rt.CreateAccountPage(w, r, tab)
		case "/login":
			if !newLimiterMiddleware(r, loginLimiter) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			Rt.LoginPage(w, r, tab)
		case "/logout":
			if !newLimiterMiddleware(r, limiter) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			Rt.LogOutHandler(w, r, tab)
		case "/home":
			if !newLimiterMiddleware(r, limiter) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			Rt.HomeHandler(w, r, tab)
		case "/myprofil/posts":
			if !newLimiterMiddleware(r, limiter) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			Rt.Profil(w, r, tab)
		case "/myprofil/favorites":
			if !newLimiterMiddleware(r, limiter) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			Rt.Profil_fav(w, r, tab)
		case "/myprofil/comments":
			if !newLimiterMiddleware(r, limiter) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			Rt.Profil_comment(w, r, tab)
		case "/filter":
			if !newLimiterMiddleware(r, limiter) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			Rt.Filter(w, r, tab)
		case "/index":
			if !newLimiterMiddleware(r, loginLimiter) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			Rt.Indexfilter(w, r, tab)
		default:
			auth.Snippets(w, http.StatusNotFound)
		}
	}))

	// Lancement du serveur

	fmt.Println("📡----------------------------------------------------📡")
	fmt.Println("|                                                    |")
	fmt.Println("| 🌐 Server has started at \033[32mhttps://localhost\033[0m 🟢  |")
	fmt.Println("|                                                    |")
	fmt.Println("📡----------------------------------------------------📡")
	errr := http.ListenAndServeTLS(":https", "Security/server.crt", "Security/server.key", nil)
	if errr != nil {
		fmt.Printf("Erreur de serveur HTTPS : %s\n", errr)
	}
}

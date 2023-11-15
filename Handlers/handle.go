package hdle

import (
    "fmt"
    "net/http"
    "time"

	sec "forum/Security"
    auth "forum/Authentification"
    db "forum/Database"
    Rt "forum/Routes"
)


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

    // Initialisation des paramètres de rate limiting
    windowSize := time.Minute // Fenêtre de temps d'une minute
    maxRequests := 30         // Nombre maximum de requêtes autorisées
	maxLoginTimeout := 6             // Nombre maximum de tentatives de connexions autorisées
	// Utilisation du middleware rate limiting pour chaque endpoint
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		switch r.URL.Path {
		case "/":
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
                http.Error(w, "Limite de taux dépassée", http.StatusTooManyRequests)
                return
            }
			Rt.Index(w, r, tab)
		case "/create":
		 if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
                http.Error(w, "Limite de taux dépassée", http.StatusTooManyRequests)
                return
            }
			Rt.CreateAccountPage(w, r, tab)
		case "/login":
			if !sec.NewLimiterMiddleware(r, windowSize, maxLoginTimeout) {
                http.Error(w, "Limite de taux dépassée", http.StatusTooManyRequests)
                return
            }
			Rt.LoginPage(w, r, tab)
		case "/logout":
		 if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
                http.Error(w, "Limite de taux dépassée", http.StatusTooManyRequests)
                return
            }
			Rt.LogOutHandler(w, r, tab)
		case "/home":
		 if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
                http.Error(w, "Limite de taux dépassée", http.StatusTooManyRequests)
                return
            }
			Rt.HomeHandler(w, r, tab)
		case "/myprofil/posts":
		 if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
                http.Error(w, "Limite de taux dépassée", http.StatusTooManyRequests)
                return
            }
			Rt.Profil(w, r, tab)
		case "/myprofil/favorites":
		 if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
                http.Error(w, "Limite de taux dépassée", http.StatusTooManyRequests)
                return
            }
			Rt.Profil_fav(w, r, tab)
		case "/myprofil/comments":
		 if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
                http.Error(w, "Limite de taux dépassée", http.StatusTooManyRequests)
                return
            }
			Rt.Profil_comment(w, r, tab)
		case "/filter":
		 if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
                http.Error(w, "Limite de taux dépassée", http.StatusTooManyRequests)
                return
            }
			Rt.Filter(w, r, tab)
		case "/index":
			if !sec.NewLimiterMiddleware(r, windowSize, maxRequests) {
                http.Error(w, "Limite de taux dépassée", http.StatusTooManyRequests)
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

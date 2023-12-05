package sec

import (
	"net/http"
	"sync"
	"time"
	"fmt"
)

// LoginwindowData est une structure qui contient les informations sur les requêtes et la dernière requête pour une adresse IP.
type LoginwindowData struct {
	requests    []time.Time
	lastRequest time.Time
}

// ipMap est une carte qui stocke les LoginwindowData pour chaque adresse IP.
var (
	ipMapp      = make(map[string]*LoginwindowData)
	ipMapMutexx sync.Mutex
)

// newLimiterMiddleware est une fonction middleware qui implémente le rate limiting.

func LoginLimiterMiddleware(r *http.Request, windowSize time.Duration, maxRequests int) bool {
	clientIP := r.RemoteAddr
	now := time.Now()

	// Verrouiller le mutex pour éviter l'accès concurrent à ipMap.
	ipMapMutexx.Lock()
	defer ipMapMutexx.Unlock()

	// Vérifier si l'adresse IP du client est déjà dans ipMap.
	if _, ok := ipMapp[clientIP]; !ok {
		// Si l'adresse IP du client n'est pas dans ipMap, créer une nouvelle LoginwindowData et l'ajouter à la carte.
		ipMapp[clientIP] = &LoginwindowData{
			requests:    []time.Time{now},
			lastRequest: now,
		}
		return true
	}

	// Si l'adresse IP du client est déjà dans ipMap, mettre à jour les champs requests et lastRequest.
	data := ipMapp[clientIP]
	data.requests = append(data.requests, now)
	data.lastRequest = now
	fmt.Println(now)

	// Supprimer les requêtes plus anciennes que windowSize.
	for i := len(data.requests) - 1; i >= 0; i-- {
		if now.Sub(data.requests[i]) > windowSize {
			data.requests = data.requests[:i+1]
			break
		}
	}

	// Vérifier si le nombre de requêtes est supérieur à maxRequests.
	fmt.Println(len(data.requests))
	if (len(data.requests) > maxRequests){
		data.requests = nil
		return false
	}
	return true
}


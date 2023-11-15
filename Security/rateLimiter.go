package sec

import (
	"net/http"
	"sync"
	"time"
)
// windowData est une structure qui contient les informations sur les requêtes et la dernière requête pour une adresse IP.
type windowData struct {
	requests    []time.Time
	lastRequest time.Time
}

// ipMap est une carte qui stocke les windowData pour chaque adresse IP.
var (
	ipMap      = make(map[string]*windowData)
	ipMapMutex sync.Mutex
)

// newLimiterMiddleware est une fonction middleware qui implémente le rate limiting.

func NewLimiterMiddleware(r *http.Request, windowSize time.Duration, maxRequests int) bool {
	clientIP := r.RemoteAddr
	now := time.Now()

	// Verrouiller le mutex pour éviter l'accès concurrent à ipMap.
	ipMapMutex.Lock()
	defer ipMapMutex.Unlock()

	// Vérifier si l'adresse IP du client est déjà dans ipMap.
	if _, ok := ipMap[clientIP]; !ok {
		// Si l'adresse IP du client n'est pas dans ipMap, créer une nouvelle windowData et l'ajouter à la carte.
		ipMap[clientIP] = &windowData{
			requests:    []time.Time{now},
			lastRequest: now,
		}
		return true
	}

	// Si l'adresse IP du client est déjà dans ipMap, mettre à jour les champs requests et lastRequest.
	data := ipMap[clientIP]
	data.requests = append(data.requests, now)
	data.lastRequest = now

	// Supprimer les requêtes plus anciennes que windowSize.
	for i := len(data.requests) - 1; i >= 0; i-- {
		if now.Sub(data.requests[i]) > windowSize {
			data.requests = data.requests[:i+1]
			break
		}
	}

	// Vérifier si le nombre de requêtes est supérieur à maxRequests.
	if len(data.requests) > maxRequests {
		return false
	}

	return true
}

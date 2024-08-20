package handlers

import (
	"fmt"
	"net/http"

	"github.com/JoshElias/go-web-server/internal"
)

func HandleMetricsAdmin(w http.ResponseWriter, r *http.Request) {
	var template = `<html>
<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>
</html>
`
	m := internal.GetMetrics()
	w.Header().Set("Content-Type", "text/html")
	m.Mu.Lock()
	hits := m.FileserverHits
	m.Mu.Unlock()
	w.Write([]byte(fmt.Sprintf(template, hits)))
}

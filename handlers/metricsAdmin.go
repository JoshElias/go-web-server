package handlers

import (
	"fmt"
	"net/http"

	"github.com/JoshElias/chirpy/config"
)

func HandleMetricsAdmin(w http.ResponseWriter, r *http.Request) {
	var template = `<html>
<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>
</html>
`
	c := config.GetConfig()
	w.Header().Set("Content-Type", "text/html")
	c.Mu.Lock()
	hits := c.FileserverHits
	c.Mu.Unlock()
	w.Write([]byte(fmt.Sprintf(template, hits)))
}

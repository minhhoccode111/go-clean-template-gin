package integration_test

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
)

// translationStubAddr is the port the app is pointed at by
// docker-compose-integration-test.yml. The test binary hosts the stub so the
// suite never depends on the public Google Translate endpoint.
const translationStubAddr = ":8443"

// startTranslationStub serves the two responses go-googletrans expects: the
// home page with a tkk token and the translate_a/single JSON payload.
func startTranslationStub() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/translate_a/single", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")

		response := struct {
			Sentences []struct {
				Trans   string `json:"trans"`
				Orig    string `json:"orig"`
				Backend int    `json:"backend"`
			} `json:"sentences"`
			Src string `json:"src"`
		}{
			Sentences: []struct {
				Trans   string `json:"trans"`
				Orig    string `json:"orig"`
				Backend int    `json:"backend"`
			}{{Trans: "translated: " + query, Orig: query, Backend: 1}},
			Src: r.URL.Query().Get("sl"),
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("translation stub: encode response: %v", err)
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("<html><script>tkk:'406398.2087938574'</script></html>")); err != nil {
			log.Printf("translation stub: write home page: %v", err)
		}
	})

	server := httptest.NewUnstartedServer(mux)

	listener, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", translationStubAddr)
	if err != nil {
		log.Fatalf("translation stub: listen on %s: %v", translationStubAddr, err)
	}

	server.Listener = listener
	server.StartTLS()

	log.Printf("translation stub: listening on %s", translationStubAddr)

	return server
}

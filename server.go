package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const hotReloadScript = `
<script>
(function() {
    const socket = new WebSocket('ws://' + window.location.host + '/ws');
    
    socket.onopen = function() {
        console.log('WebSocket connected');
    };
    
    socket.onmessage = function(event) {
        if (event.data === 'reload') {
            console.log('Reloading page...');
            window.location.reload();
        }
    };
    
    socket.onclose = function() {
        console.log('WebSocket disconnected, reconnecting in 1s...');
        setTimeout(function() {
            window.location.reload();
        }, 1000);
    };
})();
</script>
`

type FileWatcher struct {
	dir           string
	fileModTimes  map[string]time.Time
	clients       map[chan bool]bool
	checkInterval time.Duration
}

func NewFileWatcher(dir string) *FileWatcher {
	return &FileWatcher{
		dir:           dir,
		fileModTimes:  make(map[string]time.Time),
		clients:       make(map[chan bool]bool),
		checkInterval: 500 * time.Millisecond,
	}
}

func (fw *FileWatcher) AddClient(client chan bool) {
	fw.clients[client] = true
}

func (fw *FileWatcher) RemoveClient(client chan bool) {
	delete(fw.clients, client)
	close(client)
}

func (fw *FileWatcher) Start() {
	fw.updateFileModTimes()

	go func() {
		for {
			if fw.checkForChanges() {
				for client := range fw.clients {
					select {
					case client <- true:
					default:
					}
				}
			}
			time.Sleep(fw.checkInterval)
		}
	}()
}

func (fw *FileWatcher) updateFileModTimes() {
	newModTimes := make(map[string]time.Time)

	filepath.Walk(fw.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() || strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".html" || ext == ".css" || ext == ".js" || ext == ".md" {
			newModTimes[path] = info.ModTime()
		}

		return nil
	})

	fw.fileModTimes = newModTimes
}

func (fw *FileWatcher) checkForChanges() bool {
	changed := false

	filepath.Walk(fw.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() || strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".html" || ext == ".css" || ext == ".js" || ext == ".md" {
			if modTime, exists := fw.fileModTimes[path]; !exists || modTime != info.ModTime() {
				changed = true
				fw.fileModTimes[path] = info.ModTime()
			}
		}

		return nil
	})

	return changed
}

type WebSocketHandler struct {
	clients   map[chan bool]bool
	upgrader  *websocket.Upgrader
	fileWatch *FileWatcher
}

func NewWebSocketHandler(fileWatch *FileWatcher) *WebSocketHandler {
	return &WebSocketHandler{
		clients: make(map[chan bool]bool),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		fileWatch: fileWatch,
	}
}

func (wsh *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	client := make(chan bool, 1)
	wsh.fileWatch.AddClient(client)
	defer wsh.fileWatch.RemoveClient(client)

	for {
		select {
		case <-client:
			if err := conn.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
				return
			}
		case <-time.After(60 * time.Second):
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func HotReloadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".html") || r.URL.Path == "/" {
			recorder := &responseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           &strings.Builder{},
			}

			next.ServeHTTP(recorder, r)

			body := recorder.body.String()

			if strings.Contains(body, "</body>") {
				body = strings.Replace(body, "</body>", hotReloadScript+"</body>", 1)
			} else {
				body = body + hotReloadScript
			}

			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(recorder.statusCode)
			w.Write([]byte(body))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *strings.Builder
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	rr.body.Write(b)
	return rr.ResponseWriter.Write(b)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Printf("Error opening browser: %v", err)
	}
}

func main() {
	port := flag.Int("port", 8080, "Port to serve on")
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}

	fileWatcher := NewFileWatcher(dir)
	fileWatcher.Start()

	wsHandler := NewWebSocketHandler(fileWatcher)

	fs := http.FileServer(http.Dir(dir))

	http.Handle("/ws", wsHandler)
	http.Handle("/", HotReloadMiddleware(fs))

	serverAddr := fmt.Sprintf(":%d", *port)
	url := fmt.Sprintf("http://localhost%s", serverAddr)
	fmt.Printf("Starting server at %s\n", url)
	fmt.Printf("Serving files from: %s\n", dir)
	fmt.Println("Hot reload enabled - changes to HTML, CSS, JS, and MD files will automatically refresh the browser")
	fmt.Println("Press Ctrl+C to stop the server")

	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser(url)
	}()

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

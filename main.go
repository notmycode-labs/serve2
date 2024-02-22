package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

var (
	host   = flag.String("host", "0.0.0.0", "Host to serve on")
	port   = flag.Int("port", 8080, "Port to serve on")
	folder = flag.String("dir", ".", "Folder to serve")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", logRequest(handleRequest))
	addr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Serving files from %s on %s\n", *folder, addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ip := strings.Split(r.RemoteAddr, ":")[0]
		method := r.Method
		url := r.URL.Path
		userAgent := r.UserAgent()

		log.Printf("[%s] %s %s from %s (User-Agent: %s)\n", time.Now().Format(time.RFC3339), method, url, ip, userAgent)

		next.ServeHTTP(w, r)

		elapsed := time.Since(start)
		log.Printf("[%s] %s %s from %s completed in %s\n", time.Now().Format(time.RFC3339), method, url, ip, elapsed)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	filePath := path.Join(*folder, r.URL.Path[1:])
	if filePath == "" {
		renderDirectory(w, *folder)
		return
	}

	info, err := os.Stat(filePath)
	if err == nil && info.IsDir() {
		renderDirectory(w, filePath)
		return
	}

	http.ServeFile(w, r, filePath)
}

func humanizeSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

type PageData struct {
	DirPath  string
	FileList []string
}

func renderDirectory(w http.ResponseWriter, dirPath string) {
	dir, err := os.Open(dirPath)
	if err != nil {
		http.Error(w, "Failed to open directory", http.StatusInternalServerError)
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		http.Error(w, "Failed to read directory contents", http.StatusInternalServerError)
		return
	}

	type FileInfo struct {
		Name  string
		Size  string
		IsDir bool
		Path  string
	}

	var fileList []FileInfo
	for _, file := range files {
		fileInfo := FileInfo{
			Name:  file.Name(),
			Size:  humanizeSize(file.Size()),
			IsDir: file.IsDir(),
			Path:  path.Join("/", dirPath, file.Name()),
		}
		fileList = append(fileList, fileInfo)
	}

	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Index of {{.Name}}</title>
		
		<style>
		body {
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
			color: rgb(240, 240, 240);
			background-color: rgb(0, 0, 0);
			margin: 0;
			padding: 30px;
			-webkit-font-smoothing: antialiased;
		}
		tr td {
			list-style: none;
			justify-content: space-between;
			font-weight: 600;
			
		}

		td a {
			color: rgb(254, 242, 255);
			text-decoration: none;
		}
		
		
		</style>	
    </head>
    <body>
        <h1>Index of {{.Name}}</h1>
        <table>
            <thead>
                <tr>
                    <th>Type</th>
                    <th>Name</th>
                    <th>Size</th>
                </tr>
            </thead>
            <tbody>
                {{range .Files}}
                <tr>
                    <td>{{if .IsDir}}📁{{else}}📄{{end}}</td>
                    <td><a href="{{.Path}}">{{.Name}}</a></td>
                    <td>{{.Size}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </body>
    </html>`

	t := template.Must(template.New("listing").Parse(tmpl))
	err = t.Execute(w, struct {
		Name  string
		Files []FileInfo
	}{Name: dirPath, Files: fileList})
	if err != nil {
		http.Error(w, "Failed to render directory listing", http.StatusInternalServerError)
		return
	}
}

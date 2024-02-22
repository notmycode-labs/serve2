package main

import (
	"net/http"
	"os"
	"path"
	"text/template"
)

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
		<meta name="viewport" content="width=device-width, initial-scale=1" />

		
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

package main

import (
	"html/template"
	"log"
	"net/http"
	"syscall"
)

type StorageInfo struct {
	Total   uint64
	Free    uint64
	Used    uint64
	UsedPct float64
}

func getStorageInfo(path string) (StorageInfo, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return StorageInfo{}, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - free
	usedPct := float64(used) / float64(total) * 100

	return StorageInfo{
		Total:   total,
		Free:    free,
		Used:    used,
		UsedPct: usedPct,
	}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	nasPath := "/"
	storage, err := getStorageInfo(nasPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, storage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	port := "8080"
	log.Printf("Server started on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

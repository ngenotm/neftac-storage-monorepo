package main

import (
	"fmt"
	"os"
	"path/filepath"
	"neftac/storage/internal/storage"
)

func main() {
	if len(os.Args) < 3 { fmt.Println("usage: cli <cmd> <args>"); return }
	cmd := os.Args[1]
	switch cmd {
	case "upload-folder":
		local, bucket, prefix := os.Args[2], os.Args[3], os.Args[4]
		filepath.Walk(local, func(p string, i os.FileInfo, e error) error {
			if i.IsDir() { return nil }
			rel, _ := filepath.Rel(local, p)
			f, _ := os.Open(p)
			storage.WriteObject(bucket, filepath.Join(prefix, rel), f)
			fmt.Println("Uploaded:", rel)
			return nil
		})
	}
}

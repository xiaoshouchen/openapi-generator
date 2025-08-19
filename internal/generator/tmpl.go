package generator

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

//go:embed templates/*
var embeddedTemplates embed.FS

func InitTemplates() {
	// Create the .templates directory
	if err := os.MkdirAll(".templates", 0755); err != nil {
		log.Fatal("创建 .templates 目录失败:", err)
	}

	// Get the 'templates' subdirectory from the embedded filesystem
	subFS, err := fs.Sub(embeddedTemplates, "templates")
	if err != nil {
		log.Fatal("无法获取嵌入的 templates 目录:", err)
	}

	// Walk the embedded templates and copy them
	err = fs.WalkDir(subFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Create the destination path
		destPath := filepath.Join(".templates", path)

		if d.IsDir() {
			// Create the directory if it's a directory
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return err
			}
		} else {
			// Read the file content from the embedded FS
			content, readErr := fs.ReadFile(subFS, path)
			if readErr != nil {
				return readErr
			}
			// Write the file to the destination
			if writeErr := os.WriteFile(destPath, content, 0644); writeErr != nil {
				return writeErr
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal("复制模板文件失败:", err)
	}
	log.Println("模板文件已成功生成到 .templates 目录.")
}

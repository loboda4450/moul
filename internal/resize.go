package internal

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

// Collection struct
type Collection struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Src      string `json:"src"`
	Color    string `json:"color"`
	SrcHd    string `json:"src_hd"`
	Width    int    `json:"width"`
	WidthHd  int    `json:"width_hd"`
	Height   int    `json:"height"`
	HeightHd int    `json:"height_hd"`
}

// get file path
func getFilePath(uid, prefix string, size int) string {
	return filepath.Join(".moul", "photos", uid, prefix, strconv.Itoa(size))
}

// GetFileName func
func GetFileName(fn, author string) string {
	return strings.TrimSuffix(fn, filepath.Ext(fn)) + "-by-" + author
}

// GetPhotoDimension given path
func GetPhotoDimension(path string) (int, int) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
	}

	return image.Width, image.Height
}

// resize image
func manipulate(size int, inPath, author, unique, outPrefix string) {
	src, err := imaging.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}

	fn := filepath.Base(inPath)
	name := GetFileName(fn, author)

	dir := getFilePath(unique, outPrefix, size)
	out := filepath.Join(dir, name+".jpg")
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal(err)
	}

	newImage := imaging.Resize(src, size, 0, imaging.Lanczos)

	err = imaging.Save(newImage, out)
	if err != nil {
		log.Fatal(err)
	}
}

// GetPhotos given path
func GetPhotos(path string) []string {
	var photos []string
	// folder to walk through
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".jpeg" || filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".png" {
			photos = append(photos, path)
		}
		return nil
	})

	if err != nil {
		log.Println(err)
	}

	return photos
}

// open image
func loadImage(fileInput string) (image.Image, error) {
	f, err := os.Open(fileInput)
	defer f.Close()
	if err != nil {
		log.Println("File not found:", fileInput)
		return nil, err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// Resize func
func Resize(inPath, author, outPrefix, unique string, sizes []int) {
	// weight := "100"
	// author := "sophearak-tha"
	// path := "path-to-src/photos"
	// unique := UniqueID()
	// sizes := []int{2560, 1280, 620},

	photos := GetPhotos(inPath)

	for _, photo := range photos {
		for _, size := range sizes {
			manipulate(size, photo, author, unique, outPrefix)
		}
	}
}
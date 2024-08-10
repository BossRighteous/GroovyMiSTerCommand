package generators

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func WriteGMCtoDir(dir string, filename string, content []byte) {
	os.MkdirAll(dir, os.ModePerm)
	gmcPath := filepath.Join(dir, filename+".gmc")
	fo, err := os.Create(gmcPath)
	if err != nil {
		fmt.Printf("Unable to create file, may exist %s\n", gmcPath)
		return
	}
	length, err := fo.Write(content)
	if err != nil || length != len(content) {
		log.Fatal(err)
	}
	if err := fo.Close(); err != nil {
		log.Fatal(err)
	}
}

func ExecutableDir() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal("FATAL Executable not referenced")
	}
	return filepath.Dir(ex)
}

func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
func CutSuffix(s, suffix string) (before string, found bool) {
	if !HasSuffix(s, suffix) {
		return s, false
	}
	return s[:len(s)-len(suffix)], true
}

func GetBaseFilename(path string) (before string, found bool) {
	return CutSuffix(filepath.Base(path), filepath.Ext(path))
}

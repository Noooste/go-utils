package logger

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

const logsPath = "./logs"

var errorLock = &sync.Mutex{}
var infoLock = &sync.Mutex{}
var debugLock = &sync.Mutex{}

func createFile(name string) {
	if file, err := os.Create(logsPath + name); err == nil {
		_ = file.Close()
	}
}

func writeToFile(name string, content []byte) {
	file, err := os.OpenFile(logsPath+name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	_, _ = file.Write(content)
	err = file.Close()
}

func Init() {
	_ = os.Mkdir(logsPath, os.ModePerm)
	createFile("/errors.log")
	createFile("/info.log")
	createFile("/debug.log")
}

func SaveError(error string) {
	errorLock.Lock()
	defer errorLock.Unlock()
	writeToFile("/errors.log", []byte(time.Now().Format(`Mon Jan _2 15:04:05`)+" | "+error))
}

func SaveInfo(info string) {
	infoLock.Lock()
	defer infoLock.Unlock()
	writeToFile("/info.log", []byte(time.Now().Format(`Mon Jan _2 15:04:05`)+" | "+info))
}

func SaveDebug(debug string) {
	debugLock.Lock()
	defer debugLock.Unlock()
	writeToFile("/debug.log", []byte(time.Now().Format(`Mon Jan _2 15:04:05`)+" | "+debug))
}

func ZipWriter(path, dest string) {
	files, _ := ioutil.ReadDir(path)
	if len(files) == 0 {
		return
	}
	outFile, _ := os.Create(dest)
	defer outFile.Close()
	w := zip.NewWriter(outFile)
	addFiles(w, path, "")
	w.Close()
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	files, _ := ioutil.ReadDir(basePath)
	for _, file := range files {
		if !file.IsDir() {
			dat, _ := ioutil.ReadFile(basePath + file.Name())
			f, _ := w.Create(baseInZip + file.Name())
			_, _ = f.Write(dat)
		} else if file.IsDir() {
			newBase := basePath + file.Name() + "/"
			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}

func SaveLogs(dir string) {
	if dir[len(dir)-1:] != "/" {
		dir += "/"
	}

	name := time.Now().Format("2006-01-02_15-04-05")

	if stat, err := os.Stat(dir); err == nil && stat.Size() > 0 {
		ZipWriter(dir, strings.TrimRight(dir, "/")+"_"+name+".zip")
		dirRead, _ := os.Open(dir)
		dirFiles, _ := dirRead.Readdir(0)

		for index := range dirFiles {
			fileHere := dirFiles[index]

			nameHere := fileHere.Name()
			fullPath := dir + nameHere

			os.Remove(fullPath)
			fmt.Println("Removed file:", fullPath)
		}
	}
}

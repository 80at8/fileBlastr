package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

type operationStatus int

// Constants for fileStatus
const (
	STALE = iota
	ACTIVE
)

type fileHashEntry struct {
	createTime     time.Time
	lastAccessTime time.Time
	entryAge       time.Time
	fileBytes      int64
	fileName       string
	fileVersion    int
	fileUUID       uuid.UUID
	fileHash       string
	fileStatus     operationStatus
}

type fileHashMap struct {
	files     map[string]fileHashEntry
	filesPath string
	mutex     sync.Mutex
	logBuffer bytes.Buffer
	logger    *log.Logger
}

// add a new fileHashEntry
func (f *fileHashMap) add(entry fileHashEntry) {
	f.mutex.Lock()

	entry.createTime = time.Now()
	entry.entryAge = time.Now()
	entry.fileStatus = ACTIVE
	f.files[entry.fileHash] = entry

	f.mutex.Unlock()
}

// computes hash of a file.
func (f *fileHashMap) hash(entry fileHashEntry) fileHashEntry {
	fileHandle, err := os.Open(path.Clean(f.filesPath + "/" + entry.fileName))
	if err != nil {
		f.log("[current time:" + time.Now().Format(time.RFC3339) + "] hash():os.Open() fileName=" + entry.fileName + " " + err.Error())
	}
	defer fileHandle.Close()

	hash := sha256.New()

	fileBytes, err := io.Copy(hash, fileHandle)

	if err != nil {
		f.log("[current time:" + time.Now().Format(time.RFC3339) + "] hash():io.Copy() fileBytes=" + strconv.FormatInt(fileBytes, 10) + " " + err.Error())
	}
	entry.fileHash = base64.URLEncoding.EncodeToString(hash.Sum(nil))
	entry.fileBytes = fileBytes
	return entry
}

// remove an entry (by hash key) from the fileHashMap
func (f *fileHashMap) remove(hash string) {
	f.mutex.Lock()

	delete(f.files, hash)

	f.mutex.Unlock()
}

// refresh entry time, < 30 == ACTIVE, > 30 == STALE
func (f *fileHashMap) refresh() {
	f.mutex.Lock()

	for k, v := range f.files {
		if time.Since(v.entryAge).Minutes() < 30.0 {
			v.fileStatus = ACTIVE
			f.files[k] = v
			f.log("[current time:" + time.Now().Format(time.RFC3339) + "] refresh() fileHashMap[" + k + "] key [entryAge:" + v.entryAge.Format(time.RFC3339) + "] marked ACTIVE")
		} else {
			v.fileStatus = STALE
			f.files[k] = v
			f.log("[current time:" + time.Now().Format(time.RFC3339) + "] refresh() fileHashMap[" + k + "] key [" + v.entryAge.Format(time.RFC3339) + "] marked STALE")
		}
	}

	f.mutex.Unlock()
}

// initialize the filepath and logging facilities
func (f *fileHashMap) init(filesPath string) {
	f.filesPath = filesPath
	f.logger = log.New(&f.logBuffer, "fileHashMap:", log.Ltime|log.Lshortfile)
}

// create a logging entry.
func (f *fileHashMap) log(logEntry string) {
	f.logger.Printf(logEntry)
}

// func main() {
// 	var fileBlaster fileHashMap

// 	fileBlaster.init("./data")

// 	fileBlaster.files = make(map[string]fileHashEntry)
// 	for x := 0; x < 10; x++ {
// 		var testHash fileHashEntry
// 		testHash.fileHash = "HASH" + strconv.Itoa(x)
// 		testHash.fileName = "testfile" + strconv.Itoa(x) + ".txt"
// 		fileBlaster.add(testHash)
// 	}
// 	fileBlaster.refresh()

// }

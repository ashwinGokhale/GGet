package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Downloader struct {
	client *http.Client
	filename string
	threads int
}

func (dl *Downloader) singleThreaded(startTime time.Time) error {
	Info("Starting single threaded download...")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil { return err }
	resp, err := dl.client.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil { return err }
	file, err := os.OpenFile(dl.filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil { return err }
	defer file.Close()
	Info("Writing to %s", dl.filename)
	file.Write(body)
	Info("Downloaded %d bytes to %s in %s",len(body), dl.filename, time.Now().Sub(startTime))
	return nil
}


func (dl *Downloader) threadedDownload(length int, startTime time.Time) error {
	Info("Starting threaded download...")
	size := length / dl.threads
	remainder := length % dl.threads
	Info("Downloading %s on %d threads", dl.filename, dl.threads)
	wg := &sync.WaitGroup{}
	for i := 0; i < dl.threads; i++ {
		wg.Add(1)

		start := i * size
		end := (i+1) * size

		if i == dl.threads-1 {
			end += remainder
		}

		Info("Starting thread %d", i)
		go func(start, end, i int) error {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				wg.Done()
				return err
			}
			byteRange := fmt.Sprintf("bytes=%d-%d", start, end-1)
			req.Header.Add("Range", byteRange)
			resp, err := dl.client.Do(req)
			if err != nil {
				wg.Done()
				return err
			}
			defer resp.Body.Close()
			Info("Thread: %d Reading response body", i)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				wg.Done()
				return err
			}
			file, err := os.OpenFile(dl.filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				wg.Done()
				return err
			}
			defer file.Close()
			io.Copy(file, resp.Body)
			Info("Thread: %d writing bytes %d - %d", i, start, end)
			file.WriteAt(body, int64(start))
			wg.Done()
			Info("Thread: %d done", i)
			return nil
		}(start, end, i)
	}
	wg.Wait()
	Info("Downloaded %s in %s", dl.filename, time.Now().Sub(startTime))
	return nil
}

func Info(format string, args ...interface{}) {
	fmt.Printf(color.BlueString("[INFO] ")+format+"\n", args...)
}
// Copyright © 2017 Ashwin Gokhale ashwin.gokhale98@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"runtime"
	"github.com/PuerkitoBio/purell"
	"github.com/goware/urlx"
	"github.com/spf13/cobra"
)

var (
	downloader    = &Downloader{
		threads: 1,
		client: &http.Client{},
	}
	url 		  string
	single        bool
	root           = &cobra.Command{
		Use:	"gget [url]",
		Short:	"GGet is a multithreaded accelerated downloader",
		Long:	`GGet uses go routines to split up the download
				 of files into multiple concurrent processes`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if downloader.threads < 0 {
				return errors.New("Invalid thread count: " + strconv.Itoa(downloader.threads))
			}
			if downloader.filename == "" {
				downloader.filename = url[strings.LastIndex(url, "/")+1:]
			}
			return nil
		},
		RunE: run,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Must specify a url")
			}
			u, err := urlx.NormalizeString(purell.MustNormalizeURLString(args[0],purell.FlagsUsuallySafeGreedy))
			url = u
			if err != nil { return err }
			return nil
		},
	}
)

func init()  {
	root.PersistentFlags().IntVarP(&downloader.threads,"threads","t", runtime.NumCPU(), "Number of threads to run")
	root.PersistentFlags().StringVarP(&(downloader.filename),"filename", "f", "", "Specify a filename")
	root.PersistentFlags().BoolVarP(&single, "single-threaded", "s", false, "Use single threaded download")
}

func run(cmd *cobra.Command, args []string) error {
	startTime := time.Now()
	Info("Going to download %s from %s", downloader.filename, url)
	resp, err := downloader.client.Head(url)
	if err != nil { return err }
	contentLength := resp.Header.Get("Content-Length")
	ranges := resp.Header.Get("Accept-Ranges")
	if contentLength == "" {
		Info("Content length not specified")
		return downloader.singleThreaded(startTime)
	}
	if ranges != "bytes" {
		Info("Server does not accept byte ranges")
		return downloader.singleThreaded(startTime)
	}
	if single || downloader.threads <= 1 {
		return downloader.singleThreaded(startTime)
	}
	length, err := strconv.Atoi(contentLength)
	if err != nil { return err }
	return downloader.threadedDownload(length, startTime)
}

func main() {
	root.Execute()
}
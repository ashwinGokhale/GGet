# GGet

### What is it?
* GGet is a download accelerator written in Go. By default it uses as many threads 
as there are CPUs on your machine. You can specify the number of threads and filename.

### Usage
```bash
Usage:
  gget [url] [flags]

Flags:
  -f, --filename string   Specify a filename
  -h, --help              help for gget
  -s, --single-threaded   Use single threaded download
  -t, --threads int       Number of threads to run
```

#### Todo
* Add a progress bar
* Add more options (Make more like wget)

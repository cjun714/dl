package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: dl '<URL>' <num>")
		fmt.Println("sample: dl 'http://cs248.stanford.edu/lectures/04_texture/images' 77")
		return
	}

	base := os.Args[1]
	num, e := strconv.Atoi(os.Args[2])
	if e != nil {
		panic(e)
	}

	var wg sync.WaitGroup
	for i := 1; i < num+1; i++ {
		fileName := fmt.Sprintf("slide_%03d.jpg", i)
		url := base + "/" + fileName
		wg.Add(1)
		go func() {
			defer wg.Done()
			if e = download(url); e != nil {
				fmt.Fprintln(os.Stderr, url+" failed:", e.Error())
			}
		}()
	}
	wg.Wait()
}

func download(url string) error {
	fmt.Println(url)

	resp, e := http.Get(url)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	idx := strings.LastIndex(url, "/")
	if idx == -1 {
		return errors.New("can't pase filename:" + url)
	}
	name := url[idx:]

	f, e := os.Create("./" + name)
	if e != nil {
		return e
	}
	defer f.Close()

	_, e = io.Copy(f, resp.Body)
	if e != nil {
		return e
	}

	fmt.Println(url, "  ok")

	return nil
}

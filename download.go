package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	flagList   = flag.String("l", "list.txt", "file list")
	flagThread = flag.Int("t", 4, "threads")
)

type Task struct {
	Src  string
	Dest string
}

func work(t Task) {
	if _, err := os.Stat(t.Dest); err == nil {
		// exists
		log.Println("exist", t.Dest)
		return
	}
	log.Println("working", t.Src)
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(t.Src)
	if err != nil {
		log.Println("ERROR1", err)
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERROR2", err)
		return
	}
	destDir := path.Dir(t.Dest)
	os.MkdirAll(destDir, os.ModePerm)
	check(ioutil.WriteFile(t.Dest, content, os.ModePerm))
}

func main() {
	flag.Parse()

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt)

	queue := make(chan Task)

	wg := sync.WaitGroup{}
	wg.Add(*flagThread)
	for i := 0; i < *flagThread; i++ {
		go func() {
			for t := range queue {
				work(t)
			}
			wg.Done()
		}()
	}

	file, err := os.Open(*flagList)
	check(err)
	list := bufio.NewReader(file)
	defer file.Close()
loop:
	for {
		line, err := list.ReadString('\n')
		if err != nil {
			break
		}
		parts := strings.Fields(line)
		task := Task{
			Src:  parts[0],
			Dest: parts[1],
		}
		select {
		case queue <- task:
		case <-done:
			log.Println("stopping...")
			break loop
		}
	}
	close(queue)

	wg.Wait()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

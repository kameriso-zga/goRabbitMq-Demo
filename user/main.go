package main

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"
)

var query = "test"
var matches int

var workerCount = 0
var maxWorkerCount = 32
var wg sync.WaitGroup

//channel 互相交流共享内存

//包工头指派活
var searchReq = make(chan string)

//通知包工头，自己的活是否干完了
var workDone = make(chan bool)

//传输我是否找到搜索结果的消息
var foundMatch = make(chan bool)

func main() {
	start := time.Now()
	workerCount = 1

	go search("/", true)
	waitForWorkers()
	fmt.Println(matches, "mactches")
	fmt.Println(time.Since(start))

}

//case 来收听几个channel
func waitForWorkers() {
	for {
		select {
		//有新的活干了
		case path := <-searchReq:
			//工人数量+1  指派新的工人干活
			workerCount++
			go search(path, true)
		case <-workDone:
			workerCount--
			if workerCount == 0 {
				return
			}
		case <-foundMatch:
			matches++
		}
	}
}

func search(path string, master bool) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, file := range files {
		name := file.Name()
		if name == query {
			foundMatch <- true
		}
		if file.IsDir() {
			//判断是否有空余的工人
			if workerCount < maxWorkerCount {
				searchReq <- path + name + "/ "
			} else {
				search(path+name+"/", false)
			}
		}
	}
	//master 是用来判断
	if master {
		workDone <- true
	}
}

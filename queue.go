package main

import (
    "comic/services"
)

//执行队列监听
func main() {
    services.QueueListen()
}

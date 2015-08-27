package main

import "github.com/fpgeek/docker-gc/docker"

func main() {
	docker.RemoveExistedContainers()
	docker.RemoveStaleImages()
}

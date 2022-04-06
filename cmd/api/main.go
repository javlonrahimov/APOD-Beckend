package main

import (
	"javlonrahimov/apod/internal/data"
	"javlonrahimov/apod/internal/jsonlog"
	"sync"
)

type config struct {
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {

}

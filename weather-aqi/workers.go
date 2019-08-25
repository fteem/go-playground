package main

type worker interface {
	work(out chan<- string)
}

type workerGroup struct {
	city          string
	aqiWorker     worker
	darkskyWorker worker
}

func newWorkerGroup(city string, aqiWorker worker, darkskyWorker worker) workerGroup {
	return workerGroup{
		city:          city,
		aqiWorker:     aqiWorker,
		darkskyWorker: darkskyWorker,
	}
}

func (wg workerGroup) work(out chan<- row) {
	indexChan := make(chan string)
	darkskyChan := make(chan string)

	go wg.aqiWorker.work(indexChan)
	go wg.darkskyWorker.work(darkskyChan)

	out <- row{
		wg.city,
		<-darkskyChan,
		<-indexChan,
	}
}

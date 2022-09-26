package shuttle

import (
	"log"
	"sync"
	"time"
)

const (
	// channel control operator
	opDelData    = "opDelData"
	opAddData    = "opAddData"
	opHandleData = "opHandleData"
	opExit       = "opExit"
	// default maximum message pool size
	maxSize = 20

	defaultDuration = 30 * time.Second
)

// dataPool you know
type dataPool struct {
	pool        map[string]interface{}
	totalSize   int // current number of data in the dataPool
	maxSize     int // max number of data in the dataPool
	channel     chan *handleChan
	duration    time.Duration
	shuttleFunc ShuttleFunc
}

// handleChan channel struct
type handleChan struct {
	op    string
	key   string
	value interface{}
}

// ShuttleFunc a function that will be triggered periodically
type ShuttleFunc func(m map[string]interface{}) error

var (
	mutex            sync.Mutex
	once             sync.Once
	dataPollInstance *dataPool
	timer            *time.Ticker
	defaultFunc      = func(m map[string]interface{}) error {
		return nil
	}
)

// Pool pool interface
type Pool interface {
	AddData(key string, value interface{})
	DelData(key string)
	Exit()
}

// AddData insert into data pool
func (a *dataPool) AddData(key string, value interface{}) {
	ch := &handleChan{
		op:    opAddData,
		value: value,
		key:   key,
	}
	a.channel <- ch
}

// DelData delete from data pool
func (a *dataPool) DelData(key string) {
	ch := &handleChan{
		op:  opDelData,
		key: key,
	}
	a.channel <- ch
}

// Exit stop the loop (used for data saving after a service panic)
func (a *dataPool) Exit() {
	ch := &handleChan{
		op: opExit,
	}
	a.channel <- ch
}

// newDataPool init a dataPool
func newDataPool(options ...Option) Pool {
	once.Do(func() {
		dataPollInstance = &dataPool{
			pool:      make(map[string]interface{}),
			totalSize: 0,
			channel:   make(chan *handleChan),
		}
		for k := range options {
			options[k](dataPollInstance)
		}
		dataPollInstance.initDataPoll()
	})
	return dataPollInstance
}

// initDataPoll init data pool
func (a *dataPool) initDataPoll() {
	if a.shuttleFunc == nil {
		a.shuttleFunc = defaultFunc
	}
	a.initTimer()
	go a.handleLoop()
	log.Printf("data pool init finish")
}

func (a *dataPool) registerFunc(f func(map[string]interface{}) error) {
	if a.shuttleFunc == nil {
		a.shuttleFunc = f
	}
}

// initTimer init timing trigger coroutine
func (a *dataPool) initTimer() {
	duration := defaultDuration
	if a.duration != 0 {
		duration = a.duration
	}
	timer = time.NewTicker(duration)
	go func() {
		for {
			<-timer.C
			trigger := &handleChan{
				op: opHandleData,
			}
			a.channel <- trigger
		}
	}()
}

// handleLoop master loop
func (a *dataPool) handleLoop() {
	for c := range a.channel {
		switch c.op {
		case opDelData:
			a.delData(c)
		case opAddData:
			a.addData(c)
		case opHandleData:
			a.handleData()
		case opExit:
			a.handleData()
			return
		default:
			log.Printf("unknown operation: %s", c.op)
		}
	}
}

func (a *dataPool) delData(c *handleChan) {
	if a.totalSize == 0 {
		return
	}
	delete(a.pool, c.key)
	a.totalSize -= 1
}

func (a *dataPool) addData(c *handleChan) {
	a.pool[c.key] = c.value
	a.totalSize += 1
	if a.totalSize >= maxSize {
		a.handleData()
	}
}

func (a *dataPool) handleData() {
	defer a.clearAll()
	if a.totalSize == 0 {
		return
	}
	if err := a.shuttleFunc(a.pool); err != nil {
		log.Printf("shuttle func execute error: %v", err)
	}
}

// clearAll clear data
func (a *dataPool) clearAll() {
	mutex.Lock()
	a.pool = make(map[string]interface{})
	a.totalSize = 0
	mutex.Unlock()
}

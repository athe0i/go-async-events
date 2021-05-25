package async_events

import (
	"container/list"
	"context"
	"sync"
)

type Listener interface {
	AddChannel(chn ReaderChannel) error
	GetDispatcher() EventDispatcher
	Run(ctx context.Context) error
	SetErrorHandler(eh ErrorHandler)
	HandlesErrors
}

type EventListener struct {
	repo         Repository
	disp         EventDispatcher
	errorHandler ErrorHandler

	chns []ReaderChannel

	mu    sync.Mutex
	queue *list.List

	sema       chan int
	loopSignal chan struct{}

	wg *sync.WaitGroup
}

func NewEventListener(repo Repository, disp EventDispatcher) Listener {
	return &EventListener{
		repo:         repo,
		disp:         disp,
		errorHandler: LogError,

		chns: []ReaderChannel{},

		queue: list.New(),

		sema:       make(chan int, 3),
		loopSignal: make(chan struct{}, 1),
		wg:         &sync.WaitGroup{},
	}
}

func (rl *EventListener) AddChannel(chn ReaderChannel) (err error) {
	if rl.isChannelRegistered(chn) {
		return ListenerChannelAlreadyRegistered{Chn: chn.GetName()}
	}

	rl.chns = append(rl.chns, chn)

	return
}

func (rl *EventListener) GetDispatcher() (disp EventDispatcher) {
	return rl.disp
}

func (rl *EventListener) isChannelRegistered(chn ReaderChannel) bool {
	for _, lch := range rl.chns {
		if chn.GetName() == lch.GetName() {
			return true
		}
	}

	return false
}

func (rl *EventListener) Run(ctx context.Context) (err error) {
	go rl.loop(ctx)

	go rl.fillQueue(ctx)

	rl.wg.Wait()

	return
}

func (rl *EventListener) loop(ctx context.Context) {
	for {
		select {
		case <-rl.loopSignal:
			rl.pickEventFromQueue()
		case <-ctx.Done():
			return
		}
	}
}

func (rl *EventListener) pickEventFromQueue() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.queue.Len() == 0 {
		return
	}

	select {
	case rl.sema <- 1:
		ev := rl.dequeue()
		go rl.processEvent(ev)
	}
}

func (rl *EventListener) processEvent(ev Event) {
	defer rl.replenish()
	defer rl.wg.Done()

	err := rl.disp.HandleEvent(ev)
	if err != nil {
		rl.HandleError(err)
	}
}

func (rl *EventListener) replenish() {
	<-rl.sema

	rl.tickleLoop()
}

func (rl *EventListener) dequeue() (ev Event) {
	el := rl.queue.Front()
	rl.queue.Remove(el)
	return el.Value.(Event)
}

func (rl *EventListener) tickleLoop() {
	select {
	case rl.loopSignal <- struct{}{}:
	}
}

func (rl *EventListener) fillQueue(ctx context.Context) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for _, ch := range rl.chns {
		ev, err := ch.Pop()

		// we might get redis.Nil if nothing to return, we don't care about it
		if err != nil {
			rl.HandleError(err)
		}
		if ev != nil {
			rl.queue.PushBack(ev)
			rl.tickleLoop()
		}
	}

	if rl.queue.Len() == 0 {
		ctx.Done()
	} else {
		rl.wg.Add(rl.queue.Len())
	}
}

func (rl *EventListener) HandleError(err error) {
	rl.errorHandler(err)
}
func (rl *EventListener) SetErrorHandler(errHandl ErrorHandler) {
	rl.errorHandler = errHandl
}

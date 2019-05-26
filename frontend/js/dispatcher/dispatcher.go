package dispatcher

type Listener func(data interface{})

type Dispatcher struct {
	listeners map[int][]Listener
}

func New() *Dispatcher {
	return &Dispatcher{
		listeners: make(map[int][]Listener),
	}
}

func (this *Dispatcher) Dispatch(eventId int, data interface{}) bool {
	listeners := this.listeners[eventId]
	if len(listeners) < 0 {
		return false
	}

	for _, listener := range listeners {
		listener(data)
	}

	return true
}

func (this *Dispatcher) AddListener(eventId int, listener Listener) {
	listeners := this.listeners[eventId]
	listeners = append(listeners, listener)
	this.listeners[eventId] = listeners
}

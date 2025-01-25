package broadcast

var (
	channels = make(map[string]chan interface{})
)

func SendSignal(qualifier string) {
	if channels[qualifier] != nil {
		channels[qualifier] <- struct{}{}
	}
}

func Channel(qualifier string) chan interface{} {
	if channels[qualifier] == nil {
		channels[qualifier] = make(chan interface{})
	}

	return channels[qualifier]
}

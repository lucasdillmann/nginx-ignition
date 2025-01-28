package broadcast

import "dillmann.com.br/nginx-ignition/core/common/log"

var (
	channels = make(map[string]chan interface{})
)

func SendSignal(qualifier string) {
	if channels[qualifier] != nil {
		channels[qualifier] <- struct{}{}
	} else {
		log.Warnf("Signal ignored: qualifier %s has no listeners yet", qualifier)
	}
}

func Listen(qualifier string) chan interface{} {
	if channels[qualifier] == nil {
		channels[qualifier] = make(chan interface{})
	}

	return channels[qualifier]
}

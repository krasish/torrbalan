package logutil

import "log"

func LogOnErr(f func() error) {
	LogAllOnErr(f)
}

func LogAllOnErr(fs ...func() error) {
	for _, f := range fs {
		err := f()
		if err != nil {
			log.Println(err)
		}
	}
}

package logutil

import "log"

func LogOnErr(f func() error) {
	LogAllOnErr(f)
}

//LogAllOnErr executes all the functions from fs and logs all the errors that
//those functions may have returned.
func LogAllOnErr(fs ...func() error) {
	for _, f := range fs {
		err := f()
		if err != nil {
			log.Println(err)
		}
	}
}

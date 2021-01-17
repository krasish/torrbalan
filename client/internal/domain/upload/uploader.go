package upload

type Uploader struct {
	q chan string
}

func NewUploader(concurrentUploads uint) Uploader {
	return Uploader{
		q: make(chan string, concurrentUploads),
	}
}

func (u Uploader) Start() {
	for {
		filename := <-u.q
		u.processUploading(filename)
	}
}

func (u Uploader) Upload(filename, partnerAddress string) {
	u.q <- filename
}

func (u Uploader) processUploading(filename string) {
	//TODO: Consider what happens when two simultaneous
	// downloads for the same file are started.
}

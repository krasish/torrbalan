package download

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"

	"github.com/krasish/torrbalan/client/pkg/eofutil"

	"github.com/krasish/torrbalan/client/internal/logutil"
)

const BufferSize = 1024

//Info represents all the information needed to download a file from other client
type Info struct {
	Filename    string
	PeerAddress string
	PathToSave  string
}

//Downloader is responsible for downloading files from other clients.
type Downloader struct {
	q chan Info
}

func NewDownloader(concurrentDownloads uint) Downloader {
	return Downloader{
		q: make(chan Info, concurrentDownloads),
	}
}

//Start should be started in a separate goroutine. It waits on the Info chan of d
//and starts goroutines handling the download from the given input.
func (d Downloader) Start() {
	for {
		info := <-d.q
		rootPath := path.Clean(info.PathToSave)
		filePath := rootPath + "/" + info.Filename
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Could not create file: %v\n", err)
			continue
		}

		conn, err := d.connectToPeer(info)
		if err != nil {
			fmt.Printf("Could not connect: %v\n", err)
			logutil.LogOnErr(file.Close)
			continue
		}
		go d.processDownloading(file, conn)
	}
}

func (d Downloader) Download(info Info) {
	d.q <- info
}

func (d Downloader) connectToPeer(info Info) (net.Conn, error) {
	conn, err := net.Dial("tcp", info.PeerAddress)
	if err != nil {
		return nil, fmt.Errorf("while trying to connect to %s: %w", info.PeerAddress, err)
	}
	return conn, nil
}

//initialContract should be called before processDownloading.
//initialContract sends/reads messages to/form client on conn by which both clients
//can determine which file is about to get transferred and determine whetehr
//that is possible, returning an error if not.
func (d Downloader) initialContract(conn net.Conn, fileName string) error {
	readWriter := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	handler := eofutil.LoggingEOFHandler{DestName: conn.RemoteAddr().String()}

	err := eofutil.WriteCheckEOF(readWriter.Writer, fileName+"$", handler)
	if err != nil {
		return fmt.Errorf("while writing filename: %v", err)
	}
	resp, err := eofutil.ReadCheckEOF(readWriter.Reader, '$', handler)
	if err != nil {
		return fmt.Errorf("while waiting for first response from peer: %v", err)
	}
	if resp != "OK$" {
		return fmt.Errorf("file %q cannot be downloaded from %s", fileName, conn.RemoteAddr().String())
	}
	return nil
}

//processDownloading loops reading from client and writing to the file.
//Both the net.Conn and *file.File are closed when downloading is done or an error occurred.
func (d Downloader) processDownloading(file *os.File, conn net.Conn) {
	defer logutil.LogAllOnErr(file.Close, conn.Close)
	fileName := path.Base(file.Name())
	remoteAddr := conn.RemoteAddr().String()
	if err := d.initialContract(conn, fileName); err != nil {
		log.Printf("An error occurred while establishinng initial contract with %s: %v", remoteAddr, err)
		return
	}

	reader, writer := bufio.NewReader(conn), bufio.NewWriter(file)
	errorMessages := [5]string{
		fmt.Sprintf("Stopping download of file %q from %s. EOF read through connection.", file.Name(), remoteAddr),
		fmt.Sprintf("An error occurred while reading from %s: %%v", remoteAddr),
		fmt.Sprintf("Stopping download of file %q from %s. EOF while writing to file.", file.Name(), remoteAddr),
		fmt.Sprintf("An error occurred while writing in file %s: %%v", file.Name()),
		fmt.Sprintf("final flush to file %s failed", file.Name()),
	}

	ReadWriteLoop(reader, writer, errorMessages)
	log.Printf("Finished downloading %s", fileName)
}

//ReadWriteLoop is an loop which reads from reader and writes to writer until and error or
//io.EOF occurs logging respective messages for any scenario. It flushes writer after error.
func ReadWriteLoop(reader *bufio.Reader, writer *bufio.Writer, errorMessages [5]string) {
	bytes := make([]byte, BufferSize)
	for {
		n, err := reader.Read(bytes)
		if err == io.EOF {
			log.Println(errorMessages[0])
			break
		} else if err != nil {
			log.Printf(errorMessages[1], err)
			break
		}
		_, err = writer.Write(bytes[:n])
		if err == io.EOF {
			log.Println(errorMessages[2])
			break
		} else if err != nil {
			log.Printf(errorMessages[3], err)
			break
		}
	}
	if err := writer.Flush(); err != nil {
		log.Println(errorMessages[4])
	}
}

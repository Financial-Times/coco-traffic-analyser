package analyser

import (
	"bufio"
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type httpStreamFactory struct {
	networkGraph *WeightedGraph
}

func newHttpStreamFactory(networkGraph *WeightedGraph) *httpStreamFactory {
	return &httpStreamFactory{networkGraph}
}

type httpRequestCollector struct {
	httpRequests chan *http.Request
}

// httpStream will handle the actual decoding of http requests.
type httpStream struct {
	net, transport gopacket.Flow
	streamReader   tcpreader.ReaderStream
	networkGraph   *WeightedGraph
}

func (f *httpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &httpStream{
		net:          net,
		transport:    transport,
		streamReader: tcpreader.NewReaderStream(),
		networkGraph: f.networkGraph,
	}

	go hstream.run() // Important... we must guarantee that data from the reader stream is read.
	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &hstream.streamReader
}

func (stream *httpStream) run() {
	buf := bufio.NewReader(&stream.streamReader)
	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			// We must read until we see an EOF... very important!
			return
		} else if err != nil {
			log.Error("Error reading stream", stream.net, stream.transport, ":", err)
		} else {
			bodyBytes := tcpreader.DiscardBytesToEOF(req.Body)
			req.Body.Close()
			source := req.UserAgent()
			target := req.URL.Host
			log.Info(source, target)
			stream.networkGraph.add(source, target)
			log.Debug("Received request from stream", stream.net, stream.transport, ":", req, "with", bodyBytes, "bytes in request body")
		}
	}
}

package analyser

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

type Analyser interface {
	TrafficGraph() *WeightedGraph
}

type StandardAnalyser struct {
	iface string
	graph *WeightedGraph
}

func New(iface string) *StandardAnalyser {
	return &StandardAnalyser{
		iface: iface,
		graph: newWeightedGraph(),
	}
}

func (a *StandardAnalyser) TrafficGraph() *WeightedGraph {
	return a.graph
}

func (a *StandardAnalyser) Start() {
	handle, err := pcap.OpenLive(a.iface, 1600, true, pcap.BlockForever)

	if err != nil {
		log.Error(err)
	}

	// Set up assembly
	streamFactory := newHttpStreamFactory(a.graph)
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	log.Info("reading in packets")
	// Read in packets, pass to assembler.
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates the end of a pcap file.
			if packet == nil {
				return
			}

			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Debug("Unusable packet")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		case <-ticker:
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}

}

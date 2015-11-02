package mtr

import (
	"time"
)

const DEFAULT_PORT = 33434
const DEFAULT_MAX_HOPS = 64
const DEFAULT_TIMEOUT_MS = 500

const DEFAULT_RETRIES = 5
const DEFAULT_PACKET_SIZE = 52
const DEFAULT_SNT_SIZE = 10

type TracerouteReturn struct {
	Success bool
	Addr    string
	Elapsed time.Duration
}

type TracerouteHop struct {
	Success  bool
	Address  string
	Host     string
	N        int
	TTL      int
	LastTime time.Duration
	AvgTime  time.Duration
	BestTime time.Duration
	WrstTime time.Duration
	Loss     float32
}

type TracerouteResult struct {
	DestAddress [4]byte
	Hops        []TracerouteHop
}

type TracerouteOptions struct {
	port       int
	maxHops    int
	timeoutMs  int
	retries    int
	packetSize int
	sntSize    int
}

func (options *TracerouteOptions) Port() int {
	if options.port == 0 {
		options.port = DEFAULT_PORT
	}
	return options.port
}

func (options *TracerouteOptions) SetPort(port int) {
	options.port = port
}

func (options *TracerouteOptions) MaxHops() int {
	if options.maxHops == 0 {
		options.maxHops = DEFAULT_MAX_HOPS
	}
	return options.maxHops
}

func (options *TracerouteOptions) SetMaxHops(maxHops int) {
	options.maxHops = maxHops
}

func (options *TracerouteOptions) TimeoutMs() int {
	if options.timeoutMs == 0 {
		options.timeoutMs = DEFAULT_TIMEOUT_MS
	}
	return options.timeoutMs
}

func (options *TracerouteOptions) SetTimeoutMs(timeoutMs int) {
	options.timeoutMs = timeoutMs
}

func (options *TracerouteOptions) Retries() int {
	if options.retries == 0 {
		options.retries = DEFAULT_RETRIES
	}
	return options.retries
}

func (options *TracerouteOptions) SetRetries(retries int) {
	options.retries = retries
}

func (options *TracerouteOptions) SntSize() int {
	if options.sntSize == 0 {
		options.sntSize = DEFAULT_SNT_SIZE
	}
	return options.sntSize
}

func (options *TracerouteOptions) SetSntSize(sntSize int) {
	options.sntSize = sntSize
}

func (options *TracerouteOptions) PacketSize() int {
	if options.packetSize == 0 {
		options.packetSize = DEFAULT_PACKET_SIZE
	}
	return options.packetSize
}

func (options *TracerouteOptions) SetPacketSize(packetSize int) {
	options.packetSize = packetSize
}

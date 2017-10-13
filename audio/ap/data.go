package ap

// audio file decoder
type Decoder interface {
	// helper method for in-memory decode
	FullDecode() (d []byte, numChan, freq int32, err error)

	// stream decode
	Decode() int
	NumOfChan() int
	BitDepth() int
	SampleRate() int32
	Buffer() []byte
	ReachEnd() bool
}

// decoder factory, we use'll used it to
// create new decoder by file-type
type DecoderFactory interface {
	NewDecoder(name string, fileType FileType) (Decoder, error)
}

// sound represent a audio segment
type Sound struct {
	Type SourceType
	Format uint32
	Priority uint16

	Data interface{}
}

// static in-memory data
type StaticData struct {
	Buffer BufferAL

	SampleRate int32
	BitDepth   int32
	NumOfChan  int32
}

// streamed from file
type StreamData struct {
	Decoder
}



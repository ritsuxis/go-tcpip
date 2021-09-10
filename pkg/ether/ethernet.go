package ethernet


const (
	headerSize = 14
	trailerSize = 0 // FCSの部分だけど、ここではFCSを抜いているので0
	maxPayloadSize = 1500
	minPayloadSize = 14

	maxFrameSize = headerSize + maxPayloadSize + trailerSize
	minFrameSize = headerSize + minPayloadSize + trailerSize
)


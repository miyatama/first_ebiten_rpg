package audio

import (
	_ "embed"
)

var (
	//go:embed jab.wav
	Jab_wav []byte

	//go:embed jab8.wav
	Jab8_wav []byte

	//go:embed jump.ogg
	Jump_ogg []byte

	//go:embed ragtime.mp3
	Ragtime_mp3 []byte

	//go:embed ragtime.ogg
	Ragtime_ogg []byte
)

const (
	DEFAULT_SAMPLE_RATE = 48000
)

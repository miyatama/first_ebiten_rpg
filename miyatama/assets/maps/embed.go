package maps

import (
	"bytes"
	_ "embed"
	"image"
	"log/slog"

	"first_rpg/miyatama/assets/images"
)

var (
	//go:embed title_layer01.csv
	TitleLayer01 string
)

func LoadMovableMap() (map[MovableKey]bool, error) {
	slog.Info("maps LoadMovableMap()")
	maps := make(map[MovableKey]bool)
	img, _, err := image.Decode(bytes.NewReader(images.TitleMovableMapImage))
	if err != nil {
		slog.Error("title_movable_map_image.png decode error")
		slog.String("error: {}", err.Error())
		return maps, err
	}
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	for i := 0; i < (width / images.MAP_TILE_WIDTH); i++ {
		for j := 0; j < (height / images.MAP_TILE_WIDTH); j++ {
			left := i * images.MAP_TILE_WIDTH
			top := j * images.MAP_TILE_WIDTH
			_, _, _, a := img.At(left, top).RGBA()
			maps[MovableKey{X: i, Y: j}] = a == 0
		}
	}

	return maps, nil
}

type MovableKey struct {
	X int
	Y int
}

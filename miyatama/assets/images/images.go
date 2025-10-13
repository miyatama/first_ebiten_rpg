package images

import (
	"bytes"
	_ "embed"
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MAP_TILE_COLS  = 8
	MAP_TILE_WIDTH = 32
)

var (
	//go:embed player.png
	Player_png []byte

	//go:embed black_cat.png
	BlackCat_png []byte

	//go:embed title_map_image.png
	TitleMapImage []byte

	//go:embed title_movable_map_image.png
	TitleMovableMapImage []byte
)

func GetTitleMapImage() (*ebiten.Image, error) {
	slog.Info("images GetTitleMapImage()")
	img, _, err := image.Decode(bytes.NewReader(TitleMapImage))
	if err != nil {
		slog.Error("title_map_image.png decode error")
		slog.String("error: {}", err.Error())
		return ebiten.NewImage(1, 1), err
	}
	return ebiten.NewImageFromImage(img), nil
}

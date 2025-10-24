package scenes

import (
	MiyatamaImages "first_rpg/miyatama/ui/assets/images"
	gamestatus "first_rpg/miyatama/ui/game_status"
	"image"
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

type SceneSwitcherBehavior int

const (
	SCENE_SWITCHER_BEHAVIOR_IDLE SceneSwitcherBehavior = iota
	SCENE_SWITCHER_BEHAVIOR_IRIS_IN
	SCENE_SWITCHER_BEHAVIOR_IRIS_OUT
)

const (
	// マスク画像の透過最小矩形幅
	SCENE_SWITCH_MASK_TRANSPARENT_MIN_WIDTH = 197
	SCENE_SWITCH_MASK_CENTER                = 256
	SCENE_SWITCH_IRIS_FRAME_COUNT           = 30
)

type SceneSwitcher struct {
	calcPanelRectHeight   int
	calcPanelRectWidth    int
	switchBackgroundImage *ebiten.Image
	maskImage             *ebiten.Image
	blackImage            *ebiten.Image
	maskIrisOutStartScale float64
	irisFrameProcessCount int
	Behavior              SceneSwitcherBehavior
}

func (s *SceneSwitcher) Init() error {
	maskImage, err := MiyatamaImages.GetImage(MiyatamaImages.SceneSwitchMask)
	if err != nil {
		slog.Error("failed to get scene switch mask image", "err", err)
		return err
	}

	slog.Info("SceneSwitcher.Init()",
		slog.Int("mask bound x", maskImage.Bounds().Dx()),
		slog.Int("mask bound y", maskImage.Bounds().Dy()),
	)
	alphamap := image.NewAlpha(image.Rectangle{image.Point{}, maskImage.Bounds().Max})
	for i := 0; i < maskImage.Bounds().Dx(); i++ {
		for j := 0; j < maskImage.Bounds().Dy(); j++ {
			_, _, _, a := maskImage.At(i, j).RGBA()
			if a != 0 {
				alphamap.Set(i, j, color.Alpha{0x0})
			} else {
				alphamap.Set(i, j, color.Alpha{0xff})
			}
		}
	}
	s.maskImage = ebiten.NewImageFromImage(alphamap)

	s.Behavior = SCENE_SWITCHER_BEHAVIOR_IDLE
	return nil
}

func (s *SceneSwitcher) Update(data *gamestatus.GameData) {
	if s.calcPanelRectHeight != data.LayoutHeight || s.calcPanelRectWidth != data.LayoutWidth {
		// generate foreground image
		s.switchBackgroundImage = ebiten.NewImage(data.LayoutWidth, data.LayoutHeight)
		s.blackImage = ebiten.NewImage(data.LayoutWidth, data.LayoutHeight)
		s.blackImage.Fill(color.Black)

		// calc mask in scale
		s.maskIrisOutStartScale = float64(max(data.LayoutWidth, data.LayoutHeight)) / float64(SCENE_SWITCH_MASK_TRANSPARENT_MIN_WIDTH)

		slog.Info(
			"SceneSwitcher.Update()",
			slog.Float64("mask iris out scale", s.maskIrisOutStartScale),
		)
		s.calcPanelRectHeight = data.LayoutHeight
		s.calcPanelRectWidth = data.LayoutWidth
	}
	switch s.Behavior {
	case SCENE_SWITCHER_BEHAVIOR_IRIS_OUT:
		s.irisFrameProcessCount--
		if s.irisFrameProcessCount <= 0 {
			s.Behavior = SCENE_SWITCHER_BEHAVIOR_IDLE
			s.switchBackgroundImage.Fill(color.Black)
		}
	case SCENE_SWITCHER_BEHAVIOR_IRIS_IN:
		s.irisFrameProcessCount++
		if s.irisFrameProcessCount >= SCENE_SWITCH_IRIS_FRAME_COUNT {
			s.Behavior = SCENE_SWITCHER_BEHAVIOR_IDLE
		}
	}
}

func (s *SceneSwitcher) Draw(screen *ebiten.Image, data *gamestatus.GameData) {
	if s.calcPanelRectHeight != data.LayoutHeight || s.calcPanelRectWidth != data.LayoutWidth {
		return
	}
	switch s.Behavior {
	case SCENE_SWITCHER_BEHAVIOR_IDLE:
		screen.DrawImage(s.switchBackgroundImage, nil)
	case SCENE_SWITCHER_BEHAVIOR_IRIS_OUT, SCENE_SWITCHER_BEHAVIOR_IRIS_IN:
		s.switchBackgroundImage.Fill(color.White)
		op := &ebiten.DrawImageOptions{}
		op.Blend = ebiten.BlendCopy
		scale := float64(s.irisFrameProcessCount) / float64(SCENE_SWITCH_IRIS_FRAME_COUNT) * s.maskIrisOutStartScale
		scaledWidth := float64(s.maskImage.Bounds().Dx()) * scale
		scaledHeight := float64(s.maskImage.Bounds().Dy()) * scale
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(-scaledWidth/2, -scaledHeight/2)
		op.GeoM.Translate(float64(data.LayoutWidth/2), float64(data.LayoutHeight/2))
		s.switchBackgroundImage.DrawImage(s.maskImage, op)
		op = &ebiten.DrawImageOptions{}
		op.Blend = ebiten.BlendSourceIn
		s.switchBackgroundImage.DrawImage(s.blackImage, op)

		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0)
		screen.DrawImage(s.switchBackgroundImage, op)
	}
}

func (s *SceneSwitcher) IsIdle() bool {
	return s.Behavior == SCENE_SWITCHER_BEHAVIOR_IDLE
}

func (s *SceneSwitcher) StartIrisOut() {
	s.Behavior = SCENE_SWITCHER_BEHAVIOR_IRIS_OUT
	s.irisFrameProcessCount = SCENE_SWITCH_IRIS_FRAME_COUNT
}

func (s *SceneSwitcher) StartIrisIn() {
	s.Behavior = SCENE_SWITCHER_BEHAVIOR_IRIS_IN
	s.irisFrameProcessCount = 0
}

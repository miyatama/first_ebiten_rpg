package scenes

import (
	"first_rpg/miyatama/assets/images"
	maps "first_rpg/miyatama/assets/maps"
	gamestatus "first_rpg/miyatama/game_status"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TITLE_MAP_ROWS     = 50
	TITLE_MAP_COLS     = 50
	MOVING_FRAME_COUNT = 30
)

type TitleSceneStatus int

const (
	IDLE TitleSceneStatus = iota
	MOVING
)

/**
* player ユーザが操作するプレイヤー
* mapLayer マップの配置情報
* mapParts 描画で利用するマップのパーツ
**/
type TitleScene struct {
	player                *Player
	mapImage              *ebiten.Image
	movableMap            map[maps.MovableKey]bool
	currentPlayerPosition PlayerPosition
	nextPlayerPosition    PlayerPosition
	beforeImageDrawX      int
	beforeImageDrawY      int
	sceneStatus           TitleSceneStatus
	movingFrame           int
}

type PlayerPosition struct {
	X int
	Y int
}

func (t *TitleScene) Init() error {
	t.player = &Player{}
	if err := t.player.Init(); err != nil {
		return err
	}
	slog.Info("TitleScene.Init",
		slog.Bool("initialized Player", true),
	)

	// マップの配置情報をロード
	mapImage, err := images.GetTitleMapImage()
	if err != nil {
		return err
	}
	t.mapImage = mapImage

	// 移動可能なマップの情報
	t.movableMap, err = maps.LoadMovableMap()
	if err != nil {
		return err
	}

	// プレイヤー初期位置
	t.currentPlayerPosition = PlayerPosition{
		X: 24,
		Y: 21,
	}
	return nil
}

func (t *TitleScene) Update(data *gamestatus.GameData) {
	// キャラクタの移動
	if t.sceneStatus == IDLE {
		if isMove(data.UserAction) {
			nextX, nextY := getNextPosition(t.currentPlayerPosition.X, t.currentPlayerPosition.Y, data.UserAction)
			key := maps.MovableKey{X: nextX, Y: nextY}
			if t.movableMap[key] {
				t.nextPlayerPosition = PlayerPosition{
					X: nextX,
					Y: nextY,
				}
				t.sceneStatus = MOVING
				t.movingFrame = 0
			}
		}
	}
	t.player.Update(data)
}

func (t *TitleScene) Draw(screen *ebiten.Image, data *gamestatus.GameData) {
	mapSx := (t.currentPlayerPosition.X*images.MAP_TILE_WIDTH + images.MAP_TILE_WIDTH/2)
	mapSy := (t.currentPlayerPosition.Y*images.MAP_TILE_WIDTH + images.MAP_TILE_WIDTH/2)
	// 移動中のフレーム判定
	if t.sceneStatus == MOVING {
		deltaX := (t.nextPlayerPosition.X - t.currentPlayerPosition.X) * images.MAP_TILE_WIDTH
		deltaY := (t.nextPlayerPosition.Y - t.currentPlayerPosition.Y) * images.MAP_TILE_WIDTH
		deltaX = int(float32(deltaX) / float32(MOVING_FRAME_COUNT) * float32(t.movingFrame))
		deltaY = int(float32(deltaY) / float32(MOVING_FRAME_COUNT) * float32(t.movingFrame))
		mapSx += deltaX
		mapSy += deltaY
		t.movingFrame++
		if t.movingFrame > MOVING_FRAME_COUNT {
			t.currentPlayerPosition = t.nextPlayerPosition
			t.sceneStatus = IDLE
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(mapSx), -float64(mapSy))
	op.GeoM.Translate(float64(data.LayoutWidth)/2, float64(data.LayoutHeight)/2)
	screen.DrawImage(t.mapImage, op)

	if t.beforeImageDrawX != mapSx || t.beforeImageDrawY != mapSy {
		t.beforeImageDrawX = mapSx
		t.beforeImageDrawY = mapSy
	}
	t.player.Draw(screen, data)
}

func isMove(userAction gamestatus.UserAction) bool {
	return userAction == gamestatus.USER_ACTION_LEFT ||
		userAction == gamestatus.USER_ACTION_RIGHT ||
		userAction == gamestatus.USER_ACTION_UP ||
		userAction == gamestatus.USER_ACTION_DOWN
}

func getNextPosition(currentX, currentY int, userAction gamestatus.UserAction) (int, int) {
	nextX, nextY := currentX, currentY

	if userAction == gamestatus.USER_ACTION_LEFT {
		nextX = currentX - 1
	}
	if userAction == gamestatus.USER_ACTION_RIGHT {
		nextX = currentX + 1
	}
	if userAction == gamestatus.USER_ACTION_UP {
		nextY = currentY - 1
	}
	if userAction == gamestatus.USER_ACTION_DOWN {
		nextY = currentY + 1
	}

	if nextX < 0 {
		nextX = 0
	}
	if nextY < 0 {
		nextY = 0
	}
	if nextX >= TITLE_MAP_COLS {
		nextX = TITLE_MAP_COLS - 1
	}
	if nextY >= TITLE_MAP_ROWS {
		nextY = TITLE_MAP_ROWS - 1
	}
	return nextX, nextY
}

package scenes

import (
	"bytes"
	miyatamaAudio "first_rpg/miyatama/assets/audio"
	"first_rpg/miyatama/assets/events"
	"first_rpg/miyatama/assets/images"
	maps "first_rpg/miyatama/assets/maps"
	gamestatus "first_rpg/miyatama/game_status"
	"first_rpg/miyatama/util"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
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
	TALK_MOB
)

/**
* player ユーザが操作するプレイヤー
* mapLayer マップの配置情報
* mapParts 描画で利用するマップのパーツ
**/
type TitleScene struct {
	gameStateMsg          gamestatus.GameStateMsg
	player                *Player
	talkPanel             *TalkPanel
	mapImage              *ebiten.Image
	movableMap            map[util.MapPosition]bool
	currentPlayerPosition util.MapPosition
	nextPlayerPosition    util.MapPosition
	beforeImageDrawX      int
	beforeImageDrawY      int
	sceneStatus           TitleSceneStatus
	movingFrame           int
	mobs                  []*MobCharacter
	events                []*events.Event
	// Audio
	audioContext *audio.Context
	audioPlayer  *audio.Player
}

func (t *TitleScene) Init() error {
	t.gameStateMsg = gamestatus.GAME_STATE_MSG_NONE
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
	t.currentPlayerPosition = util.MapPosition{
		X: 24,
		Y: 21,
	}

	// モブキャラクター
	t.mobs = []*MobCharacter{}
	t.mobs = append(t.mobs, generateMobCharacter()...)
	for _, m := range t.mobs {
		m.Init()
	}

	// イベント
	t.events = []*events.Event{}
	t.events = append(t.events, generateEvents()...)

	// トーク用パネル
	t.talkPanel = &TalkPanel{}
	t.talkPanel.Init()

	// BGM
	audioContext := audio.NewContext(miyatamaAudio.DEFAULT_SAMPLE_RATE)
	t.audioContext = audioContext
	audioStream, err := mp3.DecodeF32(bytes.NewReader(miyatamaAudio.Ragtime_mp3))
	if err != nil {
		return err
	}
	audioPlayer, err := audioContext.NewPlayerF32(audioStream)
	if err != nil {
		return err
	}
	t.audioPlayer = audioPlayer
	t.audioPlayer.Play()
	return nil
}

func (t *TitleScene) Update(data *gamestatus.GameData) {
	// キャラクタの移動
	if t.sceneStatus == IDLE {
		f := func() {
			if t.movePlayer(data) {
				return
			}

			if t.actionMobCharacter(data) {
				return
			}
		}
		f()
	} else if t.sceneStatus == TALK_MOB {
		f := func() {
			// イベントを次に進める
			if data.UserAction == gamestatus.USER_ACTION_DECIDE {
				t.sceneStatus = IDLE
				data.Event = nil
			}
		}
		f()
	}
	for _, m := range t.mobs {
		m.Update(data)
	}
	t.player.Update(data)
	t.talkPanel.Update(data)
}

func (t *TitleScene) Draw(screen *ebiten.Image, data *gamestatus.GameData) {
	// マップの描画
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

	// モブの描画
	for _, m := range t.mobs {
		m.SetDrawCorrection(mapSx, mapSy)
		m.Draw(screen, data)
	}

	if t.beforeImageDrawX != mapSx || t.beforeImageDrawY != mapSy {
		t.beforeImageDrawX = mapSx
		t.beforeImageDrawY = mapSy
	}
	t.player.Draw(screen, data)

	// モブ会話の描画
	if t.sceneStatus == TALK_MOB {
		t.talkPanel.Draw(screen, data)
	}
}

func (t *TitleScene) Msg() gamestatus.GameStateMsg {
	return t.gameStateMsg
}

func (t *TitleScene) movePlayer(data *gamestatus.GameData) bool {
	if !isInputDirection(data.UserAction) {
		return false
	}

	nextX, nextY := getNextPosition(t.currentPlayerPosition.X, t.currentPlayerPosition.Y, data.UserAction)
	key := util.MapPosition{X: nextX, Y: nextY}
	existsMobCharacter, _ := t.existsMobCharacter(nextX, nextY)
	movable := t.movableMap[key] && !existsMobCharacter
	if movable {
		t.nextPlayerPosition = util.MapPosition{
			X: nextX,
			Y: nextY,
		}
		t.sceneStatus = MOVING
		t.movingFrame = 0
		t.player.SetUserAction(data.UserAction)
		return true
	} else {
		return false
	}
}

func (t *TitleScene) actionMobCharacter(data *gamestatus.GameData) bool {
	if !isInputDirection(data.UserAction) {
		return false
	}

	// 進行方向にモブが存在するか
	nextX, nextY := getNextPosition(t.currentPlayerPosition.X, t.currentPlayerPosition.Y, data.UserAction)
	existsMobCharacter, mobIndex := t.existsMobCharacter(nextX, nextY)
	if existsMobCharacter {
		t.player.SetUserAction(data.UserAction)
		t.sceneStatus = TALK_MOB
		eventId := t.mobs[mobIndex].EventId
		for _, e := range t.events {
			if e.Id == eventId {
				data.Event = e
				break
			}
		}
		slog.Info("TitleScene.actionMobCharacter",
			slog.Int("mob index", mobIndex),
		)
		return true
	} else {
		return false
	}
}

func isInputDirection(userAction gamestatus.UserAction) bool {
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

func (t *TitleScene) existsMobCharacter(x, y int) (bool, int) {
	for i, m := range t.mobs {
		if m.Position.X == x && m.Position.Y == y {
			return true, i
		}
	}
	return false, 0
}

func generateMobCharacter() []*MobCharacter {
	return []*MobCharacter{
		&MobCharacter{
			MobType: MOB_TYPE_BLACK_CAT,
			Position: util.MapPosition{
				X: 39,
				Y: 5,
			},
			Direction: util.DIRECTION_DOWN,
			EventId:   0,
		},
		&MobCharacter{
			MobType: MOB_TYPE_BLACK_CAT,
			Position: util.MapPosition{
				X: 24,
				Y: 20,
			},
			Direction: util.DIRECTION_DOWN,
			EventId:   0,
		},
		&MobCharacter{
			MobType: MOB_TYPE_NONE,
			Position: util.MapPosition{
				X: 44,
				Y: 24,
			},
			Direction: util.DIRECTION_DOWN,
			EventId:   1,
		},
		&MobCharacter{
			MobType: MOB_TYPE_VILLAGE_BOY,
			Position: util.MapPosition{
				X: 44,
				Y: 23,
			},
			Direction: util.DIRECTION_DOWN,
			EventId:   1,
		},
	}
}

func generateEvents() []*events.Event {
	return []*events.Event{
		&events.Event{
			Id:        0,
			EventType: events.EVENT_TYPE_MOB_TALK,
			TalkTexts: []string{"シャー！"},
		},
		&events.Event{
			Id:        1,
			EventType: events.EVENT_TYPE_RESTAULANT,
			TalkTexts: []string{"いらっしゃい、なににしますか？"},
		},
	}
}

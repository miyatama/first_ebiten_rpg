package scenes

import (
	"bytes"
	"errors"
	miyatamaAudio "first_rpg/miyatama/assets/audio"
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
	HOUSE_MAP_ROWS = 10
	HOUSE_MAP_COLS = 10
)

type HouseSceneStatus int

const (
	HOUSE_SCENE_STATUS_IDLE HouseSceneStatus = iota
	HOUSE_SCENE_STATUS_START_SCENE
	HOUSE_SCENE_STATUS_MOVING
	HOUSE_SCENE_STATUS_PROCESS_EVENT
	HOUSE_SCENE_CHANGE_SCENE
)

/**
* player ユーザが操作するプレイヤー
* mapLayer マップの配置情報
* mapParts 描画で利用するマップのパーツ
**/
type HouseScene struct {
	gameStateMsg          gamestatus.GameStateMsg
	player                *Player
	talkPanel             *TalkPanel
	gamepad               *GamePad
	sceneSwitcher         *SceneSwitcher
	mapImage              *ebiten.Image
	movableMap            map[util.MapPosition]bool
	currentPlayerPosition util.MapPosition
	nextPlayerPosition    util.MapPosition
	beforeImageDrawX      int
	beforeImageDrawY      int
	sceneStatus           HouseSceneStatus
	movingFrame           int
	mobs                  []*MobCharacter
	processEvent          *gamestatus.Event
	processEventTalkSeq   int
	events                []*gamestatus.Event
	userInputFrame        int
	// Audio
	audioPlayer *audio.Player
}

func (h *HouseScene) Init(data *gamestatus.GameData) error {
	h.gameStateMsg = gamestatus.GAME_STATE_MSG_NONE
	h.gamepad = &GamePad{}
	if err := h.gamepad.Init(); err != nil {
		return err
	}

	h.player = &Player{}
	if err := h.player.Init(); err != nil {
		return err
	}
	slog.Info("HouseScene.Init",
		slog.Bool("initialized Player", true),
	)

	// マップの配置情報をロード
	mapImage, err := images.GetHouseMapImage()
	if err != nil {
		return err
	}
	h.mapImage = mapImage

	// 移動可能なマップの情報
	h.movableMap, err = maps.LoadHouseMovableMap()
	if err != nil {
		return err
	}

	// プレイヤー初期位置
	h.currentPlayerPosition = util.MapPosition{
		X: 5,
		Y: 7,
	}

	// モブキャラクター
	h.mobs = []*MobCharacter{}
	h.mobs = append(h.mobs, h.generateMobCharacter()...)
	for _, m := range h.mobs {
		m.Init()
	}

	// イベント
	h.events = []*gamestatus.Event{}
	h.events = append(h.events, h.generateEvents()...)

	// トーク用パネル
	h.talkPanel = &TalkPanel{}
	h.talkPanel.Init()

	// BGM
	audioStream, err := mp3.DecodeF32(bytes.NewReader(miyatamaAudio.Ragtime_mp3))
	if err != nil {
		return err
	}
	audioPlayer, err := data.AudioContext.NewPlayerF32(audioStream)
	if err != nil {
		return err
	}
	h.audioPlayer = audioPlayer
	h.audioPlayer.Play()

	// シーン切り替え
	h.sceneSwitcher = &SceneSwitcher{}
	h.sceneSwitcher.Init()
	h.sceneSwitcher.StartIrisIn()
	return nil
}

func (h *HouseScene) Update(data *gamestatus.GameData) {
	h.gamepad.Update(data)
	// キャラクタの移動
	switch h.sceneStatus {
	case HOUSE_SCENE_STATUS_IDLE:
		{
			f := func() {
				if h.movePlayer(data) {
					return
				}

				if h.actionMobCharacter(data) {
					return
				}
			}
			f()
		}
	case HOUSE_SCENE_STATUS_PROCESS_EVENT:
		{
			f := func() {
				switch h.processEvent.EventType {
				case gamestatus.EVENT_TYPE_MOB_TALK:
					{
						// イベントを次に進める
						if !h.allowUserInput() {
							return
						}
						if data.UserAction != gamestatus.USER_ACTION_DECIDE {
							return
						}
						if len(h.processEvent.TalkTexts) > h.processEventTalkSeq+1 {
							h.processEventTalkSeq++
							h.talkPanel.SetText(h.processEvent.TalkTexts[h.processEventTalkSeq])
							return
						}
						if h.processEvent.NextEventId <= 0 {
							h.sceneStatus = HOUSE_SCENE_STATUS_IDLE
							h.processEvent = nil
							return
						}
						evevnt, _, err := h.getEvent(h.processEvent.NextEventId)
						if err != nil {
							slog.Error(err.Error(),
								slog.Int("event id", h.processEvent.NextEventId),
							)
							return
						}
						h.processEvent = evevnt
					}
				case gamestatus.EVENT_TYPE_CHANGE_SCENE:
					{
						h.sceneStatus = HOUSE_SCENE_CHANGE_SCENE
						h.sceneSwitcher.StartIrisOut()
					}
				}
			}
			f()
		}
	}
	for _, m := range h.mobs {
		m.Update(data)
	}
	h.player.Update(data)
	h.talkPanel.Update(data, h.gamepad.GetPadRect())
	h.sceneSwitcher.Update(data)
}

func (h *HouseScene) Draw(screen *ebiten.Image, data *gamestatus.GameData) {
	// マップの描画
	mapSx := (h.currentPlayerPosition.X*images.MAP_TILE_WIDTH + images.MAP_TILE_WIDTH/2)
	mapSy := (h.currentPlayerPosition.Y*images.MAP_TILE_WIDTH + images.MAP_TILE_WIDTH/2)
	// 移動中のフレーム判定
	if h.sceneStatus == HOUSE_SCENE_STATUS_MOVING {
		deltaX := (h.nextPlayerPosition.X - h.currentPlayerPosition.X) * images.MAP_TILE_WIDTH
		deltaY := (h.nextPlayerPosition.Y - h.currentPlayerPosition.Y) * images.MAP_TILE_WIDTH
		deltaX = int(float32(deltaX) / float32(PLAYER_MOVING_FRAME_COUNT) * float32(h.movingFrame))
		deltaY = int(float32(deltaY) / float32(PLAYER_MOVING_FRAME_COUNT) * float32(h.movingFrame))
		mapSx += deltaX
		mapSy += deltaY
		h.movingFrame++
		if h.movingFrame > PLAYER_MOVING_FRAME_COUNT {
			h.currentPlayerPosition = h.nextPlayerPosition
			h.sceneStatus = HOUSE_SCENE_STATUS_IDLE
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(mapSx), -float64(mapSy))
	op.GeoM.Translate(float64(data.LayoutWidth)/2, float64(data.LayoutHeight)/2)
	screen.DrawImage(h.mapImage, op)

	// モブの描画
	for _, m := range h.mobs {
		m.SetDrawCorrection(mapSx, mapSy)
		m.Draw(screen, data)
	}

	if h.beforeImageDrawX != mapSx || h.beforeImageDrawY != mapSy {
		h.beforeImageDrawX = mapSx
		h.beforeImageDrawY = mapSy
	}
	h.player.Draw(screen, data)

	// モブ会話の描画
	if h.processEvent != nil {
		switch h.processEvent.EventType {
		case gamestatus.EVENT_TYPE_MOB_TALK:
			{
				h.talkPanel.Draw(screen, data)
			}
		}
	}

	// コントローラーの描画
	h.gamepad.Draw(screen, data)
	if h.userInputFrame > 0 {
		h.userInputFrame--
	}
	h.sceneSwitcher.Draw(screen, data)
}

func (h *HouseScene) Msg() gamestatus.GameStateMsg {
	return h.gameStateMsg
}

func (h *HouseScene) movePlayer(data *gamestatus.GameData) bool {
	if !h.isInputDirection(data.UserAction) {
		return false
	}

	nextX, nextY := h.getNextPosition(h.currentPlayerPosition.X, h.currentPlayerPosition.Y, data.UserAction)
	key := util.MapPosition{X: nextX, Y: nextY}
	existsMobCharacter, _ := h.existsMobCharacter(nextX, nextY)
	movable := h.movableMap[key] && !existsMobCharacter
	if movable {
		h.nextPlayerPosition = util.MapPosition{
			X: nextX,
			Y: nextY,
		}
		h.sceneStatus = HOUSE_SCENE_STATUS_MOVING
		h.movingFrame = 0
		h.player.SetUserAction(data.UserAction)
		return true
	} else {
		return false
	}
}

func (h *HouseScene) actionMobCharacter(data *gamestatus.GameData) bool {
	if !h.isInputDirection(data.UserAction) {
		return false
	}

	// 進行方向にモブが存在するか
	nextX, nextY := h.getNextPosition(h.currentPlayerPosition.X, h.currentPlayerPosition.Y, data.UserAction)
	existsMobCharacter, mobIndex := h.existsMobCharacter(nextX, nextY)
	if existsMobCharacter {
		h.player.SetUserAction(data.UserAction)
		h.sceneStatus = HOUSE_SCENE_STATUS_PROCESS_EVENT
		eventId := h.mobs[mobIndex].EventId
		for _, e := range h.events {
			if e.Id == eventId {
				h.processEvent = e
				h.processEventTalkSeq = 0
				h.talkPanel.SetText(h.processEvent.TalkTexts[h.processEventTalkSeq])
				break
			}
		}
		slog.Info("HouseScene.actionMobCharacter",
			slog.Int("mob index", mobIndex),
		)
		return true
	} else {
		return false
	}
}

func (h *HouseScene) allowUserInput() bool {
	if h.userInputFrame <= 0 {
		h.userInputFrame = gamestatus.USER_INPUT_WAIT_FRAME_COUNT
		return true
	} else {
		return false
	}
}

func (h *HouseScene) getEvent(id int) (*gamestatus.Event, int, error) {
	for i, e := range h.events {
		if e.Id == id {
			return e, i, nil
		}
	}
	return &gamestatus.Event{}, -1, errors.New("event not found")
}

func (h *HouseScene) isInputDirection(userAction gamestatus.UserAction) bool {
	return userAction == gamestatus.USER_ACTION_LEFT ||
		userAction == gamestatus.USER_ACTION_RIGHT ||
		userAction == gamestatus.USER_ACTION_UP ||
		userAction == gamestatus.USER_ACTION_DOWN
}

func (h *HouseScene) getNextPosition(currentX, currentY int, userAction gamestatus.UserAction) (int, int) {
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
	if nextX >= HOUSE_MAP_COLS {
		nextX = HOUSE_MAP_COLS - 1
	}
	if nextY >= HOUSE_MAP_ROWS {
		nextY = HOUSE_MAP_ROWS - 1
	}
	return nextX, nextY
}

func (h *HouseScene) existsMobCharacter(x, y int) (bool, int) {
	for i, m := range h.mobs {
		if m.Position.X == x && m.Position.Y == y {
			return true, i
		}
	}
	return false, 0
}

func (h *HouseScene) generateMobCharacter() []*MobCharacter {
	return []*MobCharacter{
		&MobCharacter{
			MobType: MOB_TYPE_BLACK_CAT,
			Position: util.MapPosition{
				X: 1,
				Y: 2,
			},
			Direction: util.DIRECTION_DOWN,
			EventId:   0,
		},
	}
}

func (h *HouseScene) generateEvents() []*gamestatus.Event {
	return []*gamestatus.Event{
		&gamestatus.Event{
			Id:          0,
			EventType:   gamestatus.EVENT_TYPE_MOB_TALK,
			TalkTexts:   []string{"シャー！"},
			NextEventId: -1,
			StoreId:     -1,
		},
	}
}

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
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const (
	TITLE_MAP_ROWS = 50
	TITLE_MAP_COLS = 50
)

type TitleSceneStatus int

const (
	TITLE_SCENE_STATUS_IDLE TitleSceneStatus = iota
	TITLE_SCENE_STATUS_MOVING
	TITLE_SCENE_PROCESS_EVENT
	TITLE_SCENE_CHANGE_SCENE
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
	gamepad               *GamePad
	sceneSwitcher         *SceneSwitcher
	mapImage              *ebiten.Image
	movableMap            map[util.MapPosition]bool
	currentPlayerPosition util.MapPosition
	nextPlayerPosition    util.MapPosition
	beforeImageDrawX      int
	beforeImageDrawY      int
	sceneStatus           TitleSceneStatus
	movingFrame           int
	mobs                  []*MobCharacter
	processEvent          *gamestatus.Event
	processEventTalkSeq   int
	events                []*gamestatus.Event
	stores                []*Store
	storeIndex            int
	userInputFrame        int
	audioPlayer           *audio.Player
}

func (t *TitleScene) Init(data *gamestatus.GameData) error {
	t.gameStateMsg = gamestatus.GAME_STATE_MSG_NONE
	t.gamepad = &GamePad{}
	if err := t.gamepad.Init(); err != nil {
		return err
	}

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
	t.movableMap, err = maps.LoadTitleMovableMap()
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
	t.mobs = append(t.mobs, t.generateMobCharacter()...)
	for _, m := range t.mobs {
		m.Init()
	}

	// イベント
	t.events = []*gamestatus.Event{}
	t.events = append(t.events, t.generateEvents()...)

	// ストア
	t.stores = []*Store{}
	t.stores = append(t.stores, t.generateStores()...)
	for _, s := range t.stores {
		s.Init()
	}

	// トーク用パネル
	t.talkPanel = &TalkPanel{}
	t.talkPanel.Init()

	// BGM
	audioStream, err := mp3.DecodeF32(bytes.NewReader(miyatamaAudio.Ragtime_mp3))
	if err != nil {
		return err
	}
	audioPlayer, err := data.AudioContext.NewPlayerF32(audioStream)
	if err != nil {
		return err
	}
	t.audioPlayer = audioPlayer
	t.audioPlayer.Play()

	// シーン切り替え
	t.sceneSwitcher = &SceneSwitcher{}
	t.sceneSwitcher.Init()
	return nil
}

func (t *TitleScene) Update(data *gamestatus.GameData) {
	t.gamepad.Update(data)
	if !t.audioPlayer.IsPlaying() {
		t.audioPlayer.SetPosition(time.Duration(0))
		t.audioPlayer.Play()
	}

	// キャラクタの移動
	switch t.sceneStatus {
	case TITLE_SCENE_STATUS_IDLE:
		{
			f := func() {
				if t.movePlayer(data) {
					return
				}

				if t.actionMobCharacter(data) {
					return
				}
			}
			f()
		}
	case TITLE_SCENE_PROCESS_EVENT:
		{
			f := func() {
				switch t.processEvent.EventType {
				case gamestatus.EVENT_TYPE_MOB_TALK:
					{
						// イベントを次に進める
						if !t.allowUserInput() {
							return
						}
						if data.UserAction != gamestatus.USER_ACTION_DECIDE {
							return
						}
						if len(t.processEvent.TalkTexts) > t.processEventTalkSeq+1 {
							t.processEventTalkSeq++
							t.talkPanel.SetText(t.processEvent.TalkTexts[t.processEventTalkSeq])
							return
						}
						if t.processEvent.NextEventId <= 0 {
							t.sceneStatus = TITLE_SCENE_STATUS_IDLE
							t.processEvent = nil
							return
						}
						evevnt, _, err := t.getEvent(t.processEvent.NextEventId)
						if err != nil {
							slog.Error(err.Error(),
								slog.Int("event id", t.processEvent.NextEventId),
							)
							return
						}
						t.processEvent = evevnt

						if evevnt.EventType == gamestatus.EVENT_TYPE_STORE {
							_, storeIndex, err := t.getStore(evevnt.StoreId)
							if err != nil {
								slog.Error(err.Error(),
									slog.Int("store id", evevnt.StoreId),
								)
								return
							}
							t.storeIndex = storeIndex
						}
					}
				case gamestatus.EVENT_TYPE_STORE:
					{
						t.stores[t.storeIndex].Update(data)
					}
				case gamestatus.EVENT_TYPE_CHANGE_SCENE:
					{
						t.sceneStatus = TITLE_SCENE_CHANGE_SCENE
						t.sceneSwitcher.StartIrisOut()
					}
				}
			}
			f()
		}
	case TITLE_SCENE_CHANGE_SCENE:
		{
			if t.sceneSwitcher.IsIdle() {
				t.audioPlayer.Close()
				t.gameStateMsg = gamestatus.GAME_STATE_MSG_HOUSE
			}
		}
	}
	for _, m := range t.mobs {
		m.Update(data)
	}
	t.player.Update(data)
	t.talkPanel.Update(data, t.gamepad.GetPadRect())
	t.sceneSwitcher.Update(data)
}

func (t *TitleScene) Draw(screen *ebiten.Image, data *gamestatus.GameData) {
	// マップの描画
	mapSx := (t.currentPlayerPosition.X*images.MAP_TILE_WIDTH + images.MAP_TILE_WIDTH/2)
	mapSy := (t.currentPlayerPosition.Y*images.MAP_TILE_WIDTH + images.MAP_TILE_WIDTH/2)
	// 移動中のフレーム判定
	if t.sceneStatus == TITLE_SCENE_STATUS_MOVING {
		deltaX := (t.nextPlayerPosition.X - t.currentPlayerPosition.X) * images.MAP_TILE_WIDTH
		deltaY := (t.nextPlayerPosition.Y - t.currentPlayerPosition.Y) * images.MAP_TILE_WIDTH
		deltaX = int(float32(deltaX) / float32(PLAYER_MOVING_FRAME_COUNT) * float32(t.movingFrame))
		deltaY = int(float32(deltaY) / float32(PLAYER_MOVING_FRAME_COUNT) * float32(t.movingFrame))
		mapSx += deltaX
		mapSy += deltaY
		t.movingFrame++
		if t.movingFrame > PLAYER_MOVING_FRAME_COUNT {
			t.currentPlayerPosition = t.nextPlayerPosition
			t.sceneStatus = TITLE_SCENE_STATUS_IDLE
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
	if t.processEvent != nil {
		switch t.processEvent.EventType {
		case gamestatus.EVENT_TYPE_MOB_TALK:
			{
				t.talkPanel.Draw(screen, data)
			}
		case gamestatus.EVENT_TYPE_STORE:
			{
				t.stores[t.storeIndex].Draw(screen, data)
			}
		}
	}

	// コントローラーの描画
	t.gamepad.Draw(screen, data)
	if t.userInputFrame > 0 {
		t.userInputFrame--
	}

	// シーン切り替えの描画
	if t.sceneStatus == TITLE_SCENE_CHANGE_SCENE {
		t.sceneSwitcher.Draw(screen, data)
	}
}

func (t *TitleScene) Msg() gamestatus.GameStateMsg {
	return t.gameStateMsg
}

func (t *TitleScene) movePlayer(data *gamestatus.GameData) bool {
	if !t.isInputDirection(data.UserAction) {
		return false
	}

	nextX, nextY := t.getNextPosition(t.currentPlayerPosition.X, t.currentPlayerPosition.Y, data.UserAction)
	key := util.MapPosition{X: nextX, Y: nextY}
	existsMobCharacter, _ := t.existsMobCharacter(nextX, nextY)
	movable := t.movableMap[key] && !existsMobCharacter
	if movable {
		t.nextPlayerPosition = util.MapPosition{
			X: nextX,
			Y: nextY,
		}
		t.sceneStatus = TITLE_SCENE_STATUS_MOVING
		t.movingFrame = 0
		t.player.SetUserAction(data.UserAction)
		return true
	} else {
		return false
	}
}

func (t *TitleScene) actionMobCharacter(data *gamestatus.GameData) bool {
	if !t.isInputDirection(data.UserAction) {
		return false
	}

	// 進行方向にモブが存在するか
	nextX, nextY := t.getNextPosition(t.currentPlayerPosition.X, t.currentPlayerPosition.Y, data.UserAction)
	existsMobCharacter, mobIndex := t.existsMobCharacter(nextX, nextY)
	if existsMobCharacter {
		t.player.SetUserAction(data.UserAction)
		t.sceneStatus = TITLE_SCENE_PROCESS_EVENT
		eventId := t.mobs[mobIndex].EventId
		event, _, err := t.getEvent(eventId)
		if err != nil {
			slog.Error(err.Error(),
				slog.Int("event id", eventId),
			)
			return false
		}
		t.processEvent = event
		t.processEventTalkSeq = 0
		if event.EventType == gamestatus.EVENT_TYPE_MOB_TALK {
			t.talkPanel.SetText(t.processEvent.TalkTexts[0])
		}
		slog.Info("TitleScene.actionMobCharacter",
			slog.Int("mob index", mobIndex),
		)
		return true
	} else {
		return false
	}
}

func (t *TitleScene) allowUserInput() bool {
	if t.userInputFrame <= 0 {
		t.userInputFrame = gamestatus.USER_INPUT_WAIT_FRAME_COUNT
		return true
	} else {
		return false
	}
}

func (t *TitleScene) getEvent(id int) (*gamestatus.Event, int, error) {
	for i, e := range t.events {
		if e.Id == id {
			return e, i, nil
		}
	}
	return &gamestatus.Event{}, -1, errors.New("event not found")

}

func (t *TitleScene) getStore(id int) (*Store, int, error) {
	for i, s := range t.stores {
		if s.Id == id {
			return s, i, nil
		}
	}
	return &Store{}, -1, errors.New("store not found")

}

func (t *TitleScene) isInputDirection(userAction gamestatus.UserAction) bool {
	return userAction == gamestatus.USER_ACTION_LEFT ||
		userAction == gamestatus.USER_ACTION_RIGHT ||
		userAction == gamestatus.USER_ACTION_UP ||
		userAction == gamestatus.USER_ACTION_DOWN
}

func (t *TitleScene) getNextPosition(currentX, currentY int, userAction gamestatus.UserAction) (int, int) {
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

func (t *TitleScene) generateMobCharacter() []*MobCharacter {
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
			EventId:   4,
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
		&MobCharacter{
			MobType: MOB_TYPE_NONE,
			Position: util.MapPosition{
				X: 26,
				Y: 5,
			},
			Direction: util.DIRECTION_DOWN,
			EventId:   3,
		},
	}
}

func (t *TitleScene) generateEvents() []*gamestatus.Event {
	return []*gamestatus.Event{
		&gamestatus.Event{
			Id:          0,
			EventType:   gamestatus.EVENT_TYPE_MOB_TALK,
			TalkTexts:   []string{"シャー！"},
			NextEventId: -1,
			StoreId:     -1,
		},
		&gamestatus.Event{
			Id:        1,
			EventType: gamestatus.EVENT_TYPE_MOB_TALK,
			TalkTexts: []string{
				"いらっしゃい、きょうは さかなが はいってるよ",
				"なににしますか？",
			},
			NextEventId: 2,
			StoreId:     -1,
		},
		&gamestatus.Event{
			Id:          2,
			EventType:   gamestatus.EVENT_TYPE_STORE,
			TalkTexts:   []string{},
			NextEventId: -1,
			StoreId:     1,
		},
		&gamestatus.Event{
			Id:          3,
			EventType:   gamestatus.EVENT_TYPE_CHANGE_SCENE,
			TalkTexts:   []string{},
			NextEventId: -1,
			EventTag:    gamestatus.EVENT_TAG_CHANGE_TO_HOUSE_SCENE,
		},
		&gamestatus.Event{
			Id:          4,
			EventType:   gamestatus.EVENT_TYPE_MOB_TALK,
			TalkTexts:   []string{"シャー！\nシャー！\nシャー！"},
			NextEventId: -1,
			StoreId:     -1,
		},
	}
}

func (t *TitleScene) generateStores() []*Store {
	return []*Store{
		&Store{
			Id: 1,
			StoreInfo: &StoreInfo{
				showpName: "みせや",
				baseCommand: []string{
					"かう",
					"店を出る",
				},
				items: []string{
					"しろきり",
					"くろきり",
					"あかきり",
				},
			},
		},
	}
}

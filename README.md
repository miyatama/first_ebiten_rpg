# outline

ebitengineでRPG実装できるかのチャレンジ

参考

+ [example/touch](https://github.com/hajimehoshi/ebiten/tree/main/examples/touch)

## のこりチャレンジ

+ アイテム購入的な演出してみる <- Next
+ 効果音つける
+ シーン相互切り替え
+ デバイス対応
  + デバイスのスクリーンサイズ取得 & 設定
  + 縦型 or 横型
+ デバイスから写真を取得する
+ アプリローカルに保存した写真を表示する
+ ローカルファイル保存(save & load)
  + もしかしたらデバイス依存
+ Androidで遅いのなんとかする
+ ユーザの上にくるマップ(看板とか)
+ パッケージ構成を変える
  + ui
  + usecase
  + domain
  + infra

## 対応済みチャレンジ

+ ~~テキスト領域を文字の大きさに合わせる~~
+ ~~音出す~~
+ ~~シーン変える~~
  + ~~アイリスイン・アイリスアウト(透過ロゴがギューってなるやつ~~
  + ~~[Masking](https://ebitengine.org/ja/examples/masking.html)~~
+ ~~コントローラー表示する~~
  + ~~十字 + 決定ボタン~~
+ ~~イベント周りの整備~~
+ ~~音楽~~
  +  ~~ループ再生~~
  +  ~~シーン切り替わりで停止~~
+ ~~文字の大きさをスクリーンサイズによって変動させる~~
+ ~~FPS表示~~

## for Android

android用のビルド

prepaire

```shell
# add path & get build command
PATH=${PATH}:C:\Program\ Files\Android\Android Studio\jbr\bin
go get github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile@v2.8.9
```

build

```shell
go run github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile bind -target android -javapkg com.miyatama.game_main -o ./mobile/android/ebitengine/game_main.aar ./mobile
```

emulator

```shell
emulator -list-avds
Pixel_6a
emulator -avd Pixel_6a
```

android apk install

```shell
cd mobile/android
./gradlew assembleDebug
adb install .\app\build\outputs\apk\debug\app-debug.apk
```

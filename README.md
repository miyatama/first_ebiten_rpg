# outline

ebitengineでRPG実装できるかのチャレンジ

参考

+ [example/touch](https://github.com/hajimehoshi/ebiten/tree/main/examples/touch)

## のこりチャレンジ

+ ~~音出す~~
+ アイテム購入的な演出してみる
+ シーン変える
+ デバイス対応
  + デバイスのスクリーンサイズ取得 & 設定
+ コントローラー表示する
  + 十字 + 決定ボタン
+ ローカルファイル保存(save & load)
  + もしかしたらデバイス依存

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
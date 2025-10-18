package gamestatus

type MobileInterface struct {
	// logger
	OutputDebugLog func(text string)
	OutputInfoLog  func(text string)
	OutputErrorLog func(text string)
}

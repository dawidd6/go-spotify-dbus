package spotify

// Listeners is a struct of the events we are listening for
type Listeners struct {
	OnMetadata       func(*Metadata)
	OnPlaybackStatus func(PlaybackStatus)
	OnServiceStart   func()
	OnServiceStop    func()
	OnError          func(error)
}

// NewListeners returns default listeners
func NewListeners() *Listeners {
	return &Listeners{
		OnMetadata:       func(*Metadata) {},
		OnPlaybackStatus: func(PlaybackStatus) {},
		OnServiceStart:   func() {},
		OnServiceStop:    func() {},
		OnError:          func(error) {},
	}
}

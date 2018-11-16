package playlists

// OnPlaylistItemEvent - an event that occurs when another playlist item is parsed
type OnPlaylistItemEvent func(item *PlaylistItem)

// IPlaylistParser - common playlist parser interface
type IPlaylistParser interface {
	Parse(data []byte) error
	AsyncParse(data []byte, onItem OnPlaylistItemEvent) error
	Guide() string
	Items() []*PlaylistItem
}

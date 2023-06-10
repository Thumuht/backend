package db

// stores all cache for thumuht
type AppCache struct {
	Sessions Cache[string, int] // map token -> userId
	// PostLike    Cache[int, int]    // map postid -> likes
	CommentLike Cache[int, int] // map commentid -> likes
	// PostView    Cache[int, int] // map postid -> view number
	Notifier Cache[int, chan *Message]
}

func NewAppCache() AppCache {
	return AppCache{
		Sessions: NewMapCache[string, int](),
		// PostLike:    NewMapCache[int, int](),
		CommentLike: NewMapCache[int, int](),
		// PostView:    NewMapCache[int, int](),
		Notifier: NewMapCacheCSP[int, chan *Message](),
	}
}

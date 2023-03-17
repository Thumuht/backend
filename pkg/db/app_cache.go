package db

// stores all cache for thumuht
type AppCache struct {
	Sessions    Cache[string, string] // map token -> user
	PostLike    Cache[int, int]       // map postid -> likes
	CommentLike Cache[int, int]       // map commentid -> likes
	PostView    Cache[int, int]       // map postid -> view number
}

func NewAppCache() AppCache {
	return AppCache{
		Sessions:    NewMapCache[string, string](),
		PostLike:    NewMapCache[int, int](),
		CommentLike: NewMapCache[int, int](),
		PostView:    NewMapCache[int, int](),
	}
}

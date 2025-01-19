// Code generated by goctl. DO NOT EDIT.
package types

type AddEmojiReq struct {
	UserID    string `header:"Beaver-User-Id"`
	FileUrl   string `json:"fileUrl"`
	Title     string `json:"title"`
	PackageID uint   `json:"packageId,optional"`
}

type AddEmojiRes struct {
}

type EmojiItem struct {
	EmojiID   uint   `json:"emojiId"`
	FileUrl   string `json:"fileUrl"`
	Title     string `json:"title"`
	PackageID *uint  `json:"packageId"`
}

type GetEmojisListReq struct {
	UserID string `header:"Beaver-User-Id"`
}

type GetEmojisListRes struct {
	Count int64       `json:"count"`
	List  []EmojiItem `json:"list"`
}

type UpdateFavoriteEmojiReq struct {
	UserID  string `header:"Beaver-User-Id"`
	EmojiID uint   `json:"emojiId"`
	Type    string `json:"type"` // "favorite" or "unfavorite"
}

type UpdateFavoriteEmojiRes struct {
}

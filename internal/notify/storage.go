package notify

// Storage describes persistent datastorage.
type Storage interface {
	GetChats() (ids []int64)
	SaveChats(ids []int64) error
}

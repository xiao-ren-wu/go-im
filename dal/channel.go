package dal

type ChannelMap interface {
	Add(Channel)
	Remove(id string)
	Get(id string) (Channel, bool)
	All() []Channel
}

type Channel interface {
}

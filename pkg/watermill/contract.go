package watermill

type Topic string

const (
	TopicMessageRaw     Topic = "message_raw"
	TopicMessageNoAlert Topic = "message_no_alert"
)

type Config struct {
	Debug bool
	Trace bool
}

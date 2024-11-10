package enums

import "time"

type PublishTime time.Duration

const (
	Latest        PublishTime = 0
	Earliest      PublishTime = -1
	WithinASecond PublishTime = PublishTime(time.Second * -1)
	WithinAnHour  PublishTime = PublishTime(time.Hour * -1)
	WithinADay    PublishTime = PublishTime(time.Hour * -24)
)

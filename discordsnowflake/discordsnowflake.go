package discordsnowflake

import "time"

func Snowflake2utc(sf int64) time.Time {
	timestamp := int64(((sf >> 22) + 1420070400000) / 1000)
	return time.Unix(timestamp, 0).UTC()
}

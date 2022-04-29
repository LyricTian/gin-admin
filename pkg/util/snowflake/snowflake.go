package snowflake

import (
	"os"
	"strconv"
	"time"

	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	settings := sonyflake.Settings{
		StartTime: time.Date(2021, 11, 8, 0, 0, 0, 0, time.UTC),
	}

	machineID := os.Getenv("MACHINE_ID")
	if machineID != "" {
		id, err := strconv.ParseUint(machineID, 10, 16)
		if err != nil {
			os.Stderr.WriteString("[warn] parse machine id failed: " + err.Error())
		} else {
			settings.MachineID = func() (uint16, error) {
				return uint16(id), nil
			}
		}
	}

	sf = sonyflake.NewSonyflake(settings)
	if sf == nil {
		panic("sonyflake not created")
	}
}

// Get unique id from Twitter's Snowflake
func MustID() uint64 {
	id, err := sf.NextID()
	if err == nil {
		return id
	}

	sleep := 1
	for {
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		id, err := sf.NextID()
		if err == nil {
			return id
		}
		sleep *= 2
	}
}

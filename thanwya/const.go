package thanwya

import (
	"time"
)

const (
	BaseURL                        = "http://natega.youm7.com/Home/GetResultStage1/"
	NumberOfGoRoutines             = 50
	InitialSeatNumber              = 1
	LastSeatNumber                 = 1000000
	NumberOfCuncurrentInserts      = 1000
	DatabaseName                   = "thanwya"
	DatabaseUser                   = "postgres"
	DatabaseHost                   = "localhost"
	DatabaseSSLMode                = "disable"
	NumberOfProgressPartitions     = 100
	MaxNumberOfArraySizeBeforeSave = 50000
	RefreshRate                    = 64 * time.Millisecond
)

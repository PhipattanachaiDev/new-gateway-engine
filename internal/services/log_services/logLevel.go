package logservices

import "github.com/rs/zerolog"

// GetLevelNameByInt takes an integer and returns the corresponding log level name
func GetLevelNameByInt(levelInt int) zerolog.Level {
	var level zerolog.Level
	switch levelInt {
	case 0:
		level = zerolog.DebugLevel
	case 1:
		level = zerolog.InfoLevel
	case 2:
		level = zerolog.WarnLevel
	case 3:
		level = zerolog.ErrorLevel
	case 4:
		level = zerolog.FatalLevel
	case 5:
		level = zerolog.PanicLevel
	case 6:
		level = zerolog.NoLevel
	case 7:
		level = zerolog.Disabled
	case -1:
		level = zerolog.TraceLevel
	default:
		level = zerolog.InfoLevel // ค่าเริ่มต้นเมื่อกำหนดค่าไม่ถูกต้อง
	}

	return level
}

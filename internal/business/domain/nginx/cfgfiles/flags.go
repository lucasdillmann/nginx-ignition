package cfgfiles

const (
	onFlag  = "on"
	offFlag = "off"
)

func flag(enabled bool, trueValue, falseValue string) string {
	if enabled {
		return trueValue
	}

	return falseValue
}

func statusFlag(enabled bool) string {
	return flag(enabled, onFlag, offFlag)
}

package valuerange

type ValueRange struct {
	Min int
	Max int
}

func New(minValue, maxValue int) *ValueRange {
	return &ValueRange{
		Min: minValue,
		Max: maxValue,
	}
}

func (vr *ValueRange) Contains(value int) bool {
	return value >= vr.Min && value <= vr.Max
}

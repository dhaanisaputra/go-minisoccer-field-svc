package constants

type FieldScheduleStatusName string
type FieldScheduleStatus int

const (
	Available FieldScheduleStatus = 100
	Booked    FieldScheduleStatus = 200

	AvailableString FieldScheduleStatusName = "Available"
	BookedString    FieldScheduleStatusName = "Booked"
)

var mapFieldScheduleStatusIntToStr = map[FieldScheduleStatus]FieldScheduleStatusName{
	Available: AvailableString,
	Booked:    BookedString,
}

var mapFieldScheduleStatusStrToInt = map[FieldScheduleStatusName]FieldScheduleStatus{
	AvailableString: Available,
	BookedString:    Booked,
}

func (f FieldScheduleStatus) GetStatusString() FieldScheduleStatusName {
	return mapFieldScheduleStatusIntToStr[f]
}

func (f FieldScheduleStatusName) GetStatusInt() FieldScheduleStatus {
	return mapFieldScheduleStatusStrToInt[f]
}

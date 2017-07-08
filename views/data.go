package views

const (
	LevelError   = "danger"
	LevelWarning = "warning"
	LevelInfo    = "info"
	LevelSuccess = "success"
)

var (
	AlertGeneric *Alert = AlertError("Sorry! something went wrong")
)

type alertLevel string

type Alert struct {
	Level   alertLevel
	Message string
}

type Data struct {
	Alert *Alert
	Yield interface{}
}

func AlertError(msg string) *Alert {
	return &Alert{LevelError, msg}
}

func AlertWarning(msg string) *Alert {
	return &Alert{LevelWarning, msg}
}
func AlertInfo(msg string) *Alert {
	return &Alert{LevelInfo, msg}
}
func AlertSuccess(msg string) *Alert {
	return &Alert{LevelSuccess, msg}
}

func (d *Data) AddAlert(alert *Alert) {
	d.Alert = alert
}

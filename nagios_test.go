package nagios

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNagiosStatus_Aggregate(t *testing.T) {
	Convey("Aggregates statuses together", t, func() {

		otherStatuses := []*NagiosStatus{
			&NagiosStatus{"ok", NAGIOS_OK},
			&NagiosStatus{"Not so bad", NAGIOS_WARNING},
		}

		Convey("Picks the worst status", func() {
			status := &NagiosStatus{"Uh oh", NAGIOS_CRITICAL}
			status.Aggregate(otherStatuses)

			So(status.Value, ShouldEqual, NAGIOS_CRITICAL)
		})

		Convey("Aggregates the messages", func() {
			status := &NagiosStatus{"Uh oh", NAGIOS_CRITICAL}
			status.Aggregate(otherStatuses)

			So(status.Message, ShouldEqual, "Uh oh - ok - Not so bad")
		})

		Convey("Handles an empty slice", func() {
			status := &NagiosStatus{"Uh oh", NAGIOS_CRITICAL}
			status.Aggregate([]*NagiosStatus{})

			So(status.Value, ShouldEqual, NAGIOS_CRITICAL)
			So(status.Message, ShouldEqual, "Uh oh")
		})

	})
}

func TestValMessages(t *testing.T) {
	Convey("Maps the correct strings to values", t, func() {
		So(valMessages[NAGIOS_UNKNOWN], ShouldEqual, "UNKNOWN:")
		So(valMessages[NAGIOS_CRITICAL], ShouldEqual, "CRITICAL:")
		So(valMessages[NAGIOS_WARNING], ShouldEqual, "WARNING:")
		So(valMessages[NAGIOS_OK], ShouldEqual, "OK:")
	})
}

func TestConstructedNagiosMessage(t *testing.T) {
	Convey("Constructs a Nagios message without performance data", t, func() {
		statusUnknown := &NagiosStatus{"Shrug dunno", NAGIOS_UNKNOWN}
		So(statusUnknown.constructedNagiosMessage(), ShouldEqual, "UNKNOWN: Shrug dunno")

		statusCritical := &NagiosStatus{"Uh oh", NAGIOS_CRITICAL}
		So(statusCritical.constructedNagiosMessage(), ShouldEqual, "CRITICAL: Uh oh")

		statusWarning := &NagiosStatus{"Not so bad", NAGIOS_WARNING}
		So(statusWarning.constructedNagiosMessage(), ShouldEqual, "WARNING: Not so bad")

		statusOK := &NagiosStatus{"ok", NAGIOS_OK}
		So(statusOK.constructedNagiosMessage(), ShouldEqual, "OK: ok")
	})

	Convey("Constructs a Nagios message with performance data", t, func() {
		statusUnknown := &NagiosStatus{"Shrug dunno", NAGIOS_UNKNOWN}
		perfdata1 := NagiosPerformanceVal{"metric", "1234", "ms", "12", "3400", "0", "99999"}
		statusUnknownPerf := &NagiosStatusWithPerformanceData{statusUnknown, perfdata1}
		So(statusUnknownPerf.constructedNagiosMessage(), ShouldEqual, "UNKNOWN: Shrug dunno | 'metric'=1234ms;12;3400;0;99999")

		statusCritical := &NagiosStatus{"Uh oh", NAGIOS_CRITICAL}
		perfdata2 := NagiosPerformanceVal{"metric", "1234", "ms", "12", "3400", "", ""}
		statusCriticalPerf := &NagiosStatusWithPerformanceData{statusCritical, perfdata2}
		So(statusCriticalPerf.constructedNagiosMessage(), ShouldEqual, "CRITICAL: Uh oh | 'metric'=1234ms;12;3400;;")

		statusWarning := &NagiosStatus{"Not so bad", NAGIOS_WARNING}
		perfdata3 := NagiosPerformanceVal{"metric", "1234", "ms", "", "", "0", "99999"}
		statusWarningPerf := &NagiosStatusWithPerformanceData{statusWarning, perfdata3}
		So(statusWarningPerf.constructedNagiosMessage(), ShouldEqual, "WARNING: Not so bad | 'metric'=1234ms;;;0;99999")

		statusOK := &NagiosStatus{"ok", NAGIOS_OK}
		perfdata4 := NagiosPerformanceVal{"metric", "1234", "", "12", "3400", "0", "99999"}
		statusOKPerf := &NagiosStatusWithPerformanceData{statusOK, perfdata4}
		So(statusOKPerf.constructedNagiosMessage(), ShouldEqual, "OK: ok | 'metric'=1234;12;3400;0;99999")
	})
}

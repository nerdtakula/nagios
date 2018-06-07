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
		status_unknown := &NagiosStatus{"Shrug dunno", NAGIOS_UNKNOWN}
		So(status_unknown.constructedNagiosMessage(), ShouldEqual, "UNKNOWN: Shrug dunno")

		status_critical := &NagiosStatus{"Uh oh", NAGIOS_CRITICAL}
		So(status_critical.constructedNagiosMessage(), ShouldEqual, "CRITICAL: Uh oh")

		status_warning := &NagiosStatus{"Not so bad", NAGIOS_WARNING}
		So(status_warning.constructedNagiosMessage(), ShouldEqual, "WARNING: Not so bad")

		status_ok := &NagiosStatus{"ok", NAGIOS_OK}
		So(status_ok.constructedNagiosMessage(), ShouldEqual, "OK: ok")
	})

	Convey("Constructs a Nagios message with performance data", t, func() {
		status_unknown := &NagiosStatus{"Shrug dunno", NAGIOS_UNKNOWN}
		perfdata1 := NagiosPerformanceVal{"metric", "1234", "ms", "12", "3400", "0", "99999"}
		status_unknown_perf := &NagiosStatusWithPerformanceData{status_unknown, perfdata1}
		So(status_unknown_perf.constructedNagiosMessage(), ShouldEqual, "UNKNOWN: Shrug dunno | 'metric'=1234ms;12;3400;0;99999")

		status_critical := &NagiosStatus{"Uh oh", NAGIOS_CRITICAL}
		perfdata2 := NagiosPerformanceVal{"metric", "1234", "ms", "12", "3400", "", ""}
		status_critical_perf := &NagiosStatusWithPerformanceData{status_critical, perfdata2}
		So(status_critical_perf.constructedNagiosMessage(), ShouldEqual, "CRITICAL: Uh oh | 'metric'=1234ms;12;3400;;")

		status_warning := &NagiosStatus{"Not so bad", NAGIOS_WARNING}
		perfdata3 := NagiosPerformanceVal{"metric", "1234", "ms", "", "", "0", "99999"}
		status_warning_perf := &NagiosStatusWithPerformanceData{status_warning, perfdata3}
		So(status_warning_perf.constructedNagiosMessage(), ShouldEqual, "WARNING: Not so bad | 'metric'=1234ms;;;0;99999")

		status_ok := &NagiosStatus{"ok", NAGIOS_OK}
		perfdata4 := NagiosPerformanceVal{"metric", "1234", "", "12", "3400", "0", "99999"}
		status_ok_perf := &NagiosStatusWithPerformanceData{status_ok, perfdata4}
		So(status_ok_perf.constructedNagiosMessage(), ShouldEqual, "OK: ok | 'metric'=1234;12;3400;0;99999")
	})
}

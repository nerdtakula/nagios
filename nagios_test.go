package nagios

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestState_String(t *testing.T) {
	Convey("Maps the correct strings to values", t, func() {
		So(stateStrings[STATE_UNKNOWN], ShouldEqual, "UNKNOWN")
		So(stateStrings[STATE_CRITICAL], ShouldEqual, "CRITICAL")
		So(stateStrings[STATE_WARNING], ShouldEqual, "WARNING")
		So(stateStrings[STATE_OK], ShouldEqual, "OK")
	})
}

func TestStatus_String(t *testing.T) {
	Convey("Maps the correct strings to values", t, func() {
	})
}

func TestStatusWithPerformanceData_String(t *testing.T) {
	Convey("Maps the correct strings to values", t, func() {
		// Convey("Aggregates basic statuses together", func() {}

		Convey("String when no perfdata available", func() {
			status1 := &StatusWithPerformanceData{Status: New()}
			status1.Message = "FooBar"
			So(status1.Int(), ShouldEqual, STATE_OK)
			So(status1.Message, ShouldEqual, "FooBar")
			So(status1.String(), ShouldEqual, "OK: FooBar")

			status2 := New().WithPerfdata()
			status2.Message = "FooBar"
			So(status2.Int(), ShouldEqual, STATE_OK)
			So(status2.Message, ShouldEqual, "FooBar")
			So(status2.String(), ShouldEqual, "OK: FooBar")
		})
	})
}

func Test_Aggregate(t *testing.T) {
	Convey("Aggregates statuses together", t, func() {
		Convey("Aggregates basic statuses together", func() {
			statuses := []*Status{
				&Status{"ok", STATE_OK},
				&Status{"Not so bad", STATE_WARNING},
			}

			s, _ := Aggregate(statuses...)
			So(s.State.Int(), ShouldEqual, STATE_WARNING)
			So(s.Message, ShouldEqual, "ok - Not so bad")
			So(s.String(), ShouldEqual, "WARNING: ok - Not so bad")
		})

		Convey("Aggregates statuses with perfdata together", func() {
			statuses := []*StatusWithPerformanceData{
				&StatusWithPerformanceData{
					Status: &Status{
						Message: "ok",
						State:   STATE_OK,
					},
					Perfdata: []Perfdata{
						Perfdata{},
					},
				},
				&StatusWithPerformanceData{
					Status: &Status{
						Message: "Not so bad",
						State:   STATE_WARNING,
					},
					Perfdata: []Perfdata{
						Perfdata{},
					},
				},
				&StatusWithPerformanceData{
					Status: &Status{
						Message: "unknown",
						State:   STATE_UNKNOWN,
					},
					Perfdata: []Perfdata{
						Perfdata{},
					},
				},
			}

			s, _ := AggregateWithPerfdata(statuses...)
			So(s.State.Int(), ShouldEqual, STATE_UNKNOWN)
			So(s.Message, ShouldEqual, "ok - Not so bad - unknown")
			So(s.String(), ShouldEqual, "UNKOWN: ok - Not so bad - unknown")
		})

		// Convey("Aggregates basic statuses and statuses with perfdata together", func() {
		// 	statuses := []StatusType{
		// 		&Status{"ok", STATE_OK},
		// 		&StatusWithPerformanceData{
		// 			Status: &Status{
		// 				Message: "Not so bad",
		// 				State:   STATE_WARNING,
		// 			},
		// 			Perfdata: nil,
		// 		},
		// 	}

		// 	s, _ := Aggregate(statuses...)
		// 	So(s.State.Int(), ShouldEqual, STATE_WARNING)
		// 	So(s.Message, ShouldEqual, "ok - Not so bad")
		// 	So(s.String(), ShouldEqual, "WARNING: ok - Not so bad")
		// })
	})
}

func TestStatus_Aggregate(t *testing.T) {
	Convey("Aggregates statuses together into existing status", t, func() {
		otherStatuses := []*Status{
			&Status{"ok", STATE_OK},
			&Status{"Not so bad", STATE_WARNING},
		}

		Convey("Picks the worst status", func() {
			status := &Status{"Uh oh", STATE_CRITICAL}
			status.Aggregate(otherStatuses...)

			So(status.State, ShouldEqual, STATE_CRITICAL)
		})

		Convey("Aggregates the messages", func() {
			status := &Status{"Uh oh", STATE_CRITICAL}
			status.Aggregate(otherStatuses...)

			So(status.Message, ShouldEqual, "Uh oh - ok - Not so bad")
		})

		Convey("Handles an empty slice", func() {
			status := &Status{"Uh oh", STATE_CRITICAL}
			emptySlice := make([]*Status, 0)
			status.Aggregate(emptySlice...)

			So(status.State, ShouldEqual, STATE_CRITICAL)
			So(status.Message, ShouldEqual, "Uh oh")
		})
	})
}

func TestStatusWithPerformanceData_Aggregate(t *testing.T) {}

func TestStatusType_Aggregate(t *testing.T) {}

func TestPerfdata(t *testing.T) {}

func TestConstructedNagiosMessage(t *testing.T) {
	Convey("Constructs a Nagios message without performance data", t, func() {
		statusUnknown := &Status{"Shrug dunno", STATE_UNKNOWN}
		So(statusUnknown.String(), ShouldEqual, "UNKNOWN: Shrug dunno")

		statusCritical := &Status{"Uh oh", STATE_CRITICAL}
		So(statusCritical.String(), ShouldEqual, "CRITICAL: Uh oh")

		statusWarning := &Status{"Not so bad", STATE_WARNING}
		So(statusWarning.String(), ShouldEqual, "WARNING: Not so bad")

		statusOK := &Status{"ok", STATE_OK}
		So(statusOK.String(), ShouldEqual, "OK: ok")
	})

	Convey("Constructs a Nagios message with performance data", t, func() {
		statusUnknown := &Status{"Shrug dunno", STATE_UNKNOWN}

		perfdata1 := Perfdata{Label: "metric", Value: "1234", Uom: "ms", WarnThreshold: "12", CritThreshold: "3400", MinValue: "0", MaxValue: "99999"}
		statusUnknownPerf := &StatusWithPerformanceData{statusUnknown, []Perfdata{perfdata1}}
		So(statusUnknownPerf.String(), ShouldEqual, "UNKNOWN: Shrug dunno | 'metric'=1234ms;12;3400;0;99999")

		statusCritical := &Status{"Uh oh", STATE_CRITICAL}
		perfdata2 := Perfdata{Label: "metric", Value: "1234", Uom: "ms", WarnThreshold: "12", CritThreshold: "3400", MinValue: "", MaxValue: ""}
		statusCriticalPerf := &StatusWithPerformanceData{statusCritical, []Perfdata{perfdata2}}
		So(statusCriticalPerf.String(), ShouldEqual, "CRITICAL: Uh oh | 'metric'=1234ms;12;3400;;")

		statusWarning := &Status{"Not so bad", STATE_WARNING}
		perfdata3 := Perfdata{Label: "metric", Value: "1234", Uom: "ms", WarnThreshold: "", CritThreshold: "", MinValue: "0", MaxValue: "99999"}
		statusWarningPerf := &StatusWithPerformanceData{statusWarning, []Perfdata{perfdata3}}
		So(statusWarningPerf.String(), ShouldEqual, "WARNING: Not so bad | 'metric'=1234ms;;;0;99999")

		statusOK := &Status{"ok", STATE_OK}
		perfdata4 := Perfdata{Label: "metric", Value: "1234", Uom: "", WarnThreshold: "12", CritThreshold: "3400", MinValue: "0", MaxValue: "99999"}
		statusOKPerf := &StatusWithPerformanceData{statusOK, []Perfdata{perfdata4}}
		So(statusOKPerf.String(), ShouldEqual, "OK: ok | 'metric'=1234;12;3400;0;99999")
	})
}

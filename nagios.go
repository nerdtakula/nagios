package nagios

import (
	"fmt"
	"os"
)

type NagiosStatusVal int

// The values with which a Nagios check can exit
const (
	NAGIOS_OK NagiosStatusVal = iota
	NAGIOS_WARNING
	NAGIOS_CRITICAL
	NAGIOS_UNKNOWN
)

// Maps the NagiosStatusVal entries to output strings
var (
	valMessages = []string{
		"OK:",
		"WARNING:",
		"CRITICAL:",
		"UNKNOWN:",
	}
)

//--------------------------------------------------------------
// A type representing a Nagios check status. The Value is a the exit code
// expected for the check and the Message is the specific output string.
type NagiosStatus struct {
	Message string
	Value   NagiosStatusVal
}

// Take a bunch of NagiosStatus pointers and find the highest value, then
// combine all the messages. Things win in the order of highest to lowest.
func (status *NagiosStatus) Aggregate(otherStatuses []*NagiosStatus) {
	for _, s := range otherStatuses {
		if status.Value < s.Value {
			status.Value = s.Value
		}

		status.Message += " - " + s.Message
	}
}

// Construct the Nagios message
func (status *NagiosStatus) constructedNagiosMessage() string {
	return valMessages[status.Value] + " " + status.Message
}

// NagiosStatus: Issue a Nagios message to stdout and exit with appropriate Nagios code
func (status *NagiosStatus) NagiosExit() {
	fmt.Fprintln(os.Stdout, status.constructedNagiosMessage())
	os.Exit(int(status.Value))
}

//--------------------------------------------------------------
// A type representing a Nagios performance data value.
// https://nagios-plugins.org/doc/guidelines.html#AEN200
// http://docs.pnp4nagios.org/pnp-0.6/about#system_requirements
type NagiosPerformanceVal struct {
	Label         string
	Value         string
	Uom           string
	WarnThreshold string
	CritThreshold string
	MinValue      string
	MaxValue      string
}

//--------------------------------------------------------------
// A type representing a Nagios check status and performance data.
type NagiosStatusWithPerformanceData struct {
	*NagiosStatus
	Perfdata NagiosPerformanceVal
}

// Construct the Nagios message with performance data
func (status *NagiosStatusWithPerformanceData) constructedNagiosMessage() string {
	msg := fmt.Sprintf("%s %s | '%s'=%s%s;%s;%s;%s;%s",
		valMessages[status.Value],
		status.Message,
		status.Perfdata.Label,
		status.Perfdata.Value,
		status.Perfdata.Uom,
		status.Perfdata.WarnThreshold,
		status.Perfdata.CritThreshold,
		status.Perfdata.MinValue,
		status.Perfdata.MaxValue)
	return msg
}

// Issue a Nagios message (with performance data) to stdout and exit with appropriate Nagios code
func (status *NagiosStatusWithPerformanceData) NagiosExit() {
	fmt.Fprintln(os.Stdout, status.constructedNagiosMessage())
	os.Exit(int(status.Value))
}

//--------------------------------------------------------------

// Exit with an UNKNOWN status and appropriate message
func Unknown(output string) {
	ExitWithStatus(&NagiosStatus{output, NAGIOS_UNKNOWN})
}

// Exit with an CRITICAL status and appropriate message
func Critical(err error) {
	ExitWithStatus(&NagiosStatus{err.Error(), NAGIOS_CRITICAL})
}

// Exit with an WARNING status and appropriate message
func Warning(output string) {
	ExitWithStatus(&NagiosStatus{output, NAGIOS_WARNING})
}

// Exit with an OK status and appropriate message
func Ok(output string) {
	ExitWithStatus(&NagiosStatus{output, NAGIOS_OK})
}

// Exit with a particular NagiosStatus
func ExitWithStatus(status *NagiosStatus) {
	status.NagiosExit()
}

package nagios

import (
	"fmt"
	"os"
	"strings"
)

// State is a representation of the available exit codes supported by the
// Nagios spec.
type State int

// String returns a string representation of the State.
func (s State) String() string { return stateStrings[s] }

// Int returns the type casted value
func (s State) Int() int { return int(s) }

const (
	// The values with which a Nagios check can exit.
	STATE_OK State = iota
	STATE_WARNING
	STATE_CRITICAL
	STATE_UNKNOWN
)

var (
	// Maps the State entries to output strings.
	stateStrings = map[State]string{
		STATE_OK:       "OK",
		STATE_WARNING:  "WARNING",
		STATE_CRITICAL: "CRITICAL",
		STATE_UNKNOWN:  "UNKNOWN",
	}
)

// StatusType is a basic interface
type StatusType interface {
	String() string
	Int() int
	Aggregate(statuses ...StatusType)
	Exit()
}

// Status is a type representing a Nagios check status.
type Status struct {
	Message string
	State   State
}

// String returns a string representation of the Status.
func (s Status) String() string { return fmt.Sprintf("%s: %s", s.State, s.Message) }

// Int returns the State integer value
func (s Status) Int() int { return s.State.Int() }

// Aggregate takes multiple Status structs and combines them into this struct.
// Uses the highest State value and combines all the messages
func (s *Status) Aggregate(statuses ...*Status) {
	for _, o := range statuses {
		if o.State > s.State {
			s.State = o.State
		}
		s.Message += " - " + o.Message
	}
}

// Exit is designed to be called via the `defer` keyword. Prints a Nagios
// message to STDOUT and exits with appropriate Nagios code.
func (s Status) Exit() {
	fmt.Fprintln(os.Stdout, s)
	os.Exit(s.State.Int())
}

// Aggregate takes multiple Status structs and combines them. Uses the highest
// State value and combines all the messages.
func Aggregate(statuses ...*Status) (*Status, error) {
	if len(statuses) == 0 {
		return nil, fmt.Errorf("no statuses provided to aggregate")
	}

	t := &Status{}
	msgs := make([]string, len(statuses))

	for i, s := range statuses {
		if s.State > t.State {
			t.State = s.State
		}
		msgs[i] = s.Message
	}

	t.Message = strings.Join(msgs, " - ")
	return t, nil
}

// AggregateWithPerfdata takes multiple Status structs with Performance data
// and combines them. Uses the highest State value and combines all the
// messages.
func AggregateWithPerfdata(statuses ...*StatusWithPerformanceData) (*StatusWithPerformanceData, error) {
	if len(statuses) == 0 {
		return nil, fmt.Errorf("no statuses provided to aggregate")
	}

	t := &StatusWithPerformanceData{Status: &Status{}, Perfdata: make([]Perfdata, 0)}
	msgs := make([]string, len(statuses))

	for i, s := range statuses {
		if s.State > t.State {
			t.State = s.State
		}
		msgs[i] = s.Message
		for _, p := range s.Perfdata {
			t.Perfdata = append(t.Perfdata, p)
		}
	}

	t.Message = strings.Join(msgs, " - ")
	return t, nil
}

// Perfdata is a type representing the Nagios performance data structure.
// > https://nagios-plugins.org/doc/guidelines.html#AEN200
// >    'label'=value[UOM];[warn];[crit];[min];[max]
// > http://docs.pnp4nagios.org/pnp-0.6/about#system_requirements
type Perfdata struct {
	Label         string
	Value         string
	Uom           string
	WarnThreshold string
	CritThreshold string
	MinValue      string
	MaxValue      string
}

func (p Perfdata) String() string {
	return fmt.Sprintf("'%s'=%s%s;%s;%s;%s;%s",
		p.Label,
		p.Value,
		p.Uom,
		p.WarnThreshold,
		p.CritThreshold,
		p.MinValue,
		p.MaxValue)
}

// StatusWithPerformanceData provides a type representation of a Nagios check
// status containing performance data.
type StatusWithPerformanceData struct {
	*Status
	Perfdata []Perfdata
}

func (s StatusWithPerformanceData) String() string {
	if s.Perfdata == nil || len(s.Perfdata) == 0 {
		return fmt.Sprintf("%s: %s", s.State, s.Message)
	}

	pd := make([]string, len(s.Perfdata))
	for i, p := range s.Perfdata {
		pd[i] = p.String()
	}

	pdStr := strings.Join(pd, "; ")
	return fmt.Sprintf("%s: %s | %s", s.State, s.Message, pdStr)
}

// Int returns the value of the state.
func (s StatusWithPerformanceData) Int() int { return s.State.Int() }

// Aggregate takes multiple Status with Performance data structs and combines
// them into this struct. Uses the highest State value and combines all the
// messages.
func (s *StatusWithPerformanceData) Aggregate(statuses ...*StatusWithPerformanceData) {
	for _, o := range statuses {
		if o.State > s.State {
			s.State = o.State
		}
		s.Message += " - " + o.Message
		for _, p := range o.Perfdata {
			s.Perfdata = append(s.Perfdata, p)
		}
	}
}

// Exit is designed to be called via the `defer` keyword. Prints a Nagios
// message to STDOUT and exits with appropriate Nagios code.
func (s StatusWithPerformanceData) Exit() {
	fmt.Fprintln(os.Stdout, s)
	os.Exit(s.State.Int())
}

// Unknown provides a quick way to exit with an UNKNOWN state and appropriate
// message.
func Unknown(output string) {
	ExitWithStatus(&Status{output, STATE_UNKNOWN})
}

// Critical provides a quick way to exit with an CRITICAL state and
// appropriate message.
func Critical(err error) {
	ExitWithStatus(&Status{err.Error(), STATE_CRITICAL})
}

// Warning provides a quick way to exit with an WARNING state and appropriate
// message.
func Warning(output string) {
	ExitWithStatus(&Status{output, STATE_WARNING})
}

// OK provides a quick way to exit with an OK state and appropriate message.
func OK(output string) {
	ExitWithStatus(&Status{output, STATE_OK})
}

// ExitWithStatus ...
func ExitWithStatus(status *Status) {
	status.Exit()
}

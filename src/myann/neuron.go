package myann

import "fmt"

type message struct {
	message string
	payload interface{}
	reply   chan string
}

/* object to communicate with neuron with */
type Neuron struct {
	command chan message
	input   chan Data
}

/**
 * Neuron functions
 */
func (neuron *Neuron) AddOutput(conn *Connection) {
	ch := make(chan string)
	neuron.command <- message{message: AddOutput, payload: conn, reply: ch}
	// synchronise
	_ = <-ch
}

// Add a neuron as an input with weight
func (neuron *Neuron) AddInput(input *Neuron, w Data) {
	input.AddOutput(&Connection{c: neuron.input, weight: w})
	ch := make(chan string)
	neuron.command <- message{message: AddInput, payload: input, reply: ch}
	_ = <-ch
}

/**
 * log functions
 */
func (neuron *Neuron) SetLogLevel(level int) {
	ch := make(chan string)
	neuron.command <- message{message: SetLogLevel, payload: level, reply: ch}
	_ = <-ch
}

func (neuron *neuron) info(format string, a ...interface{}) {
	real_format := fmt.Sprintf("\t\t%v: %v\n", neuron.name, format)
	log_message(neuron.log_level, 1, real_format, a...)
}

func (neuron *neuron) debug(format string, a ...interface{}) {
	real_format := fmt.Sprintf("\t\t%v: %v\n", neuron.name, format)
	log_message(neuron.log_level, 2, real_format, a...)
}

func (neuron *neuron) trace(format string, a ...interface{}) {
	real_format := fmt.Sprintf("\t\t%v: %v\n", neuron.name, format)
	log_message(neuron.log_level, 3, real_format, a...)
}

/*

	input 1 -\-/---neuron 1
			  X
	input 2 -/-\---neuron 2

*/

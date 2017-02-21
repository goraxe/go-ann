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
	realFormat := fmt.Sprintf("\t\t%v: %v\n", neuron.name, format)
	logMessage(neuron.logLevel, 1, realFormat, a...)
}

func (neuron *neuron) debug(format string, a ...interface{}) {
	realFormat := fmt.Sprintf("\t\t%v: %v\n", neuron.name, format)
	logMessage(neuron.logLevel, 2, realFormat, a...)
}

func (neuron *neuron) trace(format string, a ...interface{}) {
	realFormat := fmt.Sprintf("\t\t%v: %v\n", neuron.name, format)
	logMessage(neuron.logLevel, 3, realFormat, a...)
}

/*

	input 1 -\-/---neuron 1
			  X
	input 2 -/-\---neuron 2

*/

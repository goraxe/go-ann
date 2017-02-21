package myann

import (
	"container/list"
	"fmt"
)

const NCPU = 4
const Reset = "Reset"
const AddInput = "AddInput"
const AddOutput = "AddOutput"
const SetLogLevel = "SetLogLevel"

type Data float64

/* the real neuron inside a goroutine */
type neuron struct {
	name          string
	input         chan Data
	outputs       *list.List
	command       chan message
	inputs        int
	inputsWaiting int
	accumulator   Data
	threshold     Data
	logLevel      int
}

/**
 * construct functions
 */
func MakePattern(a ...Data) *list.List {
	pattern := new(list.List)
	for _, v := range a {
		pattern.PushBack(v)
	}
	return pattern
}

func CompareList(v1 *list.List, v2 *list.List) bool {
	if v1.Len() != v2.Len() {
		return false
	}
	e1 := v1.Front()
	e2 := v2.Front()
	for i := 0; i < v1.Len(); i++ {
		if e1.Value != e2.Value {
			return false
		}
		e1 = e1.Next()
		e2 = e2.Next()
	}
	return true
}

func PrintList(l *list.List) {
	e := l.Front()
	for e != nil {
		fmt.Printf("%v ", e.Value)
		e = e.Next()
	}
	fmt.Print("\n")
}

func (v *ConnectionArray) Len() int {
	return len(*v)
}

func neuronLoop(neuron *neuron) {
	for {
		select {
		case v := <-neuron.input:
			neuron.trace("received %v", v)
			neuron.inputs--
			neuron.trace("neuron.inputs remaining is %v", neuron.inputs)
			neuron.accumulator += v
			neuron.trace("neuron.accumulator is %v", neuron.accumulator)
			if neuron.inputs == 0 {
				var fire_val Data
				if neuron.accumulator >= neuron.threshold {
					fire_val = 1.0
				} else {
					fire_val = 0.0
				}
				neuron.trace("is fireing value %v", fire_val)
				for e := neuron.outputs.Front(); e != nil; e = e.Next() {
					output := e.Value.(*Connection)
					output.c <- output.weight * fire_val
				}
			}
		case command := <-neuron.command:
			switch command.message {
			case AddOutput:
				neuron.trace("adding output %v", command.payload)
				neuron.outputs.PushBack(command.payload)
				command.reply <- "Ok"
			case AddInput:
				neuron.inputsWaiting++
				neuron.inputs++
				neuron.trace("input_wating is %v inputs is %v", neuron.inputsWaiting, neuron.inputs)
				command.reply <- "Ok"
			case Reset:
				neuron.inputs = neuron.inputsWaiting
				neuron.accumulator = 0
				command.reply <- "Ok"
			case SetLogLevel:
				neuron.logLevel = command.payload.(int)
				command.reply <- "Ok"
			default:
				if command.reply != nil {
					command.reply <- fmt.Sprintf("%v Error: Unknown command: %v", neuron.name, command.message)
				} else {
					fmt.Printf("%v Error: Unknown command: %v", neuron.name, command.message)
				}
			}
		} //select
	} // for
}

func logMessage(logLevel int, message_level int, format string, a ...interface{}) {
	if logLevel >= message_level {
		fmt.Printf(format, a...)
	}
}

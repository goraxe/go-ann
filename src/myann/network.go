package myann

import (
	"container/list"
	"errors"
)

// Connection between two neuron
type Connection struct {
	c      chan Data
	weight Data
}
type ConnectionArray []*Connection

type Network struct {
	inputs   ConnectionArray         // for presenting a pattern to the network
	outputs  ConnectionArray         // for collecting the results
	command  map[string]chan message // for sending a command to all neurons
	logLevel int
}

/**
 * network functions
 */
func newConnection(w Data) *Connection {
	return &Connection{c: make(chan Data, NCPU), weight: w}
}

func newConnectionArray(connections int) *ConnectionArray {
	ca := make(ConnectionArray, connections)
	/*	for i := 0; i < connections; i++ {
		ca[i] = newConnection(0)
	}*/
	return &ca
}

func CreateNetwork(inputs int, outputs int) *Network {
	n := &Network{command: make(map[string]chan message), logLevel: 1}
	n.info("Creating network %v inputs, %v outputs\n", inputs, outputs)
	n.inputs = *newConnectionArray(inputs)
	n.outputs = *newConnectionArray(outputs)

	for i := 0; i < outputs; i++ {
		n.outputs[i] = newConnection(1)
	}
	return n
}

func (network *Network) CreateNeuron(name string, threshold Data) *Neuron {
	network.info("Create %v, ", name)
	network.command[name] = make(chan message, NCPU)

	n := &neuron{name: name, input: make(chan Data, NCPU),
		outputs: new(list.List), command: network.command[name],
		threshold: threshold}

	network.info("Add Neuron and start neuronLoop...\n")
	defer network.AddNeuron(n)
	return &Neuron{command: network.command[name], input: n.input}
}

func (network *Network) AddNeuron(neuron *neuron) {
	go neuronLoop(neuron)
}

// Run Pattern
func (network *Network) RunPattern(input *list.List) (output *list.List, err error) {

	network.debug("resetting network\n")
	// send reset
	ch := make(chan string)
	for _, v := range network.command {
		v <- message{message: Reset, reply: ch}
		// wait for response
		_ = <-ch
	}
	// make sure pattern is same size as inputs
	if input.Len() != network.inputs.Len() {
		return nil, errors.New("invalid argument")
	}
	network.info("Run pattern to network: ")
	PrintList(input)
	i := uint(0)
	for v := input.Front(); v != nil; v = v.Next() {
		network.trace("\tsending %v to input %v\n", v, i)
		network.inputs[i].c <- v.Value.(Data) * network.inputs[i].weight
		i++
	}
	// get output
	output = new(list.List)
	for i := 0; i < network.outputs.Len(); i++ {
		o := <-network.outputs[i].c
		output.PushBack(o)
	}
	network.info(" \\- output is: ")
	PrintList(output)
	return output, err
}

func (network *Network) AddInput(number int, neuron *Neuron, weight Data) {
	network.inputs[number] = &Connection{c: neuron.input, weight: weight}
	ch := make(chan string)
	neuron.command <- message{message: "AddInput", payload: network.inputs[number], reply: ch}
	_ = <-ch
}

func (network *Network) Input(number int) *Connection {
	return network.inputs[number]
}

func (network *Network) Output(number int) *Connection {
	return network.outputs[number]
}

/**
 * log functions
 */
func (network *Network) LogLevel() int {
	return network.logLevel
}

func (network *Network) SetLogLevel(level int) {
	network.logLevel = level
}

func (network *Network) info(format string, a ...interface{}) {
	logMessage(network.logLevel, 1, format, a...)
}

func (network *Network) debug(format string, a ...interface{}) {
	logMessage(network.logLevel, 2, format, a...)
}

func (network *Network) trace(format string, a ...interface{}) {
	logMessage(network.logLevel, 3, format, a...)
}

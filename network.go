package network

import (
    "os"
    "fmt"
    "container/vector"
)

const NCPU = 4

type Data float64
type Connection struct {
    c chan Data
    weight Data
}
type ConnectionArray []*Connection


type message struct {
    message string
    payload interface{}
    reply chan string
}

type Network struct  {
    inputs  ConnectionArray // for presenting a pattern to the network
    outputs ConnectionArray // for collecting the results
    command map[string] chan message // for sending a command to all neurons
    log_level int
}

/* the real neuron inside a goroutine */
type neuron struct {
    name  string
    input chan Data
    outputs *vector.Vector
    command chan message
    inputs int
    inputs_waiting int
    accumulator Data
    threshold Data
    log_level int
}

/* object to communicate with neuron with */
type Neuron struct {
    command chan message
    input   chan Data
}

func newNetwork() *Network {
    return &Network{ command: make (map[string] chan message), log_level: 1}
}

func newConnection(w Data) *Connection {
    return &Connection {c: make( chan Data, NCPU), weight: w }
}

func newNeuron(name string,  outputs *vector.Vector, 
               command chan message, threshold Data) *neuron {
    return &neuron{name: name, input: make(chan Data, NCPU) , outputs: outputs, command: command , threshold: threshold}
}

func newConnectionArray(connections int) *ConnectionArray {
    ca := make(ConnectionArray, connections) 
/*    for i := 0; i < connections; i++ {
        ca[i] = newConnection(0)
    }*/
    return &ca
}

func PatternElm(e Data) Data {
    return e
}

func MakePattern(a...Data) (*vector.Vector){

    pattern := new(vector.Vector)
    for _, v := range a {
        pattern.Push(v)
    }
    return pattern
}

func CompareVector(v1 *vector.Vector, v2 *vector.Vector) bool {
    if (v1.Len() != v2.Len()) {
        return false
    }
    for i := 0; i < v1.Len(); i++ {
        if v1.At(i) != v2.At(i) {
            return false
        }
    }
    return true
}

func (v *ConnectionArray) Len() int{
    return len(*v)
}

func CreateNetwork(inputs int, outputs int) *Network {
    
    n := newNetwork()
    n.debug("Creating network\n")
    n.inputs = *newConnectionArray(inputs) 
    n.outputs = *newConnectionArray(outputs)

    for i := 0; i < outputs; i++ {
        n.outputs[i] = newConnection(1)
    }
    return n
}

func (network *Network) LogLevel() int {
    return network.log_level
}

func (network *Network) SetLogLevel(level int)  {
    network.log_level = level
}

func (network *Network) RunPattern(input *vector.Vector) (output *vector.Vector, err os.Error) {

     network.debug("resetting network\n");
     ch := make(chan string)
     for _, v := range network.command {
         v <- message{ message: "Reset", reply: ch } // send reset
         // wait for response
         _ = <-ch
     }
     // make sure pattern is same size as inputs
     if (input.Len() != network.inputs.Len() ) {
        return nil, os.EINVAL
     }
     network.info("presenting pattern %v to network\n", *input)
     for i := 0; i < input.Len(); i++  {
         v := input.At(i).(Data)
         network.trace("\tsending %v to input %v\n", v, i)
         network.inputs[i].c <- v * network.inputs[i].weight
     }
     output = new(vector.Vector)
     for i := 0; i < network.outputs.Len(); i++ {
        o:=  <-network.outputs[i].c
        output.Push(o)
     }
     network.info(" \\- output is %v\n", *output)
     return output, err
}


func (network *Network) AddInput(number int, neuron *Neuron, weight Data ) {
    network.inputs[number] = &Connection{ c: neuron.input, weight: weight }
    ch := make(chan string)
    neuron.command <- message{ message: "AddInput", payload: network.inputs[number], reply: ch}
    _ = <-ch
}

func (network *Network) Input(number int) *Connection {
    return network.inputs[number]
}

func (network *Network) Output(number int) *Connection {
    return network.outputs[number]
}


func (network *Network) CreateNeuron(name string, threshold Data) *Neuron {

    network.debug("\tCreating %v\n", name)
    network.command[name] = make( chan message, NCPU)
    
    n := newNeuron(name,  new(vector.Vector), network.command[name], threshold)

    defer network.addNeuron(n)
    return &Neuron{command: network.command[name], input: n.input }
}


func (network *Network) addNeuron(neuron *neuron) {
    go neuronLoop(neuron)
}

func neuronLoop(neuron *neuron) {
        for {
        select {
        case v := <-neuron.input:

                neuron.trace("recieved %v", v)
                neuron.inputs --
                neuron.trace("neuron.inputs remaing is %v", neuron.inputs)
                neuron.accumulator += v
                neuron.trace("neuron.accumulator is %v", neuron.accumulator)
                if ( neuron.inputs == 0 ) {
                    var fire_val Data
                    if ( neuron.accumulator >= neuron.threshold ) {
                        fire_val = 1.0
                    } else {
                        fire_val = 0.0
                    }
                    neuron.trace("is fireing value %v", fire_val)
                    for i := 0; i < neuron.outputs.Len(); i++ {
                        output := neuron.outputs.At(i).(*Connection)
                        output.c <- output.weight * fire_val
                    }
                }
        case command := <-neuron.command:
                      switch command.message {
                          case "AddOutput":
                            neuron.trace("adding output %v", command.payload)
                            neuron.outputs.Push(command.payload)
                            command.reply <- "Ok"
                          case "AddInput":
                            neuron.inputs_waiting++
                            neuron.inputs++
                            neuron.trace("input_wating is %v inputs is %v", neuron.inputs_waiting, neuron.inputs)
                            command.reply <- "Ok"
                          case "Reset":
                            neuron.inputs = neuron.inputs_waiting
                            neuron.accumulator = 0
                            command.reply <- "Ok"
                          case "SetLogLevel":
                            neuron.log_level = command.payload.(int)
                            command.reply <- "Ok"
                          default:
                            if command.reply != nil {
                                command.reply <- fmt.Sprintf("%v Error: Unknown command: %v", neuron.name, command.message)
                            } else {
                                fmt.Printf("%v Error: Unknown command: %v", neuron.name, command.message)
                            }
                      }

        }
        }
    } 

func (network *Network) info(format string, a ...interface{}) {
    log_message(network.log_level, 1, format, a...)
}

func (network *Network) debug(format string, a ...interface{}) {
    log_message(network.log_level, 2, format, a...)
}

func (network *Network) trace (format string, a ...interface{}) {
    log_message(network.log_level, 3, format, a...)
}

func log_message(log_level int, message_level int, format string, a...interface{}) {
    if log_level >= message_level {
        fmt.Printf(format,  a...)
    }
}

func (neuron *Neuron) AddOutput(conn *Connection ) {
    ch := make(chan string)
    neuron.command <- message{ message: "AddOutput",  payload: conn, reply: ch}
    // synchronise
    _ = <-ch
}

// Add a neuron as an input with weight 
func (neuron *Neuron) AddInput(input *Neuron, w Data) {
    input.AddOutput(&Connection {c: neuron.input, weight: w })
    ch := make(chan string)
    neuron.command <- message{ message: "AddInput", payload: input, reply: ch}
    _ = <-ch
}

func (neuron *neuron) LogLevel() int {
    return neuron.log_level
}

func (neuron *Neuron) SetLogLevel(level int)  {

    ch := make(chan string)
    neuron.command <- message{ message: "SetLogLevel", payload: level, reply: ch}
    _ = <-ch
}

func (neuron *neuron) info (format string, a ...interface{}) {
        real_format := fmt.Sprintf("\t\t%v: %v\n", neuron.name, format)
        log_message(neuron.log_level, 1, real_format,  a...)
}

func (neuron *neuron) debug(format string, a ...interface{}) {
        real_format := fmt.Sprintf("\t\t%v: %v\n", neuron.name, format)
        log_message(neuron.log_level, 2, real_format,  a...)
}

func (neuron *neuron) trace(format string, a ...interface{}) {
        real_format := fmt.Sprintf("\t\t%v: %v\n", neuron.name, format)
        log_message(neuron.log_level, 3, real_format,  a...)
}


/*

    input 1 -\-/---neuron 1 
              X
    input 2 -/-\---neuron 2 

*/

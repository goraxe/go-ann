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
}

type Network struct  {
    inputs  ConnectionArray // for presenting a pattern to the network
    outputs ConnectionArray // for collecting the results
    command map[string] chan message // for sending a command to all neurons
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
    return &Network{ command: make (map[string] chan message)}
}

func newConnection(w Data) *Connection {
    return &Connection {c: make( chan Data, NCPU), weight: w }
}

func newNeuron(name string,  outputs *vector.Vector, 
               command chan message, inputs int, threshold Data) *neuron {
    return &neuron{name: name, input: make(chan Data, NCPU) , outputs: outputs, command: command , inputs: inputs, inputs_waiting: inputs, threshold: threshold}
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

func (v *ConnectionArray) Len() int{
    return len(*v)
}

func CreateNetwork(inputs int, outputs int) *Network {
    n := newNetwork()
    n.inputs = *newConnectionArray(inputs) 
    n.outputs = *newConnectionArray(outputs)

    for i := 0; i < outputs; i++ {
        n.outputs[i] = newConnection(1)
    }
    return n
}

/*func (v *Vector) Len() int{
    return len(*v)
}*/


func (network *Network) RunPattern(input *vector.Vector) (output *vector.Vector, err os.Error) {

     fmt.Printf("resetting network\n");
     ch := make(chan message)
     for _, v := range network.command {
         v <- message{ "Reset", ch } // send reset
         // wait for response
         _ = <-ch
     }
     if (input.Len() != network.inputs.Len() ) {
        return nil, os.EINVAL
     }
     fmt.Printf("presenting pattern to network\n")
     for i := 0; i < input.Len(); i++  {
         v := input.At(i).(Data)
         fmt.Printf("\tsending %v to input %v\n", v, i)
         network.inputs[i].c <- v * network.inputs[i].weight
     }
     output = new(vector.Vector)
     for i := 0; i < network.outputs.Len(); i++ {
        o:=  <-network.outputs[i].c
        output.Push(o)
     }
     return output, err
}


func (network *Network) AddInput(number int, neuron *Neuron, weight Data ) {
    network.inputs[number] = &Connection{ c: neuron.input, weight: weight }
}

func (network *Network) Input(number int) *Connection {
    return network.inputs[number]
}

func (network *Network) Output(number int) *Connection {
    return network.outputs[number]
}


func (network *Network) CreateNeuron(name string, 
                                     inputs int, threshold Data) *Neuron {

    network.command[name] = make( chan message, NCPU)
    
    n := newNeuron(name,  new(vector.Vector), network.command[name], inputs, threshold)

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
                          case "Reset":
                            neuron.inputs = neuron.inputs_waiting
                            neuron.accumulator = 0
                            command.payload.(chan message) <- message{message: "Ok"} 
                      }

        }
        }
    } 


func (neuron *Neuron) AddOutput(conn *Connection ) {
    
    neuron.command <- message{ "AddOutput",  conn}
}

// a little backwards maybe 
func (neuron *Neuron) AddInput(input *Neuron, w Data) {
    input.AddOutput(&Connection {c: neuron.input, weight: w })
}


func (neuron *neuron) trace(format string, a ...interface{}) {
    if neuron.log_level > 5 {
        real_format := fmt.Sprintf("\t\t%v: %v\n", neuron.name, format)
        fmt.Printf(real_format,  a...)
    }
}


/*

    input 1 -\-/---neuron 1 
              X
    input 2 -/-\---neuron 2 

*/

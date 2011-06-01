package main

import (
    "os"
    "fmt"
    "./network"
    "container/vector"
)


func main() {
    // create a network with two inputs one output
    ann :=  network.CreateNetwork(2, 1)
    ann.SetLogLevel(0)

    neuron1 := ann.CreateNeuron("neuron1", 1.0)
//    neuron1.SetLogLevel(3)
    ann.AddInput(0, neuron1, 1)

    neuron2 := ann.CreateNeuron("neuron2", 1.0)
//    neuron2.SetLogLevel(3)
    ann.AddInput(1, neuron2, 1)
    
    neuron3 := ann.CreateNeuron("neuron3", 1.0)
//    neuron3.SetLogLevel(3)
    neuron3.AddInput(neuron1, 1)
    neuron3.AddInput(neuron2, -1)

    neuron4 := ann.CreateNeuron("neuron4", 1.0)
//    neuron4.SetLogLevel(3)
    neuron4.AddInput(neuron1, -1)
    neuron4.AddInput(neuron2, 1)

    neuron5 := ann.CreateNeuron("neuron5", 1.0)
//    neuron5.SetLogLevel(3)
    neuron5.AddInput(neuron3, 1)
    neuron5.AddInput(neuron4, 1)

    output1 := ann.Output(0)
    neuron5.AddOutput(output1)

    var (
        output *vector.Vector
        expected *vector.Vector
        e os.Error
    )
    pattern1 := network.MakePattern(0,0)
    expected = network.MakePattern(0)

    output, e = ann.RunPattern(pattern1)
    if(e == os.EINVAL) {
        os.Exit(1)
    }
    if(network.CompareVector(expected, output)) {
        fmt.Printf("ok: expected output pattern 1\n")
    } else {
        fmt.Printf("fail: expected output pattern 1 got: %v expected %v\n", output, expected)
    }

    pattern2 := network.MakePattern(0,1)
    expected = network.MakePattern(1)

    output, e = ann.RunPattern(pattern2)
    if(e == os.EINVAL) {
        os.Exit(1)
    }
    if(network.CompareVector(expected, output)) {
        fmt.Printf("ok: expected output pattern 2\n")
    } else {
        fmt.Printf("fail: expected output pattern 2 got: %v expected %v\n", output, expected)
    }

    pattern3 := network.MakePattern(1,0)
    expected = network.MakePattern(1)


    output, e = ann.RunPattern(pattern3)
    if(e == os.EINVAL) {
        os.Exit(1)
    }
    if(network.CompareVector(expected, output)) {
        fmt.Printf("ok: expected output pattern 3\n")
    } else {
        fmt.Printf("fail: expected output pattern 3 got: %v expected %v\n", output, expected)
    }

    pattern4 := network.MakePattern(1,1)
    expected = network.MakePattern(0)

    output, e = ann.RunPattern(pattern4)
    if(e == os.EINVAL) {
        os.Exit(1)
    }
    if(network.CompareVector(expected, output)) {
        fmt.Printf("ok: expected output pattern 4\n")
    } else {
        fmt.Printf("fail: expected output pattern 4 got: %v expected %v\n", output, expected)
    }
}


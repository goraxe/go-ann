package main

import (
	"myann"
	"os"
	"fmt"
	"errors"
	"container/list"
)

func main() {
	/* create a network with two inputs one output */
	network :=  myann.CreateNetwork(2, 1)
	network.SetLogLevel(3)

	neuron1 := network.CreateNeuron("neuron1", 1.0)
	//neuron1.SetLogLevel(3)
	network.AddInput(0, neuron1, 1)

	neuron2 := network.CreateNeuron("neuron2", 1.0)
	//neuron2.SetLogLevel(3)
	network.AddInput(1, neuron2, 1)
	
	neuron3 := network.CreateNeuron("neuron3", 1.0)
	//neuron3.SetLogLevel(3)
	neuron3.AddInput(neuron1, 1)
	neuron3.AddInput(neuron2, -1)

	neuron4 := network.CreateNeuron("neuron4", 1.0)
	//neuron4.SetLogLevel(3)
	neuron4.AddInput(neuron1, -1)
	neuron4.AddInput(neuron2, 1)

	neuron5 := network.CreateNeuron("neuron5", 1.0)
	//neuron5.SetLogLevel(3)
	neuron5.AddInput(neuron3, 1)
	neuron5.AddInput(neuron4, 1)

	output1 := network.Output(0)
	neuron5.AddOutput(output1)

	/* test network */
	var (
		output *list.List
		expected *list.List
		e error
	)
	
	// pattern 0 test
	pattern1 := myann.MakePattern(0,0)
	expected = myann.MakePattern(0)

	output, e = network.RunPattern(pattern1)
	if(e == errors.New("invalid argument")) {
		os.Exit(1)
	}
	if(myann.CompareList(expected, output)) {
		fmt.Printf("ok: expected output pattern 0\n")
	} else {
		fmt.Printf("fail: expected output pattern 0\n")
		myann.PrintList(expected)
		myann.PrintList(output)
	}
	// pattern 1 test
	pattern2 := myann.MakePattern(0,1)
	expected = myann.MakePattern(1)

	output, e = network.RunPattern(pattern2)
	if(e == errors.New("invalid argument")) {
		os.Exit(1)
	}
	if(myann.CompareList(expected, output)) {
		fmt.Printf("ok: expected output pattern 1\n")
	} else {
		fmt.Printf("fail: expected output pattern 1 got: %v expected %v\n", output, expected)
	}
	// pattern 2 test
	pattern3 := myann.MakePattern(1,0)
	expected = myann.MakePattern(1)

	output, e = network.RunPattern(pattern3)
	if(e == errors.New("invalid argument")) {
		os.Exit(1)
	}
	if(myann.CompareList(expected, output)) {
		fmt.Printf("ok: expected output pattern 2\n")
	} else {
		fmt.Printf("fail: expected output pattern 2 got: %v expected %v\n", output, expected)
	}
	// pattern 3 test
	pattern4 := myann.MakePattern(1,1)
	expected = myann.MakePattern(0)

	output, e = network.RunPattern(pattern4)
	if(e == errors.New("invalid argument")) {
		os.Exit(1)
	}
	if(myann.CompareList(expected, output)) {
		fmt.Printf("ok: expected output pattern 3\n")
	} else {
		fmt.Printf("fail: expected output pattern 3 got: %v expected %v\n", output, expected)
	}
}


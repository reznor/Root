# Invocation

After [installing Go](https://golang.org/doc/install), the program can be run in either of these 2 ways:

>$ go run main.go input.txt
>
>$ cat input.txt | go run main.go

# Overview

The central recurring theme (and guiding principle) is a focus on a production-ready architecture for future extensibility -- putting
structures into place that make it easy for (and encourages) future developers to Do The Right Thing (tm) without needing to make any
sweeping changes as long as the overall shape of the problem doesn't fundamentally change.

That can be seen in the strict interface boundaries laid out between packages (and the isolation of package-internal data domains),
the type definitions that allow for minimally-disruptive changes to the modeling of central concepts, and the data structures and
algorithms chosen that shouldn't require any modifications until we hit large scale.

There are detailed "Maintainer Notes" in many files that talk about these things in more specific detail.

# Architecture

Aside from modeling the system on components that follow very naturally from breaking down the problem statement into core responsibilities,
the architecture has leaned towards being amenable to running in a streaming+serverless environment as far as possible -- the goal is that it
should be relatively simple to slap this software onto an AWS environment composed of (say) Kinesis (for streaming infrastructure),
Lambda (for serverless compute), and DynamoDB (for storage).

## Packages

The flow of the system (and thus the interaction of the packages) is described below in terms of the elevator pitches of their main components.

### [input](input/)

Reads from an event source (in the case of our specific scenario, an input file) and generates a stream of `input.EventEnvelope` objects.

### [eventprocessor](eventprocessor/)

Receives a stream of `input.EventEnvelope` objects and a handle to an `eventstore.EventStore`, and emits a stream of `error`s that may
result from processing those events.

Parses the `input.EventEnvelope` objects just enough to be able to deduce the `eventhandler.EventType`, based off of which it delegates to the
concrete `eventhandler.Interface` implementation registered with `eventhandler.GlobalRegistry()`.

### [eventhandler](eventhandler/)

Provides independent sub-modules that implement `eventhandler.Interface` for handling all the events supported by the system, and that store
event information of value to downstream systems in the `eventstore.EventStore` that is passed in to them.

The [`eventhandler/` README](eventhandler/README.md) talks more about the reasoning for the directory structure in there, as well as the steps to add support for new events.

### [eventstore](eventstore/)

Serves as storage for retain-worthy event information.

Provides `eventstore.VisitorInterface` for inspection of all that retained information.

### [output](output/)

Provides `output.ReportGenerator` that implements `eventstore.VisitorInterface` and generates a report in the desired output format.

### [mathutils](mathutils/)

A collection of shared utilities for mathematical operations that are hard to get right.

## Concurrency

In keeping with the streaming-centric architecture of the system, there's 3 concurrent threads of execution that are naturally modeled as
stages of a streaming pipeline:

1. Reading the input from the event source to generate events (`input.StartReading()`).
2. Processing the events (`eventprocessor.EventProcessor.Process()`).
3. Handling the errors from processing the events (this happens in [main.go](main.go) on the main thread).

# Testing

There's near-100% unit test coverage for every package in the system, and the tests utilize multiple techniques (as appropriate for the
component under test):

1. Table-Driven Tests -- this  is the recommended Golang style and is used for all but [one component's](eventhandler/registry_test.go)
   tests (because the structure of each test case was different enough that trying to generalize it would make the hypothetical single unified
   test case too unreadable).

2. Fakes + Behavioral Testing -- see [eventprocessor_test.go](eventprocessor/eventprocessor_test.go).

3. Declarative Testing -- see [eventstore_test.go](eventstore/eventstore_test.go).
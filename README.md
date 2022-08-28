# Dapr components go SDK POC

This repository is a POC of a SDK for [Dapr gRPC Components (a.k.a pluggable components)](https://github.com/dapr/dapr/issues/4925)

## Running examples

Start by running `./run.sh` inside `/examples` folder. It will start the daprd runtime with pluggable components version + in memory state store implementation from components-contrib.

## Implementing a Pluggable State Store component.

To create your own implementation:

1. Create a new folder under `/examples`
2. Implement a stateStore using the sdk
3. Run `./run.sh your_folder_goes_here`

This will build your component and bootstrap the dapr runtime with the default options.

## Getting started

Creating a new component is nothing more than implement a StateStore interface and Run the dapr component server.

```golang
package main

import (
	dapr "github.com/mcandeia/dapr-components-go-sdk"
)

func main() {
	dapr.MustRun(dapr.UseStateStore(MyComponentGoesHere{}))
}
```

Bus
======

Package eventbus is an async bus with batteries included for Golang.

#### Installation
Make sure that Go is installed on your computer.
Type the following command in your terminal:

    go get github.com/tibrn/eventbus


After it the package is ready to use.

#### Import package in your project
Add following line in your `*.go` file:
```go
import "github.com/tibrn/eventbus"
```
#### Example
```go
func add(a int, b int) {
	fmt.Printf("%d\n", a + b)
}

func main() {
	bus := eventbus.New();
	bus.On("add", add);
	bus.Emit("add", 20, 40);
}
```

#### Implemented methods
* **New()**
* **On()**
* **Emit()**

#### New()
New returns new EventBus with empty handlers.
```go
bus := eventbus.New();
```

#### On(event EventName, fn interface{}) (RemoveCallback, error)
Subscribe to event. Returns error if `fn` is not a function.
```go
func Handler() { ... }
...
removeHandler, err := bus.On("event:handler", Handler)
```

#### Emit(event EventName, payload ...interface{}) WaitCallback
Emit executes callback defined for an event. Any additional argument will be transferred to the callback.
```go
arg1, arg2 := 4,5
...
waitEmit := bus.Emit("event:handler", arg1, arg2)
```

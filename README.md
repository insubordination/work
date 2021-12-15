![go workers](https://raw.githubusercontent.com/insubordination/work/assets/constworker_header_anim.gif)

# Examples
* [Quickstart](https://github.com/insubordination/work/blob/master/examples/quickstart/quickstart.go)
* [Multiple Go Workers](https://github.com/insubordination/work/blob/master/examples/multiple_workers/multiplework.go)
* [Passing Fields](https://github.com/insubordination/work/blob/master/examples/passing_fields/passingfields.go)
# Getting Started
### Pull in the dependency
```zsh
go get github.com/insubordination/work
```

### Add the import to your project
giving an alias helps since work doesn't exactly follow conventions.    
_(If you're using a JetBrains IDE it should automatically give it an alias)_
```go
import (
    "github.com/insubordination/work"
)
```
### Create a new worker <img src="https://raw.githubusercontent.com/insubordination/work/assets/constworker.png" alt="worker" width="35"/>
The NewWorker factory method returns a new worker.    
_(Method chaining can be performed on this method like calling .Work() immediately after.)_
```go
type MyWorker struct {}

func NewMyWorker() Worker {
	return &MyWorker{}
}

func (my *MyWorker) Work(in interface{}, out chan<- interface{}) error {
	// work iteration here
}

runner := work.NewRunner(ctx, NewMyWorker(), numberOfWorkers)
```
### Send work to worker
Send accepts an interface.  So send it anything you want.
```go
runner.Send("Hello World")
```
### Wait for the worker to finish and handle errors
Any error that bubbles up from your worker functions will return here.
```go
if err := runner.Wait(); err != nil {
    //Handle error
}
```

## Working With Multiple Workers
### Passing work form one worker to the next 

By using the InFrom method you can tell `workerTwo` to accept output from `workerOne`
```go
runnerOne := work.NewRunner(ctx, NewMyWorker(), 100).Work()
runnerTwo := work.NewRunner(ctx, NewMyWorkerTwo(), 100).InFrom(workerOne).Work()
```
### Accepting output from multiple workers
It is possible to accept output from more than one worker but it is up to you to determine what is coming from which worker.  (They will send on the same channel.)
```go
runnerOne := work.NewRunner(ctx, NewMyWorker(), 100).Work()
runnerTwo := work.NewRunner(ctx, NewMyWorkerTwo(), 100).Work()
runnerThree := work.NewRunner(ctx, NewMyWorkerThree(), 100).InFrom(workerOne, workerTwo).Work()
```

## Passing Fields To Workers
### Adding Values
Fields can be passed via the workers object. Be sure as with any concurrency in Golang that your variables are concurrent safe.  Most often the golang documentation will state the package or parts of it are concurrent safe.  If it does not state so there is a good chance it isn't.  Use the sync package to lock and unlock for writes on unsafe variables.  (It is good practice NOT to defer in the work function.)

<img src="https://raw.githubusercontent.com/insubordination/work/assets/constworker2.png" alt="worker" width="35"/> **ONLY** use the `Send()` method to get data into your worker. It is not shared memory unlike the worker objects values.

```go
type MyWorker struct {
	message string
}

func NewMyWorker(message string) Worker {
	return &MyWorker{message}
}

func (my *MyWorker) Work(in interface{}, out chan<- interface{}) error {
	fmt.Println(my.message)
}

runner := work.NewRunner(ctx, NewMyWorker(), 100).Work()
```

### Setting Timeouts or Deadlines
If your workers needs to stop at a deadline or you just need to have a timeout use the SetTimeout or SetDeadline methods. (These must be in place before setting the workers off to work.)
```go
 // Setting a timeout of 2 seconds
 timeoutRunner.SetTimeout(2 * time.Second)

 // Setting a deadline of 4 hours from now
 deadlineRunner.SetDeadline(time.Now().Add(4 * time.Hour))

func workerFunction(in interface{}, out chan<- interface{} error {
	fmt.Println(in)
	time.Sleep(1 * time.Second)
}
```


## Performance Hints
### Buffered Writer
If you want to write out to a file or just stdout you can use SetWriterOut(writer io.Writer).  The worker will have the following methods available
```go
runner.Println()
runner.Printf()
runner.Print()
```
The workers use a buffered writer for output and can be up to 3 times faster than the fmt package.  Just be mindful it won't write out to the console as quickly as an unbuffered writer.  It will sync and eventually flush everything at the end, making it ideal for writing out to a file.

### Using GOGC env variable
If your application is based solely around using workers, consider upping the percentage of when the scheduler will garbage collect. (ex. GOGC=200) 200% -> 300% is a good starting point. Make sure your machine has some good memory behind it.
By upping the percentage your application will interupt the workers less, meaning they get more work done.  However, be aware of the rest of your applications needs when modifying this variable.

### Using GOMAXPROCS env variable
For workers that run quick bursts of lots of simple data consider lowering the GOMAXPROCS.  Be carfeful though, this can affect your entire applicaitons performance.  Profile your application and benchmark it.  See where your application runs best.

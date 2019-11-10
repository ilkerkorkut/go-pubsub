## Go PubSub

This is an simple implementation of concurrent subscriber and publisher mechanism in go. So this tool makes repeatitive tasks(with different parameters) concurrent if needed.

## How to use
```go
// Create new instance of pubsub
// Give subscriber count, how many data will be published, your task func, if debug needed give true
func main() {
	// Data for your job
	data := [100]int{}
	for i := 0; i < len(data); i++ {
		data[i] = i
	}

	ps := pubsub.New(runtime.NumCPU(), len(data), myTask, false)
	// Start subscribers to listen channel
	ps.StartSubscribers()
	// Publish your data
	for i := 0; i < len(data); i++ {
		ps.Publish(data[i])
	}
	// If you need to wait your all tasks for done
	ps.Wait()
}

func myTask(data interface{}, done func()) {
	defer done()
	log.Println("My Task is running with data : ", data)
}
```
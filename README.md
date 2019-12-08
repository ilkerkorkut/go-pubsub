## Go PubSub

This is an simple implementation of concurrent subscriber and publisher mechanism in go. So this tool makes repetitive tasks(with different parameters) concurrent and distributed if needed. 

## How to use

#### Single Instance
```go
func main() {  
  
   data := [100]pubsub.DataPacket{}  
  
   for i := 0; i < len(data); i++ {  
      data[i] = pubsub.DataPacket{  
         Data: i,  
         Time: time.Now(),  
      }  
   }  
  
   ps, err := pubsub.New(runtime.NumCPU(), len(data), myTask, &pubsub.Config{  
      MultiNode:  false,  
      Name:    "my-pubsub-application",  
      Debug:      true,  
   })  
   if err != nil {  
      log.Fatal(err)  
   }
   for i := 0; i < len(data); i++ {  
      ps.Publish(data[i])  
   }  

   ps.Wait()  
   log.Println("Successfully completed!")
}  
  
func myTask(data interface{}, done func()) {  
   defer done()  
   log.Println("My Task is running with data : ", data)  
} 
```

#### Multi Node Instance

Go to the `./examples/pubsub_multinode` folder.

Run this commands for three nodes (or separately):

```
go run main.go 8080 8081 8082 & go run main.go 8081 8080 8082 & go run main.go 8082 8080 8081
```

```go
func main() {

	port := os.Args[1]
	otherPorts := os.Args[2:]

	data := [100]pubsub.DataPacket{}

	for i := 0; i < len(data); i++ {
		data[i] = pubsub.DataPacket{
			Data: i,
			Time: time.Now(),
		}
	}
    ps, err := pubsub.New(runtime.NumCPU(), len(data), myTask, &pubsub.Config{
        MultiNode:  true,
        ServerPort: port,
        NodePorts:  otherPorts,
        Name:       "my-pubsub-application",
        Debug:      true,
    })
    if err != nil {
        log.Fatal(err)
    }

    for i := 0; i < len(data); i++ {
        ps.Publish(data[i])
    }
    log.Println("waiting!")
    ps.Wait()
    log.Println("Successfully completed!")

}

func myTask(data interface{}, done func()) {
	defer done()
	log.Println("My Task is running with data : ", data)
}
```

### Goal
With implementation, when an instance was triggered to publish data to processing them, the other instances will consume that published data equally. So all resources will work at the same time.
Discovery service will be added to easily manage instance, kubernetes readiness for pod to pod communication. That's just an experimental.


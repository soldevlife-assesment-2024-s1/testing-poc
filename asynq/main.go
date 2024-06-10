package main

import (
	"log"
	tasks "test-asynq/task"
	"time"

	"net/http"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

const redisAddr = "localhost:6379"

// func main() {
// 	srv := asynq.NewServer(
// 		asynq.RedisClientOpt{Addr: redisAddr},
// 		asynq.Config{
// 			// Specify how many concurrent workers to use
// 			Concurrency: 10,
// 			// Optionally specify multiple queues with different priority.
// 			Queues: map[string]int{
// 				"critical": 6,
// 				"default":  3,
// 				"low":      1,
// 			},
// 			// See the godoc for other configuration options
// 		},
// 	)

// 	// mux maps a type to a handler
// 	mux := asynq.NewServeMux()
// 	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
// 	mux.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())
// 	// ...register other handlers...

// 	if err := srv.Run(mux); err != nil {
// 		log.Fatalf("could not run server: %v", err)
// 	}
// }

func main() {

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------

	task, err := tasks.NewEmailDeliveryTask(42, "some:template:id")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Example 2: Schedule task to be processed in the future.
	//            Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------

	info, err = client.Enqueue(task, asynq.ProcessIn(1*time.Minute))
	if err != nil {
		log.Fatalf("could not schedule task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// Example 3: Set other options to tune task processing behavior.
	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------

	task, err = tasks.NewImageResizeTask("https://example.com/myassets/image.jpg")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	go func() {
		h := asynqmon.New(asynqmon.Options{
			RootPath:     "/monitoring", // RootPath specifies the root for asynqmon app
			RedisConnOpt: asynq.RedisClientOpt{Addr: redisAddr, DB: 0},
		})

		// Note: We need the tailing slash when using net/http.ServeMux.
		http.Handle(h.RootPath()+"/", h)

		// Go to http://localhost:8080/monitoring to see asynqmon homepage.
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	mux.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}

}

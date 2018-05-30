package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/bench"
)

// Some sane defaults
const (
	DefaultNumMsgs     = 100000
	DefaultNumPubs     = 1
	DefaultNumSubs     = 10
	DefaultMessageSize = 128
)

func usage() {
	log.Fatalf("Usage: nats-bench [-s server (%s)] [--tls] [-np NUM_PUBLISHERS] [-ns NUM_SUBSCRIBERS] [-n NUM_MSGS] [-ms MESSAGE_SIZE] [-csv csvfile] <subject>\n", nats.DefaultURL)
}

var benchmark *bench.Benchmark

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var tls = flag.Bool("tls", false, "Use TLS Secure Connection")
	var numReqs = flag.Int("nreq", DefaultNumPubs, "Number of Concurrent Requester")
	var numResps = flag.Int("nresp", DefaultNumSubs, "Number of Concurrent Responser")

	var numMsgs = flag.Int("n", DefaultNumMsgs, "Number of Messages to Publish")
	var msgSize = flag.Int("ms", DefaultMessageSize, "Size of the message.")
	var csvFile = flag.String("csv", "", "Save bench data to csv file")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		usage()
	}

	if *numMsgs <= 0 {
		log.Fatal("Number of messages should be greater than zero.")
	}

	// Setup the option block
	opts := nats.GetDefaultOptions()
	opts.Servers = strings.Split(*urls, ",")
	for i, s := range opts.Servers {
		opts.Servers[i] = strings.Trim(s, " ")
	}
	opts.Secure = *tls

	benchmark = bench.NewBenchmark("NATS", *numReqs, *numResps)

	var startwg sync.WaitGroup
	var donewg sync.WaitGroup

	donewg.Add(*numReqs + *numResps)

	// Run Subscribers first
	startwg.Add(*numReqs)

	for i := 0; i < *numReqs; i++ {
		go runRequest(&startwg, &donewg, opts, *numMsgs, *msgSize)
	}
	startwg.Wait()

	log.Println("start resp")
	// Now Publishers
	startwg.Add(*numResps)
	pubCounts := bench.MsgsPerClient(*numMsgs, *numResps)

	log.Println("counts",pubCounts)
	for i := 0; i < *numResps; i++ {
		go runResponse(&startwg, &donewg, opts, pubCounts[i], *msgSize,i)
	}
	log.Printf("Starting benchmark [msgs=%d, msgsize=%d, reqs=%d, reps=%d]\n", *numMsgs, *msgSize, *numReqs, *numResps)

	startwg.Wait()
	donewg.Wait()

	benchmark.Close()

	fmt.Print(benchmark.Report())

	if len(*csvFile) > 0 {
		csv := benchmark.CSV()
		ioutil.WriteFile(*csvFile, []byte(csv), 0644)
		fmt.Printf("Saved metric data in csv file %s\n", *csvFile)
	}
}


func runRequest(startwg, donewg *sync.WaitGroup, opts nats.Options, numMsgs int, msgSize int) {
	nc, err := opts.Connect()
	if err != nil {
		log.Fatalf("Request Can't connect: %v\n", err)
	}
	defer nc.Close()
	startwg.Done()

	args := flag.Args()
	subj := args[0]
	var msg []byte
	if msgSize > 0 {
		msg = make([]byte, msgSize)
	}

	start := time.Now()

	for i := 0; i < numMsgs; i++ {
		//log.Println("send",subj)
		nc.Request(subj, msg, time.Microsecond*200)
		//log.Println("sendReceived",subj,i)

	}
	nc.Flush()
	benchmark.AddPubSample(bench.NewSample(numMsgs, msgSize, start, time.Now(), nc))

	donewg.Done()
	log.Println("request finish")
}

func runResponse(startwg, donewg *sync.WaitGroup, opts nats.Options, numMsgs int, msgSize int,num int) {
	log.Println("resp started",flag.Args()[0])
	nc, err := opts.Connect()
	if err != nil {
		log.Fatalf("Resp Can't connect: %v\n", err)
	}

	args := flag.Args()
	subj := args[0]

	received := 0
	start := time.Now()
	nc.Subscribe(subj, func(msg *nats.Msg) {
		log.Println("received",received)
		nc.Publish(msg.Reply, []byte("hello"))
		log.Println("responed",received,numMsgs,num)
		received++
		if received >= numMsgs {
			log.Println("receive add",num)

			benchmark.AddSubSample(bench.NewSample(numMsgs, msgSize, start, time.Now(), nc))
			log.Println("receive done",num)

			donewg.Done()
			nc.Close()
			startwg.Done()
			log.Println("receive finish",num)

		}
	})

}

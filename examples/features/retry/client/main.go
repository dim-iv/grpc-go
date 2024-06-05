/*
 *
 * Copyright 2019 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Binary client is an example client.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/features/proto/echo"
)

var (
	addr = flag.String("addr", "localhost:50052", "the address to connect to")
	// see https://github.com/grpc/grpc/blob/master/doc/service_config.md to know more about service config
	retryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "grpc.examples.echo.Echo"}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": 4,
			  "InitialBackoff": ".01s",
			  "MaxBackoff": ".01s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`
)

// use grpc.WithDefaultServiceConfig() to set service config
func retryDial() (*grpc.ClientConn, error) {
	return grpc.NewClient(
            *addr,
            grpc.WithTransportCredentials(insecure.NewCredentials()),
            grpc.WithDefaultServiceConfig(retryPolicy),
            grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
        )
}

func testRequest(c pb.EchoClient) {
	//contextTimeout := 17 * time.Minute // [dim]: above 15 minutes
	contextTimeout := 1 * time.Minute // [dim]: exactly 1 minute

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	reply, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: "Try and Success"})
	if err != nil {
		//log.Fatalf("UnaryEcho error: %v", err)
		log.Printf("UnaryEcho error: %v", err)
	}
	log.Printf("UnaryEcho reply: %v", reply)
}

func waitCountdown(waitTime int) {
        for i := 1; i <= waitTime; i++ {
            log.Printf("... waiting %d / %d ...", i, waitTime)
            time.Sleep(1 * time.Second)
        }
}

func showSomeSpace() {
        log.Printf("")
        log.Printf("")
        log.Printf("")
}

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := retryDial()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Printf("failed to close connection: %s", e)
		}
	}()

	c := pb.NewEchoClient(conn)
        requestCount := 1

        for {
            showSomeSpace()
            log.Printf("Next request (%d)...", requestCount)
            log.Printf("... Waiting (time to drop)...")
            waitCountdown(7)
            testRequest(c)
            requestCount++;
        }

	select {} // wait forever
}

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
	"google.golang.org/grpc/keepalive"
)

var addr = flag.String("addr", "localhost:50052", "the address to connect to")

var kacp = keepalive.ClientParameters{
	Time:                50 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             10 * time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

func testRequest(c pb.EchoClient) {
	contextTimeout := 17 * time.Minute // [dim]: above 15 minutes
	//contextTimeout := 1 * time.Minute // [dim]: exactly 1 minute

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	reply, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: "keepalive demo"})
	if err != nil {
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

	conn, err := grpc.NewClient(
            *addr,
            grpc.WithTransportCredentials(insecure.NewCredentials()),
            grpc.WithKeepaliveParams(kacp),
            grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
        )
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

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

	select {} // Block forever; run with GODEBUG=http2debug=2 to observe ping frames and GOAWAYs due to idleness.
}

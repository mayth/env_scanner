// Copyright 2021 Mei Akizuru (mayth)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-ble/ble"
)

var (
	namePrefix = flag.String("prefix", "", "name prefix of reporters")

	uuidEnvSensing = ble.UUID16(0x181A)

	lastResults map[string]Measurement = make(map[string]Measurement)
)

func main() {
	flag.Parse()

	dev, err := NewDevice()
	if err != nil {
		log.Fatalf("failed to create a new device. %v", err)
	}
	ble.SetDefaultDevice(dev)

	var advFilter func(a ble.Advertisement) bool
	if *namePrefix != "" {
		advFilter = func(a ble.Advertisement) bool {
			return strings.HasPrefix(a.LocalName(), *namePrefix)
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	usrCh := make(chan os.Signal, 1)
	signal.Notify(usrCh, syscall.SIGUSR1)
	go func() {
	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-usrCh:
				log.Println("SIGUSR1: showing current environment reported")
				for reporter, m := range lastResults {
					log.Printf("%s\t%s\t%.2f\u2103\t%.2f%%", reporter, m.Timestamp.Local().Format(time.RFC3339), m.Temperature, m.Humidity)
				}
			}
		}
	}()

	log.Println("start scan")
	if err := ble.Scan(ctx, true, advHandler, advFilter); err != nil {
		if !errors.Is(context.Canceled, err) {
			log.Fatalf("scan failed. %v", err)
		}
	}
	log.Println("scan finished")
}

type Measurement struct {
	Timestamp   time.Time
	Temperature float32
	Humidity    float32
}

func Decode(b []byte) (Measurement, error) {
	r := bytes.NewReader(b)
	var ts uint64
	if err := binary.Read(r, binary.LittleEndian, &ts); err != nil {
		return Measurement{}, fmt.Errorf("failed to parse timestamp. %w", err)
	}
	var rawTemp uint16
	if err := binary.Read(r, binary.LittleEndian, &rawTemp); err != nil {
		return Measurement{}, fmt.Errorf("failed to parse temperature. %w", err)
	}
	var rawHumi uint16
	if err := binary.Read(r, binary.LittleEndian, &rawHumi); err != nil {
		return Measurement{}, fmt.Errorf("failed to parse humidity. %w", err)
	}
	temp := -45.0 + 175.0*(float32(rawTemp)/65535.0)
	humi := 100.0 * (float32(rawHumi) / 65535.0)
	return Measurement{time.Unix(int64(ts), 0), temp, humi}, nil
}

func advHandler(adv ble.Advertisement) {
	for _, srv := range adv.ServiceData() {
		if srv.UUID.Equal(uuidEnvSensing) {
			m, err := Decode(srv.Data)
			if err != nil {
				log.Printf("failed to decode. %v", err)
				continue
			}
			if r, ok := lastResults[adv.LocalName()]; ok && (m.Timestamp.Equal(r.Timestamp) || m.Timestamp.Before(r.Timestamp)) {
				continue
			}
			log.Printf("Update report from %s (%s). ts=%s, temp=%.2f, humi=%.2f",
				adv.LocalName(), adv.Addr().String(), m.Timestamp.Format(time.RFC3339), m.Temperature, m.Humidity)
			lastResults[adv.LocalName()] = m
		}
	}
}

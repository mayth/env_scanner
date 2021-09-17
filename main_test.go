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
	"math"
	"testing"
	"time"
)

func TestDecode(t *testing.T) {
	tolerance := 0.005
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Measurement
		wantErr bool
	}{
		{
			name: "success",
			// 2021-09-17 09:52:13: Temp=25.34 deg, Humi=59.93 %
			args: args{[]byte{0x4D, 0x65, 0x44, 0x61, 0x00, 0x00, 0x00, 0x00, 0xe7, 0x66, 0x68, 0x99}},
			want: Measurement{
				Timestamp:   time.Unix(1631872333, 0),
				Temperature: 25.34,
				Humidity:    59.93,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Timestamp.Equal(tt.want.Timestamp) {
				t.Errorf("Decode() Timestamp = %v, want %v", got.Timestamp, tt.want.Timestamp)
			}
			if math.Abs(float64(got.Temperature)-float64(tt.want.Temperature)) > tolerance {
				t.Errorf("Decode() Temperature = %v, want %v (tolerance %v)", got.Temperature, tt.want.Temperature, tolerance)
			}
			if math.Abs(float64(got.Humidity)-float64(tt.want.Humidity)) > tolerance {
				t.Errorf("Decode() Humidity = %v, want %v (tolerance %v)", got.Humidity, tt.want.Humidity, tolerance)
			}
		})
	}
}

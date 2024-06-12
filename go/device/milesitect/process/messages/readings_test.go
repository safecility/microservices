package messages

import (
	"reflect"
	"testing"
)

func TestReadMilesiteCT(t *testing.T) {
	type args struct {
		payload []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *MilesiteCTReading
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "read all",
			args: args{
				payload: []byte{0xFF, 0x0B, 0xFF,
					0xFF, 0x01, 0x01,
					0xFF, 0x16, 0x67, 0x46, 0xD3, 0x88, 0x02, 0x58, 0x00, 0x00,
					0xFF, 0x09, 0x01, 0x00,
					0xFF, 0x0A, 0x01, 0x01,
					0x03, 0x97, 0x10, 0x27, 0x00, 0x00,
					0x84, 0x98, 0xB8, 0x0B, 0xD0, 0x07, 0xC4, 0x09, 0x05},
			},
			want: &MilesiteCTReading{
				Device: nil,
				UID:    "6746d38802580000",
				Power:  true,
				Version: Version{
					Ipso:     "0.1",
					Hardware: "1.0",
					Firmware: "1.1",
				},
				Current: Current{
					Total: 100,
					Value: 25,
					Max:   30,
					Min:   20,
					Alarms: Alarms{
						t:  true,
						tr: false,
						r:  true,
						rr: false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadMilesiteCT(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMilesiteCT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadMilesiteCT() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readAlarm(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want Alarms
	}{
		// TODO: Add test cases.
		{
			name: "read payload",
			args: args{
				b: 0x01,
			},
			want: Alarms{
				t: true,
			},
		},
		{
			name: "read payload",
			args: args{
				b: 0x05,
			},
			want: Alarms{
				t: true,
				r: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readAlarm(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readAlarm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readSlice(t *testing.T) {
	type args struct {
		r       *MilesiteCTReading
		payload []byte
		offset  int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "power",
			args: args{
				r:       &MilesiteCTReading{},
				payload: []byte{0xFF, 0x0B, 0xFF},
				offset:  0,
			},
			want: 3,
		},
		{
			name: "ipso",
			args: args{
				r:       &MilesiteCTReading{},
				payload: []byte{0xFF, 0x01, 0x01},
				offset:  0,
			},
			want: 3,
		},
		{
			name: "serial",
			args: args{
				r:       &MilesiteCTReading{},
				payload: []byte{0xFF, 0x16, 0x67, 0x46, 0xD3, 0x88, 0x02, 0x58, 0x00, 0x00},
				offset:  0,
			},
			want: 10,
		},
		{
			name: "hardware",
			args: args{
				r:       &MilesiteCTReading{},
				payload: []byte{0xFF, 0x09, 0x01, 0x00},
				offset:  0,
			},
			want: 4,
		},
		{
			name: "firmware",
			args: args{
				r:       &MilesiteCTReading{},
				payload: []byte{0xFF, 0x0A, 0x01, 0x01},
				offset:  0,
			},
			want: 4,
		},
		{
			name: "total current",
			args: args{
				r:       &MilesiteCTReading{},
				payload: []byte{0x03, 0x97, 0x10, 0x27, 0x00, 0x00},
				offset:  0,
			},
			want: 6,
		},
		{
			name: "current alarm",
			args: args{
				r:       &MilesiteCTReading{},
				payload: []byte{0x84, 0x98, 0xB8, 0x0B, 0xD0, 0x07, 0xC4, 0x09, 0x05},
				offset:  0,
			},
			want: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readSlice(tt.args.r, tt.args.payload, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("readSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readVersion(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readVersion(tt.args.b); got != tt.want {
				t.Errorf("readVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

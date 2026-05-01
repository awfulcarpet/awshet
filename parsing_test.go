package main

import "testing"

func Test_parseTime(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    int
		want2   int
		wantErr bool
	}{
		{"basic", "00:00", 0, 0, false},
		{"hh:mm", "13:49", 13, 49, false},
		{"h:mm", "5:37", 5, 37, false},
		{"0h:mm", "08:24", 8, 24, false},
		{"0h:m", "07:9", 7, 9, false},
		{"h:m", "1:3", 1, 3, false},
		{"h:0m", "9:02", 9, 2, false},
		{"hm", "13", -1, -1, true},
		{"25:mm", "25:30", -1, -1, true},
		{"12:60", "12:60", -1, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, gotErr := parseTime(tt.str)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("parseTime() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("parseTime() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("parseTime() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("parseTime() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

package scanner_test

import (
	"testing"
	"time"

	"go.senan.xyz/socr/backend/scanner"
)

func TestGuessFileCreated(t *testing.T) {
	tcases := []struct {
		filename string
		stamp    time.Time
	}{
		{filename: "20170520_012747.png", stamp: time.Date(2017, 05, 20, 01, 27, 47, 0, time.UTC)},
		{filename: "20170520_012747.mkv", stamp: time.Date(2017, 05, 20, 01, 27, 47, 0, time.UTC)},
		{filename: "IMG_20190405_134142.png", stamp: time.Date(2019, 4, 5, 13, 41, 42, 0, time.UTC)},
		{filename: "img_19940405_134142 (copy 1).png", stamp: time.Date(1994, 4, 5, 13, 41, 42, 0, time.UTC)},
		{filename: "2021-05-18T14:14:30+01:00.png", stamp: time.Date(2021, 05, 18, 13, 14, 30, 0, time.UTC)},
		{filename: "2011-05-18T14:14:30+01:00.png", stamp: time.Date(2011, 05, 18, 13, 14, 30, 0, time.UTC)},
	}

	fallback := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	for _, tcase := range tcases {
		result := scanner.GuessFileCreated(tcase.filename, fallback)
		if same := result.Equal(tcase.stamp); !same {
			t.Errorf("filename %q parsed %q expected %q", tcase.filename, result, tcase.stamp)
		}
	}
}

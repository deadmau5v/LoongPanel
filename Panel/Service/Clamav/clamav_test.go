package clamav_test

import (
	clamav "LoongPanel/Panel/Service/Clamav"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	output := `Known viruses: 0
Engine version: 0.103.8
Scanned directories: 0
Scanned files: 1
Infected files: 0
Data scanned: 0.00 MB
Data read: 0.00 MB (ratio 0.00:1)
Time: 10.153 sec (0 m 10 s)
Start Date: 2024:07:06 20:04:06
End Date:   2024:07:06 20:04:16`

	expected := &clamav.ScanResult{
		KnownViruses:       "0",
		EngineVersion:      "0.103.8",
		ScannedDirectories: "0",
		ScannedFiles:       "1",
		InfectedFiles:      "0",
		DataScanned:        "0.00",
		DataRead:           "0.00",
		Time:               "10.153",
		StartDate:          "2024:07:06 20:04:06",
		EndDate:            "2024:07:06 20:04:16",
	}

	result, err := clamav.Parse(output)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestScan(t *testing.T) {
	file, err := os.Create("/tmp/test.txt")
	require.NoError(t, err)
	defer os.Remove("/tmp/test.txt")
	defer file.Close()

	result, err := clamav.Scan(nil, []string{"/tmp/test.txt"}, false, false)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestCheck(t *testing.T) {
	file, err := os.Create("/tmp/test.txt")
	require.NoError(t, err)
	defer os.Remove("/tmp/test.txt")
	defer file.Close()

	err = clamav.Check([]string{"/tmp/test.txt"}, false)
	require.NoError(t, err)

	os.Mkdir("/tmp/test", 0755)
	defer os.Remove("/tmp/test")
	err = clamav.Check([]string{"/tmp/test"}, true)
	require.NoError(t, err)
	err = clamav.Check([]string{"/tmp/test"}, false)
	require.ErrorIs(t, err, clamav.ErrorPath)
}

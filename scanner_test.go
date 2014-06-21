package scanner_test

import (
  "github.com/Southern/scanner"
  "io/ioutil"
  "strings"
  "testing"
)

var s = make(scanner.Scanner, 0)

func TestScannerReadFile(t *testing.T) {
  Status("Reading all files in testdata directory")
  files, err := ioutil.ReadDir("testdata")

  if err != nil {
    t.Errorf("Unexpected error: %s", err)
    return
  }

  Status("Scanning all files found in testdata directory")
  for len(files) > 0 {
    file := strings.Join([]string{"testdata", files[0].Name()}, "/")
    Status("Scanning file: %s", file)

    err, s = s.ReadFile(file)
    Status("Scanned: %+v", s)
    files = files[1:]
  }
}

func TestScannerNonexistentFile(t *testing.T) {
  Status("Trying to read nonexistent file")
  err, d := s.ReadFile("idontevenexist")

  if len(d) > 0 || err == nil {
    t.Errorf("Expected this test to fail.")
  }
}

func TestJoiningLexBackToString(t *testing.T) {
  err, data := s.ReadFile(strings.Join([]string{"testdata", "html.txt"}, "/"))

  if err != nil {
    t.Errorf("Unexpected error: %s", err)
    return
  }

  Status("Scanned data: %+v", data)

  joined := data.Join()

  Status("Joined data: %+v", joined)
}

func TestInvalidDataType(t *testing.T) {
  Status("Trying to parse an invalid data type")

  err, _ := s.Parse([]int{1, 2, 3, 4})

  Status("Error returned: %s", err)

  if err == nil {
    t.Errorf("Expected this test to fail.")
  }
}

func TestScannerString(t *testing.T) {
  str := "ZoMg Ω≈∂œ™£¢˜Ωπππ¬˜£™¡¢∞•ªº < > & ; ?"
  Status("Trying to parse a string")

  err, s := s.Parse(str)

  if err != nil {
    t.Errorf("Unexpected error: %s", err)
    return
  }

  Status("String was parsed. Joining back together and checking the result.")

  joined := s.Join()

  if joined != str {
    t.Errorf("The joined string was not the same as what was input.")
  }
}

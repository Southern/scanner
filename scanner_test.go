package scanner_test

import (
  "github.com/Southern/scanner"
  "io/ioutil"
  "strings"
  "testing"
)

var s = scanner.New()

func TestScannerBasics(t *testing.T) {
  str := "test-1 test + 1 test+1 -1 1000 -1000"
  expects := [][]string{
    []string{"WORD", "test"},
    []string{"CHAR", "-"},
    []string{"NUMBER", "1"},
    []string{"WHITESPACE", " "},
    []string{"WORD", "test"},
    []string{"WHITESPACE", " "},
    []string{"CHAR", "+"},
    []string{"WHITESPACE", " "},
    []string{"NUMBER", "1"},
    []string{"WHITESPACE", " "},
    []string{"WORD", "test"},
    []string{"CHAR", "+"},
    []string{"NUMBER", "1"},
    []string{"WHITESPACE", " "},
    []string{"CHAR", "-"},
    []string{"NUMBER", "1"},
    []string{"WHITESPACE", " "},
    []string{"NUMBER", "1000"},
    []string{"WHITESPACE", " "},
    []string{"CHAR", "-"},
    []string{"NUMBER", "1000"},
  }
  Status("Parsing \"%s\"", str)

  s, err := s.Parse(str)

  if err != nil {
    t.Errorf("Unexpected error: %s", err)
    return
  }

  for i := 0; i < len(s.Tokens); i++ {
    if s.Tokens[i][0] != expects[i][0] || s.Tokens[i][1] != expects[i][1] {
      t.Errorf("Expected %+v, got %+v", expects[i], s.Tokens[i])
      return
    }
  }

  Status("Parsed: %s", s)
}

func TestScannerManipulation(t *testing.T) {
  str := "test test test"
  expects := [][]string{
    []string{"WORD", "test"},
    []string{"WHITESPACE", " "},
    []string{"WORD", "test2"},
    []string{"WHITESPACE", " "},
    []string{"WORD", "test"},
  }

  Status("Parsing \"%s\"", str)
  s, err := s.Parse(str)

  if err != nil {
    t.Errorf("Unexpected error: %s", err)
    return
  }

  Status("Parsed data: %+v\n", s)
  s.Tokens[2][1] = "test2"
  for i := 0; i < len(s.Tokens); i++ {
    if s.Tokens[i][0] != expects[i][0] || s.Tokens[i][1] != expects[i][1] {
      t.Errorf("Expected %+v, got %+v", expects[i], s.Tokens[i])
      return
    }
  }

  Status("Data after manipulation: %s", s.Join())
}

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

    s, err = s.ReadFile(file)

    if err != nil {
      t.Errorf("Unexpected error: %s", err)
      return
    }

    Status("Scanned: %+v", s)
    files = files[1:]
  }
}

func TestScannerNonexistentFile(t *testing.T) {
  s := scanner.New()
  Status("Trying to read nonexistent file")
  s, err := s.ReadFile("idontevenexist")

  Status("Scan: %+v\n", s)

  if len(s.Tokens) > 0 || err == nil {
    t.Errorf("Expected this test to fail.")
  }
}

func TestJoiningLexBackToString(t *testing.T) {
  data, err := s.ReadFile(strings.Join([]string{"testdata", "html.txt"}, "/"))

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

  _, err := s.Parse([]int{1, 2, 3, 4})

  Status("Error returned: %s", err)

  if err == nil {
    t.Errorf("Expected this test to fail.")
  }
}

func TestScannerRandomString(t *testing.T) {
  str := "ZoMg testΩ≈∂œ™£¢˜Ωπππ¬˜£™¡¢∞•ªº test< > & ; ?"
  Status("Trying to parse random string")

  s, err := s.Parse(str)

  if err != nil {
    t.Errorf("Unexpected error: %s", err)
    return
  }

  Status("Scanned: %+v", s)

  Status("String was parsed. Joining back together and checking the result.")

  joined := s.Join()

  if joined != str {
    t.Errorf("The joined string was not the same as what was input.")
    return
  }

  Status("String parsed correctly.")
}

func TestScannerRussianString(t *testing.T) {
  str := "This isn't Russian, but this is: ру́сский язы́к"
  Status("Trying to parse Russian string")

  s, err := s.Parse(str)

  if err != nil {
    t.Errorf("Unexpected error: %s", err)
    return
  }

  Status("Scanned: %+v", s)

  Status("String was parsed. Joining back together and checking the result.")

  joined := s.Join()

  if joined != str {
    t.Errorf("The joined string was not the same as what was input.")
    return
  }

  Status("String parsed correctly.")
}

func TestScannerGreekString(t *testing.T) {
  str := "This isn't Greek, but this is: ελληνικά"
  Status("Trying to parse Greek string")

  s, err := s.Parse(str)

  if err != nil {
    t.Errorf("Unexpected error: %s", err)
    return
  }

  Status("Scanned: %+v", s)

  Status("String was parsed. Joining back together and checking the result.")

  joined := s.Join()

  if joined != str {
    t.Errorf("The joined string was not the same as what was input.")
    return
  }

  Status("String parsed correctly")
}

func TestScannerArabicString(t *testing.T) {
  str := "This isn't Arabic, but this is: عربي ,عربى"
  Status("Trying to parse Arabic string")

  s, err := s.Parse(str)

  if err != nil {
    t.Errorf("Unexpected error: %s", err)
    return
  }

  Status("Scanned: %+v", s)

  Status("String was parsed. Joining back together and checking the result.")

  joined := s.Join()

  if joined != str {
    t.Errorf("The joined string was not the same as what was input.")
    return
  }

  Status("String parsed correctly")
}

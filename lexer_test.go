package lexer_test

import (
  "github.com/Southern/lexer"
  "io/ioutil"
  "os"
  "strings"
  "testing"
)

var l = make(lexer.Lexer, 0)

func TestLexerReadFile(t *testing.T) {
  Status("Getting current working directory")
  cd, err := os.Getwd()

  if err != nil {
    panic("Could not get working directory")
    return
  }

  Status("Reading all files in testdata directory")
  files, err := ioutil.ReadDir(strings.Join([]string{cd, "testdata"}, "/"))

  if err != nil {
    t.Errorf("Unexpected error: %s", err)
    return
  }

  Status("Lexing all files found in testdata directory")
  for len(files) > 0 {
    file := strings.Join([]string{cd, "testdata", files[0].Name()}, "/")
    Status("Lexing file: %s", file)

    err, l = l.ReadFile(file)
    Status("Lexed: %+v", l)
    files = files[1:]
  }
}

func TestLexerNonexistentFile(t *testing.T) {
  Status("Trying to read nonexistent file")
  err, d := l.ReadFile("idontevenexist")

  if len(d) > 0 || err == nil {
    t.Errorf("Expected this test to fail.")
  }
}

func TestJoiningLexBackToString(t *testing.T) {
  Status("Getting current working directory")
  cd, err := os.Getwd()

  if err != nil {
    panic("Could not get working directory")
    return
  }

  err, data := l.ReadFile(strings.Join([]string{cd, "testdata", "html.txt"}, "/"))

  if err != nil {
    t.Errorf("Unexpected error: %s", err)
    return
  }

  Status("Lexed data: %+v", data)

  joined := data.Join()

  Status("Joined data: %+v", joined)
}

func TestInvalidDataType(t *testing.T) {
  Status("Trying to parse an invalid data type")

  err, _ := l.Parse([]int{1, 2, 3, 4})

  Status("Error returned: %s", err)

  if err == nil {
    t.Errorf("Expected this test to fail.")
  }
}

func TestLexerString(t *testing.T) {
  str := "ZoMg Ω≈∂œ™£¢˜Ωπππ¬˜£™¡¢∞•ªº < > & ; ?"
  Status("Trying to parse a string")

  err, l := l.Parse(str)

  if err != nil {
    t.Errorf("Unexpected error: %s", err)
    return
  }

  Status("String was parsed. Joining back together and checking the result.")

  joined := l.Join()

  if joined != str {
    t.Errorf("The joined string was not the same as what was input.")
  }
}

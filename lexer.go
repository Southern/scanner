/*

  Lexer that will break down every character, word, whitespace, and number
  that is in a file.

*/
package lexer

import (
  "fmt"
  "io/ioutil"
  "regexp"
)

// Lexer definitions that are added to LexerMap
type LexerDefinition struct {
  Regex *regexp.Regexp
  Type  string
}

// LexerMap holds the regexes that we are wanting to match throughout the
// file. It comes preloaded with the
var LexerMap = []LexerDefinition{
  LexerDefinition{regexp.MustCompile("^(?i)[a-z][a-z0-9]+"), "WORD"},
  LexerDefinition{regexp.MustCompile("^\\s+"), "WHITESPACE"},
  LexerDefinition{regexp.MustCompile("^(?i)([a-z]|[^0-9])"), "CHAR"},
  LexerDefinition{regexp.MustCompile("^[0-9]+"), "NUMBER"},
}

// Lexer is just a double array of strings.
type Lexer [][]string

/*

  Join data back together once it has been lexed, if it is desired.

  For an example, let's look at the test for this function:
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

  This will allow you to do whatever you want with your lexed data, and then
  simply join it back.

*/
func (l Lexer) Join() string {
  var joined string

  for len(l) > 0 {
    joined = joined + l[0][1]
    l = l[1:]
  }

  return joined
}

func (l Lexer) Parse(data interface{}) (error, Lexer) {
  l = make([][]string, 0)

  switch data.(type) {
  case []byte:
    data = string(data.([]byte))
  case string:
    data = data.(string)
  default:
    return fmt.Errorf("Lexer.Parse only accepts []byte and string types."), l
  }

  for len(data.(string)) > 0 {
    for _, def := range LexerMap {
      r, t := def.Regex, def.Type
      str := r.FindString(data.(string))
      if len(str) > 0 {
        l = append(l, []string{t, str})
        data = data.(string)[len(str):]
        break
      }
    }
  }

  return nil, l
}

func (l Lexer) ReadFile(filename string) (error, Lexer) {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    return err, nil
  }

  return l.Parse(data)
}

/*

  Scanner that will break down every character, word, whitespace, and number
  that is passed into it.

*/
package scanner

import (
  "fmt"
  "io/ioutil"
  "regexp"
)

// Scanner definitions that are added to ScannerMap
type ScannerDefinition struct {
  Regex *regexp.Regexp
  Type  string
}

/*

ScannerMap holds the regexes that we are wanting to match throughout the
file. It also contains the type that we want to return if the regex is
matched.

*/
var ScannerMap = []ScannerDefinition{
  ScannerDefinition{regexp.MustCompile("^(?i)[a-z][a-z0-9]+"), "WORD"},
  ScannerDefinition{regexp.MustCompile("^\\s+"), "WHITESPACE"},
  ScannerDefinition{regexp.MustCompile("^(?i)([a-z]|[^0-9])"), "CHAR"},
  ScannerDefinition{regexp.MustCompile("^[0-9]+"), "NUMBER"},
}

// Scanner is just a double array of strings.
type Scanner [][]string

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
func (s Scanner) Join() string {
  var joined string

  for len(s) > 0 {
    joined = joined + s[0][1]
    s = s[1:]
  }

  return joined
}

func (s Scanner) Parse(data interface{}) (error, Scanner) {
  s = make([][]string, 0)

  switch data.(type) {
  case []byte:
    data = string(data.([]byte))
  case string:
    data = data.(string)
  default:
    return fmt.Errorf("Scanner.Parse only accepts []byte and string types."), s
  }

  for len(data.(string)) > 0 {
    for _, def := range ScannerMap {
      r, t := def.Regex, def.Type
      str := r.FindString(data.(string))
      if len(str) > 0 {
        s = append(s, []string{t, str})
        data = data.(string)[len(str):]
        break
      }
    }
  }

  return nil, s
}

func (s Scanner) ReadFile(filename string) (error, Scanner) {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    return err, nil
  }

  return s.Parse(data)
}

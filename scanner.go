/*

  Scanner that will break down every character, word, whitespace, and number
  that is passed into it.

*/
package scanner

import (
  "fmt"
  "io/ioutil"
  "regexp"
  "strings"
)

// Scanner definitions that are added to ScannerMap
type ScannerDefinition struct {
  Regex *regexp.Regexp
  Type  string
}

/*

  Accepted Unicode characters. This is exposed so that you, the end user, may
  add your own if the need arises.

  A simple way to grab all UTF-16 blocks, in case these ever get deleted for
  some reason, is to use:
    curl http://jrgraphix.net/research/unicode_blocks.php | \
    grep -iE "[a-f0-9]{4,} — [a-f0-9]{4,}" | \
    awk '{gsub("&nbsp;", " ", $0); gsub(/ {2,}/, "\n", $0); print $0}' | \
    awk '/<tr>/ {next} {FS="[<>\"]"; split(substr($7, 12), range, "-"); if (length($9) == 0) next; else printf "  // %s\n  \"\\\\x{%s}-\\\\x{%s}\",\n\n", $9, range[1], range[2]}'

  Then you can just copy it to your clipboard with another pipe and paste it
  here. Simple as that.

*/
var Unicode = []string{
  // Block Elements
  "\\x{2580}-\\x{259F}",

  // Latin-1 Supplement
  "\\x{00A0}-\\x{00FF}",

  // Geometric Shapes
  "\\x{25A0}-\\x{25FF}",

  // Latin Extended-A
  "\\x{0100}-\\x{017F}",

  // Miscellaneous Symbols
  "\\x{2600}-\\x{26FF}",

  // Latin Extended-B
  "\\x{0180}-\\x{024F}",

  // Dingbats
  "\\x{2700}-\\x{27BF}",

  // IPA Extensions
  "\\x{0250}-\\x{02AF}",

  // Miscellaneous Mathematical Symbols-A
  "\\x{27C0}-\\x{27EF}",

  // Spacing Modifier Letters
  "\\x{02B0}-\\x{02FF}",

  // Supplemental Arrows-A
  "\\x{27F0}-\\x{27FF}",

  // Combining Diacritical Marks
  "\\x{0300}-\\x{036F}",

  // Braille Patterns
  "\\x{2800}-\\x{28FF}",

  // Greek and Coptic
  "\\x{0370}-\\x{03FF}",

  // Supplemental Arrows-B
  "\\x{2900}-\\x{297F}",

  // Cyrillic
  "\\x{0400}-\\x{04FF}",

  // Miscellaneous Mathematical Symbols-B
  "\\x{2980}-\\x{29FF}",

  // Cyrillic Supplementary
  "\\x{0500}-\\x{052F}",

  // Supplemental Mathematical Operators
  "\\x{2A00}-\\x{2AFF}",

  // Armenian
  "\\x{0530}-\\x{058F}",

  // Miscellaneous Symbols and Arrows
  "\\x{2B00}-\\x{2BFF}",

  // Hebrew
  "\\x{0590}-\\x{05FF}",

  // CJK Radicals Supplement
  "\\x{2E80}-\\x{2EFF}",

  // Arabic
  "\\x{0600}-\\x{06FF}",

  // Kangxi Radicals
  "\\x{2F00}-\\x{2FDF}",

  // Syriac
  "\\x{0700}-\\x{074F}",

  // Ideographic Description Characters
  "\\x{2FF0}-\\x{2FFF}",

  // Thaana
  "\\x{0780}-\\x{07BF}",

  // CJK Symbols and Punctuation
  "\\x{3000}-\\x{303F}",

  // Devanagari
  "\\x{0900}-\\x{097F}",

  // Hiragana
  "\\x{3040}-\\x{309F}",

  // Bengali
  "\\x{0980}-\\x{09FF}",

  // Katakana
  "\\x{30A0}-\\x{30FF}",

  // Gurmukhi
  "\\x{0A00}-\\x{0A7F}",

  // Bopomofo
  "\\x{3100}-\\x{312F}",

  // Gujarati
  "\\x{0A80}-\\x{0AFF}",

  // Hangul Compatibility Jamo
  "\\x{3130}-\\x{318F}",

  // Oriya
  "\\x{0B00}-\\x{0B7F}",

  // Kanbun
  "\\x{3190}-\\x{319F}",

  // Tamil
  "\\x{0B80}-\\x{0BFF}",

  // Bopomofo Extended
  "\\x{31A0}-\\x{31BF}",

  // Telugu
  "\\x{0C00}-\\x{0C7F}",

  // Katakana Phonetic Extensions
  "\\x{31F0}-\\x{31FF}",

  // Kannada
  "\\x{0C80}-\\x{0CFF}",

  // Enclosed CJK Letters and Months
  "\\x{3200}-\\x{32FF}",

  // Malayalam
  "\\x{0D00}-\\x{0D7F}",

  // CJK Compatibility
  "\\x{3300}-\\x{33FF}",

  // Sinhala
  "\\x{0D80}-\\x{0DFF}",

  // CJK Unified Ideographs Extension A
  "\\x{3400}-\\x{4DBF}",

  // Thai
  "\\x{0E00}-\\x{0E7F}",

  // Yijing Hexagram Symbols
  "\\x{4DC0}-\\x{4DFF}",

  // Lao
  "\\x{0E80}-\\x{0EFF}",

  // CJK Unified Ideographs
  "\\x{4E00}-\\x{9FFF}",

  // Tibetan
  "\\x{0F00}-\\x{0FFF}",

  // Yi Syllables
  "\\x{A000}-\\x{A48F}",

  // Myanmar
  "\\x{1000}-\\x{109F}",

  // Yi Radicals
  "\\x{A490}-\\x{A4CF}",

  // Georgian
  "\\x{10A0}-\\x{10FF}",

  // Hangul Syllables
  "\\x{AC00}-\\x{D7AF}",

  // Hangul Jamo
  "\\x{1100}-\\x{11FF}",

  // High Surrogates
  "\\x{D800}-\\x{DB7F}",

  // Ethiopic
  "\\x{1200}-\\x{137F}",

  // High Private Use Surrogates
  "\\x{DB80}-\\x{DBFF}",

  // Cherokee
  "\\x{13A0}-\\x{13FF}",

  // Low Surrogates
  "\\x{DC00}-\\x{DFFF}",

  // Unified Canadian Aboriginal Syllabics
  "\\x{1400}-\\x{167F}",

  // Private Use Area
  "\\x{E000}-\\x{F8FF}",

  // Ogham
  "\\x{1680}-\\x{169F}",

  // CJK Compatibility Ideographs
  "\\x{F900}-\\x{FAFF}",

  // Runic
  "\\x{16A0}-\\x{16FF}",

  // Alphabetic Presentation Forms
  "\\x{FB00}-\\x{FB4F}",

  // Tagalog
  "\\x{1700}-\\x{171F}",

  // Arabic Presentation Forms-A
  "\\x{FB50}-\\x{FDFF}",

  // Hanunoo
  "\\x{1720}-\\x{173F}",

  // Variation Selectors
  "\\x{FE00}-\\x{FE0F}",

  // Buhid
  "\\x{1740}-\\x{175F}",

  // Combining Half Marks
  "\\x{FE20}-\\x{FE2F}",

  // Tagbanwa
  "\\x{1760}-\\x{177F}",

  // CJK Compatibility Forms
  "\\x{FE30}-\\x{FE4F}",

  // Khmer
  "\\x{1780}-\\x{17FF}",

  // Small Form Variants
  "\\x{FE50}-\\x{FE6F}",

  // Mongolian
  "\\x{1800}-\\x{18AF}",

  // Arabic Presentation Forms-B
  "\\x{FE70}-\\x{FEFF}",

  // Limbu
  "\\x{1900}-\\x{194F}",

  // Halfwidth and Fullwidth Forms
  "\\x{FF00}-\\x{FFEF}",

  // Tai Le
  "\\x{1950}-\\x{197F}",

  // Specials
  "\\x{FFF0}-\\x{FFFF}",

  // Khmer Symbols
  "\\x{19E0}-\\x{19FF}",

  // Linear B Syllabary
  "\\x{10000}-\\x{1007F}",

  // Phonetic Extensions
  "\\x{1D00}-\\x{1D7F}",

  // Linear B Ideograms
  "\\x{10080}-\\x{100FF}",

  // Latin Extended Additional
  "\\x{1E00}-\\x{1EFF}",

  // Aegean Numbers
  "\\x{10100}-\\x{1013F}",

  // Greek Extended
  "\\x{1F00}-\\x{1FFF}",

  // Old Italic
  "\\x{10300}-\\x{1032F}",

  // General Punctuation
  "\\x{2000}-\\x{206F}",

  // Gothic
  "\\x{10330}-\\x{1034F}",

  // Superscripts and Subscripts
  "\\x{2070}-\\x{209F}",

  // Ugaritic
  "\\x{10380}-\\x{1039F}",

  // Currency Symbols
  "\\x{20A0}-\\x{20CF}",

  // Deseret
  "\\x{10400}-\\x{1044F}",

  // Combining Diacritical Marks for Symbols
  "\\x{20D0}-\\x{20FF}",

  // Shavian
  "\\x{10450}-\\x{1047F}",

  // Letterlike Symbols
  "\\x{2100}-\\x{214F}",

  // Osmanya
  "\\x{10480}-\\x{104AF}",

  // Number Forms
  "\\x{2150}-\\x{218F}",

  // Cypriot Syllabary
  "\\x{10800}-\\x{1083F}",

  // Arrows
  "\\x{2190}-\\x{21FF}",

  // Byzantine Musical Symbols
  "\\x{1D000}-\\x{1D0FF}",

  // Mathematical Operators
  "\\x{2200}-\\x{22FF}",

  // Musical Symbols
  "\\x{1D100}-\\x{1D1FF}",

  // Miscellaneous Technical
  "\\x{2300}-\\x{23FF}",

  // Tai Xuan Jing Symbols
  "\\x{1D300}-\\x{1D35F}",

  // Control Pictures
  "\\x{2400}-\\x{243F}",

  // Mathematical Alphanumeric Symbols
  "\\x{1D400}-\\x{1D7FF}",

  // Optical Character Recognition
  "\\x{2440}-\\x{245F}",

  // CJK Unified Ideographs Extension B
  "\\x{20000}-\\x{2A6DF}",

  // Enclosed Alphanumerics
  "\\x{2460}-\\x{24FF}",

  // CJK Compatibility Ideographs Supplement
  "\\x{2F800}-\\x{2FA1F}",

  // Box Drawing
  "\\x{2500}-\\x{257F}",

  // Tags
  "\\x{E0000}-\\x{E007F}",
}

func unicode() string {
  return strings.Join(Unicode, "")
}

/*

ScannerMap holds the regexes that we are wanting to match throughout the
file. It also contains the type that we want to return if the regex is
matched.

*/
var ScannerMap = []ScannerDefinition{
  ScannerDefinition{regexp.MustCompile(
    fmt.Sprintf("^(?i)([a-z0-9][a-z0-9\\-'%s]+|[%s]+)",
      unicode(), unicode()),
  ), "WORD"},
  ScannerDefinition{regexp.MustCompile("^\\s+"), "WHITESPACE"},
  ScannerDefinition{regexp.MustCompile(
    "^(?i)([a-z]|[^0-9])",
  ), "CHAR"},
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

      err, data := s.ReadFile(strings.Join([]string{cd, "testdata", "html.txt"}, "/"))

      if err != nil {
        t.Errorf("Unexpected error: %s", err)
        return
      }

      Status("Scanned data: %+v", data)

      joined := data.Join()

      Status("Joined data: %+v", joined)
    }

  This will allow you to do whatever you want with your scanned data, and then
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

/*

  Parse []byte or string into their most basic forms. "WORD", "CHAR",
  "NUMBER", and "WHITESPACE".

  Once the data is parsed, you can manipulate the data and completely change
  the original data.

  For instance, taking a look at our test for this:
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
      err, s := s.Parse(str)

      if err != nil {
        t.Errorf("Unexpected error: %s", err)
        return
      }

      Status("Parsed data: %+v\n", s)
      s[2][1] = "test2"
      for i := 0; i < len(s); i++ {
        if s[i][0] != expects[i][0] || s[i][1] != expects[i][1] {
          t.Errorf("Manipulation failed.")
          return
        }
      }

      Status("Data after manipulation: %s", s.Join())
    }

*/
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

/*

  Reads a file and automatically runs it through Parse, returning the parsed
  results of the file that was read.

*/
func (s Scanner) ReadFile(filename string) (error, Scanner) {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    return err, nil
  }

  return s.Parse(data)
}

/*
Copyright 2012 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shlex

import (
	"strings"
	"testing"
)

var (
	// one two "three four" "five \"six\"" seven#eight # nine # ten
	// eleven 'twelve\'
	testString          = "one two \"three four\" \"five \\\"six\\\"\" seven#eight # nine # ten\n eleven 'twelve\\' thirteen=13 fourteen/14"
	testStringBackPort  = `run --command "hello \"some\value\\another\" there" finally`
	testStringBackPort2 = `dir_list "c:\\" *`
	testStringBackPort3 = `dir_list "c:\\ 'yup'" *`
	testStringBackPort4 = `dir_list "c:\\ \"yup\"" *`
	testStringBackPort5 = `log_get --file 'c:\temp\data.json' --type json`
	testStringBackPort6 = `log_get --file c:\temp\data.json --type json`
	testStringBackPort7 = `dir_list "c:\\ \'yup\'" *`
	testStringBackPort8 = `history_dump`
)

func TestClassifier(t *testing.T) {
	classifier := newDefaultClassifier()
	tests := map[rune]runeTokenClass{
		' ':  spaceRuneClass,
		'"':  escapingQuoteRuneClass,
		'\'': nonEscapingQuoteRuneClass,
		'#':  commentRuneClass}
	for runeChar, want := range tests {
		got := classifier.ClassifyRune(runeChar)
		if got != want {
			t.Errorf("ClassifyRune(%v) -> %v. Want: %v", runeChar, got, want)
		}
	}
}

func TestTokenizer(t *testing.T) {
	testInput := strings.NewReader(testString)
	expectedTokens := []*Token{
		&Token{WordToken, "one"},
		&Token{WordToken, "two"},
		&Token{WordToken, "three four"},
		&Token{WordToken, "five \"six\""},
		&Token{WordToken, "seven#eight"},
		&Token{CommentToken, " nine # ten"},
		&Token{WordToken, "eleven"},
		&Token{WordToken, "twelve\\"},
		&Token{WordToken, "thirteen=13"},
		&Token{WordToken, "fourteen/14"}}

	tokenizer := NewTokenizer(testInput)
	for i, want := range expectedTokens {
		got, err := tokenizer.Next()
		if err != nil {
			t.Error(err)
		}
		if !got.Equal(want) {
			t.Errorf("Tokenizer.Next()[%v] of %q -> %v. Want: %v", i, testString, got, want)
		}
	}
}

func TestLexer(t *testing.T) {
	testInput := strings.NewReader(testString)
	expectedStrings := []string{"one", "two", "three four", "five \"six\"", "seven#eight", "eleven", "twelve\\", "thirteen=13", "fourteen/14"}

	lexer := NewLexer(testInput)
	for i, want := range expectedStrings {
		got, err := lexer.Next()
		if err != nil {
			t.Error(err)
		}
		if got != want {
			t.Errorf("Lexer.Next()[%v] of %q -> %v. Want: %v", i, testString, got, want)
		}
	}
}

func TestSplit(t *testing.T) {
	want := []string{"one", "two", "three four", "five \"six\"", "seven#eight", "eleven", "twelve\\", "thirteen=13", "fourteen/14"}
	got, err := Split(testString)
	if err != nil {
		t.Error(err)
	}
	if len(want) != len(got) {
		t.Errorf("Split(%q) -> %v. Want: %v", testString, got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Split(%q)[%v] -> %v. Want: %v", testString, i, got[i], want[i])
		}
	}
}

func TestSplitBackSupport(t *testing.T) {
	want := []string{"run", "--command", "hello \"some\\value\\another\" there", "finally"}
	got, err := Split(testStringBackPort)
	if err != nil {
		t.Error(err)
	}
	if len(want) != len(got) {
		t.Errorf("Split(%q) -> %#v. Want: %#v", testStringBackPort, got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Split(%q)[%v]\n->\n%#v.\nWant:\n%#v", testStringBackPort, i, got[i], want[i])
		}
	}
}

func TestSplitBackSupport2(t *testing.T) {
	want := []string{"dir_list", "c:\\", "*"}
	got, err := Split(testStringBackPort2)
	if err != nil {
		t.Error(err)
	}
	if len(want) != len(got) {
		t.Errorf("Split(%q) -> %#v. Want: %#v", testStringBackPort2, got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Split(%q)[%v]\n->\n%#v.\nWant:\n%#v", testStringBackPort2, i, got[i], want[i])
		}
	}
}

func TestSplitBackSupport3(t *testing.T) {
	want := []string{"dir_list", "c:\\ 'yup'", "*"}
	got, err := Split(testStringBackPort3)
	if err != nil {
		t.Error(err)
	}
	if len(want) != len(got) {
		t.Errorf("Split(%q) -> %#v. Want: %#v", testStringBackPort3, got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Split(%q)[%v]\n->\n%#v.\nWant:\n%#v", testStringBackPort3, i, got[i], want[i])
		}
	}
}

func TestSplitBackSupport4(t *testing.T) {
	want := []string{"dir_list", "c:\\ \"yup\"", "*"}
	got, err := Split(testStringBackPort4)
	if err != nil {
		t.Error(err)
	}
	if len(want) != len(got) {
		t.Errorf("Split(%q) -> %#v. Want: %#v", testStringBackPort4, got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Split(%q)[%v]\n->\n%#v.\nWant:\n%#v", testStringBackPort4, i, got[i], want[i])
		}
	}
}

func TestSplitBackSupport5(t *testing.T) {
	want := []string{"log_get", "--file", "c:\\temp\\data.json", "--type", "json"}
	got, err := Split(testStringBackPort5)
	if err != nil {
		t.Error(err)
	}
	if len(want) != len(got) {
		t.Errorf("Split(%q) -> %#v. Want: %#v", testStringBackPort5, got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Split(%q)[%v]\n->\n%#v.\nWant:\n%#v", testStringBackPort5, i, got[i], want[i])
		}
	}
}

func TestSplitBackSupport6(t *testing.T) {
	want := []string{"log_get", "--file", "c:\\temp\\data.json", "--type", "json"}
	got, err := Split(testStringBackPort6)
	if err != nil {
		t.Error(err)
	}
	if len(want) != len(got) {
		t.Errorf("Split(%q) -> %#v. Want: %#v", testStringBackPort6, got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Split(%q)[%v]\n->\n%#v.\nWant:\n%#v", testStringBackPort6, i, got[i], want[i])
		}
	}
}

func TestSplitBackSupport7(t *testing.T) {
	want := []string{"dir_list", "c:\\ 'yup'", "*"}
	got, err := Split(testStringBackPort7)
	if err != nil {
		t.Error(err)
	}
	if len(want) != len(got) {
		t.Errorf("Split(%q) -> %#v. Want: %#v", testStringBackPort7, got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Split(%q)[%v]\n->\n%#v.\nWant:\n%#v", testStringBackPort7, i, got[i], want[i])
		}
	}
}

func TestSplitBackSupport8(t *testing.T) {
	want := []string{"history_dump"}
	got, err := Split(testStringBackPort8)
	if err != nil {
		t.Error(err)
	}
	if len(want) != len(got) {
		t.Errorf("Split(%q) -> %#v. Want: %#v", testStringBackPort8, got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Split(%q)[%v]\n->\n%#v.\nWant:\n%#v", testStringBackPort8, i, got[i], want[i])
		}
	}
}

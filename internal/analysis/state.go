package analysis

import (
	"fmt"
	"log"
	"slices"

	"github.com/exgene/cuda-autocompletes/internal/lsp"
)

type Diffs struct {
	StartRange Range
	EndRange   Range
	Text       string
}

type Range struct {
	Line      int
	Character int
}

func NewDiff(change lsp.TextDocumentContentChangeEvent) Diffs {
	return Diffs{
		StartRange: Range{
			Line:      change.Range.Start.Line,
			Character: change.Range.Start.Character,
		},
		EndRange: Range{
			Line:      change.Range.End.Line,
			Character: change.Range.End.Character,
		},
		Text: change.Text,
	}
}

type State struct {
	Documents     map[string][]string
	CurrentBuffer string
}

func NewState() State {
	return State{
		Documents:     map[string][]string{},
		CurrentBuffer: "",
	}
}

func (s *State) OpenDocument(document, data string) {
	s.Documents[document] = parseInput(data)
	s.CurrentBuffer = unwind(s.Documents[document])
}

// Leetcode type shi (Two Pointer Fine Shi)
func (s *State) GetToken(document string, position lsp.Position) string {
	lineIdx := position.Line
	characterIdx := position.Character
	line := s.Documents[document][lineIdx]
	var token []rune
	i := characterIdx
	j := characterIdx

	token = append(token, rune(line[i]))

	// Advance first before consuming
	i -= 1
	j += 1

	for i >= 0 && !isEscapeCharacter(rune(line[i])) {
		log.Printf("i_value:%s, i_index:%d", string(line[i]), i)
		token = slices.Insert(token, 0, rune(line[i]))
		i -= 1
	}

	for j < len(line) && !isEscapeCharacter(rune(line[j])) {
		log.Printf("j_value:%s, j_index:%d", string(line[j]), j)
		token = append(token, rune(line[j]))
		j += 1
	}

	return processed(string(token))
}

// Future support for escaped strings ig, Most probably not needed
// TODO: Deprecate this shit
func processed(t string) string {
	switch t {
	case "#":
		return `\#`
	default:
		return t
	}
}

func isEscapeCharacter(l rune) bool {
	if l == '\n' || l == '\r' || l == ' ' {
		return true
	}
	return false
}

func unwind(document []string) string {
	constructedString := ""
	for _, line := range document {
		constructedString += line
		constructedString += "\n"
	}
	return constructedString
}

func parseInput(text string) []string {
	output := []string{}
	initialPos := 0
	for i, char := range text {
		if char == '\n' {
			// I want to skip empty lines, But we cant do that cause it will mess up the order, Also Edits, Does it include \n characters?
			// if i+1 == initialPos {
			// 	continue
			// }
			output = append(output, text[initialPos:i])
			initialPos = i + 1
		}
	}
	if len(output) == 0 {
		return append(output, text)
	}
	return output
}

func (s *State) ApplyDiffs(document string, diffs Diffs) error {
	log.Printf("Running Diffs")
	var currentDocument []string
	currentDocument, ok := s.Documents[document]
	if !ok {
		return fmt.Errorf("Document Not found %s", document)
	}

	if diffs.Text == "" {
		updatedDocument, err := deleteTextFromDocument(currentDocument, diffs)
		if err != nil {
			return fmt.Errorf("Error while deleting items from document: %s", document)
		}
		s.Documents[document] = updatedDocument
	} else {
		updatedDocument, err := updateTextFromDocument(currentDocument, diffs)
		if err != nil {
			return fmt.Errorf("Error while updating items from document: %s", document)
		}
		s.Documents[document] = updatedDocument
	}
	// Set the current buffer to the updated document (unwinded)
	s.CurrentBuffer = unwind(s.Documents[document])
	return nil
}

func addLineIntoDocument(currentDocument []string, index int) ([]string, error) {
	if index < 0 {
		index = 0
	}
	return slices.Insert(currentDocument, index, ""), nil
}

func updateTextFromDocument(currentDocument []string, diffs Diffs) ([]string, error) {
	startCharacter := diffs.StartRange.Character
	endCharacter := diffs.EndRange.Character
	startLine := diffs.StartRange.Line
	endLine := diffs.EndRange.Line

	log.Printf("%d == %d, %d == 18, %d == 0", startLine, endLine-1, len(currentDocument[startLine]), endCharacter)
	// Inserting new line characters
	log.Printf("%q", diffs.Text)
	if diffs.Text == "\n" || diffs.Text == "\r\n" || diffs.Text == "\n\n" {
		log.Print("Inside \\n")
		// insert below
		if startLine == endLine &&
			startCharacter == len(currentDocument[startLine]) &&
			endCharacter == len(currentDocument[startLine]) {
			log.Print("Below 1")

			return addLineIntoDocument(currentDocument, endLine+1)
		}

		log.Printf("%d == %d, %d == 18, %d == 0", startLine, endLine-1, len(currentDocument[startLine]), endCharacter)
		if startLine == endLine-1 && startCharacter == len(currentDocument[startLine]) && endCharacter == 0 {
			log.Print("Below 2")
			return addLineIntoDocument(currentDocument, endLine+1)
		}
		// insert above
		if startLine == endLine-1 && endCharacter == len(currentDocument[endCharacter]) && startCharacter == 0 {
			return addLineIntoDocument(currentDocument, startLine)
		}
		if startLine == 0 && endLine == 0 && endCharacter == 0 && startCharacter == 0 {
			return addLineIntoDocument(currentDocument, startLine)
		}
	}

	if diffs.Text == "\n" && startLine < endLine {
		// TODO: Implement this
		log.Panic("This is for multi line additions and i give up for now")
	}

	if startLine == endLine {
		lineToBeModified := currentDocument[startLine]
		prefix := lineToBeModified[0:startCharacter]
		middlePart := diffs.Text
		suffix := lineToBeModified[endCharacter:]
		currentDocument[startLine] = prefix + middlePart + suffix
		return currentDocument, nil
	}

	return currentDocument, nil
}

func deleteTextFromDocument(currentDocument []string, diffs Diffs) ([]string, error) {
	startCharacter := diffs.StartRange.Character
	endCharacter := diffs.EndRange.Character
	startLine := diffs.StartRange.Line
	endLine := diffs.EndRange.Line

	// Dud Changes, TODO: Maybe bubble this up the chain to avoid overhead later
	if startLine == endLine && startCharacter == endCharacter {
		return currentDocument, nil
	}

	// Delete whole set of lines
	if startLine < endLine && startCharacter == 0 && endCharacter == 0 {
		return slices.Delete(currentDocument, startLine, endLine), nil
	}

	// Deletes Characters (In betweens)
	if startLine == endLine && startCharacter < endCharacter {
		n := len(currentDocument[startLine])
		s1 := currentDocument[startLine][0:startCharacter]
		s2 := currentDocument[startLine][endCharacter:n]
		currentDocument[startLine] = s1 + s2
		return currentDocument, nil
	}

	// Other operations pending...
	return currentDocument, nil
}

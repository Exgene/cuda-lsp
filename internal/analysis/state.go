package analysis

import (
	"fmt"
	"log"
	"slices"
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

func NewDiff() Diffs {
	return Diffs{
		StartRange: Range{
			Line:      0,
			Character: 0,
		},
		EndRange: Range{
			Line:      0,
			Character: 0,
		},
		Text: "",
	}
}

type State struct {
	Documents map[string][]string
}

func NewState() State {
	return State{
		Documents: map[string][]string{},
	}
}

func (s *State) OpenDocument(document, data string) {
	s.Documents[document] = parseInput(data)
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
	log.Printf("ITS HERE MARDOIRAI")
	// Currently lets only support single line diffs
	if diffs.StartRange.Line != diffs.EndRange.Line {
		return fmt.Errorf("Not implemented MultiLine Diffing:: StartRange: %d <==> EndRange: %d", diffs.StartRange.Line, diffs.EndRange.Line)
	}

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
		// TODO: Check if it was already empty, State management is so fucking hard
		if len(updatedDocument) == 0 {
			updatedDocument = slices.Delete(updatedDocument, diffs.EndRange.Line, diffs.EndRange.Line)
		}
		s.Documents[document] = updatedDocument
	} else {
		updatedDocument, err := updateTextFromDocument(currentDocument, diffs)
		if err != nil {
			return fmt.Errorf("Error while updating items from document: %s", document)
		}
		s.Documents[document] = updatedDocument
	}
	return nil
}

func updateTextFromDocument(currentDocument []string, diffs Diffs) ([]string, error) {
	startCharacter := diffs.StartRange.Character
	endCharacter := diffs.EndRange.Character
	line := diffs.StartRange.Line

	if line < 0 || line > len(currentDocument) {
		return nil, fmt.Errorf("Out of Bounds")
	}
	lineToBeModified := currentDocument[line]

	prefix := lineToBeModified[0:startCharacter]
	middlePart := diffs.Text
	suffix := lineToBeModified[endCharacter:]
	currentDocument[line] = prefix + middlePart + suffix
	return currentDocument, nil
}

func deleteTextFromDocument(currentDocument []string, diffs Diffs) ([]string, error) {
	startCharacter := diffs.StartRange.Character
	endCharacter := diffs.EndRange.Character
	line := diffs.StartRange.Line

	if line < 0 || line > len(currentDocument) {
		return nil, fmt.Errorf("Out of Bounds")
	}
	lineToBeModified := currentDocument[line]
	n := len(lineToBeModified) - 1
	log.Printf("Line to be modified %s", lineToBeModified)

	currentDocument[line] = lineToBeModified[0:startCharacter] + lineToBeModified[endCharacter:n]
	return currentDocument, nil
}

package analysis

import (
	"errors"
	"fmt"
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
	s.Documents[document] = append(s.Documents[document], data)
}

func (s *State) ApplyDiffs(document string, diffs Diffs) error {
	if diffs.StartRange.Line != diffs.EndRange.Line {
		return fmt.Errorf("Not implemented MultiLine Diffing:: StartRange: %d <==> EndRange: %d", diffs.StartRange.Line, diffs.EndRange.Line)
	}

	var currentDocument []string
	currentDocument, ok := s.Documents[document]
	if !ok {
		return fmt.Errorf("Document Not found %s", document)
	}
	// Currently lets only support single line diffs
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

	prefix := lineToBeModified[0:startCharacter]
	suffix := lineToBeModified[endCharacter:]
	currentDocument[line] = prefix + suffix
	return currentDocument, nil
}

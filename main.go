package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func NewRDFItem(subject string) *RDFItem {
	return &RDFItem{
		Subject:  subject,
		Children: &Set[string]{},
		Parents:  &Set[string]{},
	}
}

type RDFItem struct {
	Subject string
	// Type is the RDF type of the item, e.g. "rdfs:Class", "rdfs:Property", "owl:Class", "owl:DatatypeProperty", etc.
	Type string
	// rdfs:domain
	Domain string
	// rdfs:comment
	Comment string
	// rdfs:range
	Range    string
	Children *Set[string]
	Parents  *Set[string]
}

type Property struct {
	Name  string
	Value string
}

type Set[T comparable] struct {
	values []T
}

func (s *Set[T]) Len() int {
	return len(s.values)
}

func (s *Set[T]) Values() []T {
	return s.values
}

func (s *Set[T]) Add(value T) {
	if s.Contains(value) {
		return
	}
	s.values = append(s.values, value)
}

func (s *Set[T]) Contains(value T) bool {
	for _, v := range s.values {
		if v == value {
			return true
		}
	}
	return false
}

func run() (err error) {
	f, err := os.Open("ies.rdf")
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	subjectToItem := make(map[string]*RDFItem)

	s := bufio.NewScanner(f)
	var lineNumber int
	for s.Scan() {
		lineNumber++
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "@") {
			continue
		}
		if !(strings.HasPrefix(line, "ies") || strings.HasPrefix(line, "<")) {
			return fmt.Errorf("error at line %d: %w", lineNumber, fmt.Errorf("line %q does not start with 'ies'", line))
		}

		parts := strings.SplitN(strings.TrimSuffix(line, " ."), " ", 3)
		if len(parts) != 3 {
			return fmt.Errorf("error at line %d: %w", lineNumber, fmt.Errorf("line %q does not have 3 parts", line))
		}
		subject, predicate, object := parts[0], parts[1], parts[2]

		if subject == "" {
			return fmt.Errorf("error at line %d: %w", lineNumber, fmt.Errorf("line %q does not have a subject", line))
		}
		if predicate == "" {
			return fmt.Errorf("error at line %d: %w", lineNumber, fmt.Errorf("line %q does not have a predicate", line))
		}
		if object == "" {
			return fmt.Errorf("error at line %d: %w", lineNumber, fmt.Errorf("line %q does not have an object", line))
		}

		item := getOrCreateItem(subject, subjectToItem)
		item.Subject = subject
		switch predicate {
		case "rdf:type":
			item.Type = object
		case "rdfs:comment":
			item.Comment = object
		case "ies:powertype":
			item.Parents.Add(object)
			parent := getOrCreateItem(object, subjectToItem)
			parent.Children.Add(subject)
		case "rdfs:subClassOf":
			item.Parents.Add(object)
			parent := getOrCreateItem(object, subjectToItem)
			parent.Children.Add(subject)
		case "rdfs:subPropertyOf":
			item.Parents.Add(object)
			parent := getOrCreateItem(object, subjectToItem)
			parent.Children.Add(subject)
		case "rdfs:domain":
			item.Domain = object
		case "rdfs:range":
			item.Range = object
		default:
			fmt.Printf("unknown predicate %q on line %d\n", predicate, lineNumber)
		}
	}

	// Now, let's get all of the items that don't have parents and walk the tree.
	var roots []*RDFItem
	for _, key := range slices.Sorted(maps.Keys(subjectToItem)) {
		item := subjectToItem[key]
		if item.Parents.Len() == 0 {
			roots = append(roots, item)
		}
	}

	for _, root := range roots {
		display(root, subjectToItem, 0)
	}

	return nil
}

func display(item *RDFItem, subjectToItem map[string]*RDFItem, indent int) {
	fmt.Printf("%s%s\n", strings.Repeat(" ", indent), item.Subject)
	sortedChildren := item.Children.Values()
	slices.Sort(sortedChildren)
	for _, child := range sortedChildren {
		display(subjectToItem[child], subjectToItem, indent+2)
	}
}

func getOrCreateItem(subject string, subjectToItem map[string]*RDFItem) *RDFItem {
	if _, ok := subjectToItem[subject]; !ok {
		subjectToItem[subject] = NewRDFItem(subject)
	}
	return subjectToItem[subject]
}

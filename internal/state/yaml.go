package state

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// YAMLDoc wraps a YAML document for structured edits. Operates on a generic
// map[string]interface{} mirroring how generators construct docker-compose
// and CI files.
type YAMLDoc struct {
	root map[string]interface{}
}

func NewYAMLDoc() *YAMLDoc {
	return &YAMLDoc{root: map[string]interface{}{}}
}

func (d *YAMLDoc) Load(data []byte) error {
	if len(data) == 0 {
		d.root = map[string]interface{}{}
		return nil
	}
	var m map[string]interface{}
	if err := yaml.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("yaml: unmarshal: %w", err)
	}
	if m == nil {
		m = map[string]interface{}{}
	}
	d.root = m
	return nil
}

func (d *YAMLDoc) Marshal() ([]byte, error) {
	out, err := yaml.Marshal(d.root)
	if err != nil {
		return nil, fmt.Errorf("yaml: marshal: %w", err)
	}
	return out, nil
}

func (d *YAMLDoc) Root() map[string]interface{} { return d.root }

// SetKey sets a top-level key.
func (d *YAMLDoc) SetKey(key string, value interface{}) {
	d.root[key] = value
}

// Append appends an entry to a list at the given top-level key. The slot is
// created as a slice if missing.
func (d *YAMLDoc) Append(key string, value interface{}) error {
	existing, ok := d.root[key]
	if !ok {
		d.root[key] = []interface{}{value}
		return nil
	}
	list, ok := existing.([]interface{})
	if !ok {
		return fmt.Errorf("yaml: key %q is not a list", key)
	}
	d.root[key] = append(list, value)
	return nil
}

// Merge deep-merges src into the document.
func (d *YAMLDoc) Merge(src map[string]interface{}) {
	mergeMap(d.root, src)
}

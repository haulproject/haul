// Package types implements the types and validation of types necessary for the application to function.
package types

import (
	"encoding/json"

	"github.com/cheynewallace/tabby"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Component struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Status string   `json:"status"`

	// A component's Target should point to a kit's or assembly's ObjectID
	Target primitive.ObjectID `json:"target"`
}

type ComponentWithID struct {
	ID primitive.ObjectID `json:"_id"`
	Component
}

type Assembly struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Status string   `json:"status"`

	// An assembly's Target should point to a kit's ObjectID
	Target primitive.ObjectID `json:"target"`
}

type AssemblyWithID struct {
	ID primitive.ObjectID `json:"_id"`
	Assembly
}

type Kit struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Status string   `json:"status"`
}

type KitWithID struct {
	ID primitive.ObjectID `json:"_id"`
	Kit
}

type Components struct {
	Components []Component `json:"components"`
}

type ComponentsWithID struct {
	ComponentsWithID []ComponentWithID `json:"components"`
}

type Assemblies struct {
	Assemblies []Assembly `json:"assemblies"`
}

type AssembliesWithID struct {
	AssembliesWithID []AssemblyWithID `json:"assemblies"`
}

type Kits struct {
	Kits []Kit `json:"kits"`
}

type KitsWithID struct {
	KitsWithID []KitWithID `json:"kits"`
}

type TabbyPrinter interface {
	TabbyPrint() error
}

func (c *Component) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader("name", "tags", "status", "target")

	tags, err := json.Marshal(c.Tags)
	if err != nil {
		return err
	}

	t.AddLine(c.Name, string(tags), c.Status, c.Target)

	t.Print()
	return nil
}

func (c *ComponentWithID) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader("id", "name", "tags", "status", "target")

	tags, err := json.Marshal(c.Tags)
	if err != nil {
		return err
	}

	t.AddLine(c.ID, c.Name, string(tags), c.Status, c.Target)

	t.Print()
	return nil
}

func (c *Components) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader("name", "tags", "status", "target")

	for _, component := range c.Components {
		tags, err := json.Marshal(component.Tags)
		if err != nil {
			return err
		}

		t.AddLine(component.Name, string(tags), component.Status, component.Target)
	}

	t.Print()
	return nil
}

func (c *ComponentsWithID) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader("id", "name", "tags", "status", "target")

	for _, component := range c.ComponentsWithID {
		tags, err := json.Marshal(component.Tags)
		if err != nil {
			return err
		}

		t.AddLine(component.ID, component.Name, string(tags), component.Status, component.Target)
	}

	t.Print()
	return nil
}

func (a *Assembly) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader("name", "tags", "status", "target")

	tags, err := json.Marshal(a.Tags)
	if err != nil {
		return err
	}

	t.AddLine(a.Name, string(tags), a.Status, a.Target)

	t.Print()
	return nil
}

func (a *AssemblyWithID) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader("id", "name", "tags", "status", "target")

	tags, err := json.Marshal(a.Tags)
	if err != nil {
		return err
	}

	t.AddLine(a.ID, a.Name, string(tags), a.Status, a.Target)

	t.Print()
	return nil
}

func (a *Assemblies) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader("name", "tags", "status", "target")

	for _, assembly := range a.Assemblies {
		tags, err := json.Marshal(assembly.Tags)
		if err != nil {
			return err
		}

		t.AddLine(assembly.Name, string(tags), assembly.Status, assembly.Target)
	}

	t.Print()
	return nil
}

func (a *AssembliesWithID) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader("id", "name", "tags", "status", "target")

	for _, assembly := range a.AssembliesWithID {
		tags, err := json.Marshal(assembly.Tags)
		if err != nil {
			return err
		}

		t.AddLine(assembly.ID, assembly.Name, string(tags), assembly.Status, assembly.Target)
	}

	t.Print()
	return nil
}

func (k *Kit) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader("name", "tags", "status")

	tags, err := json.Marshal(k.Tags)
	if err != nil {
		return err
	}

	t.AddLine(k.Name, string(tags), k.Status)

	t.Print()
	return nil
}

func (k *KitWithID) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader("id", "name", "tags", "status")

	tags, err := json.Marshal(k.Tags)
	if err != nil {
		return err
	}

	t.AddLine(k.ID, k.Name, string(tags), k.Status)

	t.Print()
	return nil
}

func (k *Kits) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader("name", "tags", "status")

	for _, kit := range k.Kits {
		tags, err := json.Marshal(kit.Tags)
		if err != nil {
			return err
		}

		t.AddLine(kit.Name, string(tags), kit.Status)
	}

	t.Print()
	return nil
}

func (k *KitsWithID) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader("id", "name", "tags", "status")

	for _, kit := range k.KitsWithID {
		tags, err := json.Marshal(kit.Tags)
		if err != nil {
			return err
		}

		t.AddLine(kit.ID, kit.Name, string(tags), kit.Status)
	}

	t.Print()
	return nil
}

/*
Taking a bson.D and a reference interface{}, returns a bson.D that only contains fields whose keys are also valid keys in the reference interface{}.

If an error is encountered, it is returned, and the returned bson.D will be nil.

If no error is encountered, the validated bson.D is returned, and the error will be nil.
*/
func ValidateFields(document bson.D, reference interface{}) (bson.D, error) {
	keys, err := GetFields(reference)
	if err != nil {
		return nil, err
	}

	var validFields bson.D

	for _, value := range document {
		valid := false
		for _, key := range keys {
			if value.Key == key {
				valid = true
			}
		}
		if valid {
			validFields = append(validFields, value)
		}
	}

	return validFields, nil
}

// GetFields returns an unordered list of fields in a reference interface{}
func GetFields(reference interface{}) ([]string, error) {
	var fields []string
	fields_bytes, err := json.Marshal(reference)
	if err != nil {
		return fields, err
	}

	m := make(map[string]interface{})

	if err = json.Unmarshal(fields_bytes, &m); err != nil {
		return fields, err
	}

	for field := range m {
		fields = append(fields, field)
	}

	return fields, nil
}

func GetFieldsOrdered(reference interface{}) ([]string, error) {
	var fields []string

	fields_bytes, err := json.Marshal(reference)
	if err != nil {
		return fields, err
	}

	var fields_map map[string]interface{}

	err = json.Unmarshal(fields_bytes, &fields_map)
	if err != nil {
		return fields, err
	}

	for field := range fields_map {
		fields = append(fields, field)
	}

	return fields, nil
}

type InsertResult struct {
	Message     string        `json:"message"`
	InsertedIDs []interface{} `json:"inserted_ids,omitempty"`
}

func (r InsertResult) TabbyPrint() error {
	t := tabby.New()

	t.AddHeader(r.Message)

	for _, id := range r.InsertedIDs {
		t.AddLine(id)
	}

	t.Print()
	return nil
}

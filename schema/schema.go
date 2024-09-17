/*
Copyright © 2024 Dustin Dortch
*/
package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Attribute struct {
	Type            any    `json:"type"`
	Description     string `json:"description"`
	DescriptionKind string `json:"description_kind,omitempty"`
	Required        bool   `json:"required"`
	Optional        bool   `json:"optional"`
	Computed        bool   `json:"computed"`
	Sensitive       bool   `json:"sensitive"`
}

func (a Attribute) TypeString() string {
	return fmt.Sprintf("%s", a.Type)
}

type Block struct {
	Attributes      map[string]Attribute `json:"attributes"`
	Blocks          map[string]Block     `json:"block_types,omitempty"`
	DescriptionKind string               `json:"description_kind,omitempty"`
	MaxItems        int                  `json:"max_items,omitempty"`
	MinItems        int                  `json:"min_items,omitempty"`
}

type Resource struct {
	Version int   `json:"version"`
	Block   Block `json:"block"`
}

type Provider struct {
	Resources map[string]Resource `json:"resource_schemas"`
}

type Schema struct {
	Version   string              `json:"format_version"`
	Providers map[string]Provider `json:"provider_schemas"`
}

func ReadSchema() Schema {
	schemaJson := exec.Command("terraform", "providers", "schema", "-json")
	stdout, err := schemaJson.Output()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var schema Schema
	err = json.Unmarshal([]byte(stdout), &schema)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %s\n", err)
		os.Exit(2)
	}

	return schema
}

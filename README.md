# JSONStream ![](https://github.com/vcraescu/go-jsonstream/actions/workflows/go.yml/badge.svg)

The `jsonstream` package provides a convenient way to stream and unmarshal JSON data from an `io.Reader` source. It allows you to decode JSON objects sequentially and asynchronously, providing each decoded entry along with any encountered errors.

## Overview

The package is designed to address scenarios where you need to process large JSON datasets in a streaming manner, reading JSON objects from an input source and decoding them into Go values. This can be particularly useful when dealing with JSON data that doesn't fit entirely in memory or when you want to process data in smaller batches.

## Usage

### Installation

To use the `jsonstream` package, you need to import it into your project:

```go
import "github.com/your-username/jsonstream"
```

### Example

Here's a basic example of how to use the `jsonstream` package to stream and unmarshal JSON data from an `io.Reader`:

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/your-username/jsonstream"
)

func main() {
	data := `
		[
			{"name": "Alice", "age": 25},
			{"name": "Bob", "age": 30},
			{"name": "Charlie", "age": 28}
		]
	`

	reader := strings.NewReader(data)
	ctx := context.Background()

	entryChan, err := jsonstream.Unmarshal(ctx, reader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for entry := range entryChan {
		if entry.Err != nil {
			fmt.Println("Error:", entry.Err)
			continue
		}

		var person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		person = entry.Value
		fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)
	}
}
```

### Options

The `Unmarshal` function supports optional configuration through the use of `Option` arguments. The available options are:

- `BatchSize(int)`: Sets the number of entries to process in each batch (default is 100).
- `StartFrom(int)`: Specifies the index of the first JSON object to start processing (default is 1).

You can use these options to control the behavior of the JSON decoding process according to your needs.

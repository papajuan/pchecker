Go Single-Pass Profanity Detector
----------------------------------
A high-performance, memory-efficient profanity detector in Go that scans text in a single loop. It uses a trie-based approach for fast detection and supports customizable lists of profanities and false positives.
__________________________________
Features
----------------------------------
- Single-pass scanning for minimal allocations and high speed

- Trie-based matching for efficient prefix and full-word detection

- Support for false positives to avoid over-censoring

- Customizable replacement function for censoring words

- Written in pure Go with memory reuse for optimal performance

Installation
----------------------------------
```bash
go get github.com/papajuan/pchecker
```

Usage

```go
package main

import (
	"fmt"

	"github.com/papajuan/pchecker"
)

func main() {
	input := "Is that poop, huh?"

	// Create a default detector
	pd := pchecker.NewDefaultProfanityDetector()

	// Define replacement function
	f := func(match []rune) []rune {
		return []rune{'*', '*', '*'}
	}

	// Censor input
	censored := pd.Censor(input, f)
	fmt.Println(censored) // Output: Is that ***, huh?
}
```


Expected performance:

- ~5–6 µs per operation for average sentences

- Minimal memory allocations (~8 allocs per operation)

Contributing

Contributions are welcome! Feel free to open issues, suggest improvements, or submit pull requests.

License

MIT License © [papajuan]
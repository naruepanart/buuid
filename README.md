# buuid - Simple Random Generator Package

![Go](https://img.shields.io/badge/Go-1.16+-blue.svg)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A lightweight Go package for generating various types of random values with simple APIs.

## Features

- Generate random strings with custom character sets
- Create random integers within a range
- Generate random floating-point numbers
- Produce unique IDs with timestamp components
- Thread-safe operations

## Installation

```bash
go get github.com/naruepanart/buuid
```

## Usage

### Random Strings

```go
import "github.com/naruepanart/buuid"

// Available character set flags:
// buuid.R_NUM   - Numbers only (0-9)
// buuid.R_UPPER - Uppercase letters only (A-Z)
// buuid.R_LOWER - Lowercase letters only (a-z)
// buuid.R_All   - All characters (0-9, A-Z, a-z)

// Generate 10-character random string with numbers and lowercase letters
randomStr := buuid.String(buuid.R_NUM|buuid.R_LOWER, 10)

// Generate 6-character random string with all character types (default length)
randomStr := buuid.String(buuid.R_All)
```

### Random Numbers

```go
// Random integer between 0-100 (default)
num := buuid.Int()

// Random integer between 0-50
num := buuid.Int(50)

// Random integer between 10-20
num := buuid.Int(10, 20)
```

### Random Floating-Point Numbers

```go
// Random float with 2 decimal places, between 0.0-100.0
f := buuid.Float64(2)

// Random float with 3 decimal places, between 0.0-50.0
f := buuid.Float64(3, 50)

// Random float with 4 decimal places, between 10.0-20.0
f := buuid.Float64(4, 10, 20)
```

### Unique IDs

```go
// Numeric ID combining timestamp and random component
id := buuid.NewID() // e.g., 1651234567890123456

// Hexadecimal string version of NewID()
hexID := buuid.NewStringID() // e.g., "16f3a5b7c8d9e0f1"

// Formatted timestamp ID with random suffix
seriesID := buuid.NewSeriesID() // e.g., "2023052312453000000123456"
```

## Performance

The package uses crypto/rand for secure random generation by default, with fallback to time-based randomness if crypto/rand fails.

## License

MIT License. See [LICENSE](LICENSE) for details.
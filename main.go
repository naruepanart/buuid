package buuid

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"strconv"
	"sync"
	"time"
)

// nolint
const (
	R_NUM   = 1 // only number
	R_UPPER = 2 // only capital letters
	R_LOWER = 4 // only lowercase letters
	R_All   = 7 // numbers, upper and lower case letters
)

var (
	// Pre-calculated character sets
	numChars    = []byte("0123456789")
	upperChars  = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lowerChars  = []byte("abcdefghijklmnopqrstuvwxyz")
	allChars    = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	charSets    = [][]byte{nil, numChars, upperChars, nil, lowerChars, nil, nil, allChars}
	defaultRand = &lockedRandSource{}
)

type lockedRandSource struct {
	mu sync.Mutex
}

func (r *lockedRandSource) Int63() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return time.Now().UnixNano()
	}
	return int64(binary.BigEndian.Uint64(b[:]) & (1<<63 - 1))
}

// String generates random strings of any length of multiple types, default length is 6 if size is empty
// example: String(R_ALL), String(R_ALL, 16), String(R_NUM|R_LOWER, 16)
func String(kind int, size ...int) string {
	return string(Bytes(kind, size...))
}

// Bytes generates random strings of any length of multiple types, default length is 6 if bytesLen is empty
// example: Bytes(R_ALL), Bytes(R_ALL, 16), Bytes(R_NUM|R_LOWER, 16)
func Bytes(kind int, bytesLen ...int) []byte {
	if kind > 7 || kind < 1 {
		kind = R_All
	}

	length := 6 // default length 6
	if len(bytesLen) > 0 && bytesLen[0] > 0 {
		length = bytesLen[0]
	}

	chars := charSets[kind]
	if chars == nil {
		// Handle combined character sets
		combined := make([]byte, 0, 62)
		if kind&R_NUM != 0 {
			combined = append(combined, numChars...)
		}
		if kind&R_UPPER != 0 {
			combined = append(combined, upperChars...)
		}
		if kind&R_LOWER != 0 {
			combined = append(combined, lowerChars...)
		}
		chars = combined
	}

	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			n = big.NewInt(defaultRand.Int63() % int64(len(chars)))
		}
		result[i] = chars[n.Int64()]
	}

	return result
}

// Int generates random numbers of specified range size,
// compatible with Int(), Int(max), Int(min, max), Int(max, min) 4 ways, min<=random number<=max
func Int(rangeSize ...int) int {
	var min, max int

	switch len(rangeSize) {
	case 0:
		min, max = 0, 100 // default 0~100
	case 1:
		min, max = 0, rangeSize[0]
	default:
		if rangeSize[0] > rangeSize[1] {
			min, max = rangeSize[1], rangeSize[0]
		} else {
			min, max = rangeSize[0], rangeSize[1]
		}
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return min + int(defaultRand.Int63()%int64(max-min+1))
	}
	return min + int(n.Int64())
}

// Float64 generates a random floating point number of the specified range size,
// Four types of passing references are supported, example: Float64(dpLength), Float64(dpLength, max),
// Float64(dpLength, min, max), Float64(dpLength, max, min), min<=random numbers<=max
func Float64(dpLength int, rangeSize ...int) float64 {
	var min, max int

	switch len(rangeSize) {
	case 0:
		min, max = 0, 100 // default 0~100
	case 1:
		min, max = 0, rangeSize[0]
	default:
		if rangeSize[0] > rangeSize[1] {
			min, max = rangeSize[1], rangeSize[0]
		} else {
			min, max = rangeSize[0], rangeSize[1]
		}
	}

	// Generate decimal part
	dp := 0.0
	if dpLength > 0 {
		dpmax := big.NewInt(10)
		dpmax.Exp(dpmax, big.NewInt(int64(dpLength)), nil)
		n, err := rand.Int(rand.Reader, dpmax)
		if err != nil {
			n = big.NewInt(defaultRand.Int63() % dpmax.Int64())
		}
		dp = float64(n.Int64()) / float64(dpmax.Int64())
	}

	// Generate integer part
	intPart, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		intPart = big.NewInt(defaultRand.Int63() % int64(max-min))
	}

	return float64(min) + float64(intPart.Int64()) + dp
}

// NewID generates a milliseconds+random number ID.
func NewID() int64 {
	var buf [8]byte
	now := time.Now().UnixMilli() * 1000000

	_, err := rand.Read(buf[:])
	if err != nil {
		return now + defaultRand.Int63()%1000000
	}

	return now + int64(binary.LittleEndian.Uint64(buf[:])%1000000)
}

// NewStringID generates a string ID, the hexadecimal form of NewID(), total 16 bytes.
func NewStringID() string {
	return strconv.FormatInt(NewID(), 16)
}

// NewSeriesID generates a datetime+random string ID,
// datetime is microsecond precision, 20 bytes, random is 6 bytes, total 26 bytes.
// example: 20060102150405000000123456
func NewSeriesID() string {
	var buf [26]byte
	t := time.Now()

	// Format datetime with microsecond precision (14 bytes)
	copy(buf[:14], t.Format("20060102150405"))

	// Add microseconds (6 bytes)
	micro := t.Nanosecond() / 1000
	buf[14] = '0' + byte(micro/100000%10)
	buf[15] = '0' + byte(micro/10000%10)
	buf[16] = '0' + byte(micro/1000%10)
	buf[17] = '0' + byte(micro/100%10)
	buf[18] = '0' + byte(micro/10%10)
	buf[19] = '0' + byte(micro%10)

	// Generate a 6-digit random number
	random := Int(0, 999999)
	for i := 20; i < 26; i++ {
		buf[i] = '0' + byte(random%10)
		random /= 10
	}

	return string(buf[:])
}

//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

const (
	allowedDuration = 10 * time.Second
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	mx        sync.Mutex
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}
	u.mx.Lock()
	defer u.mx.Unlock()

	if u.TimeUsed < int64(allowedDuration) {
		done := make(chan struct{})
		var startTime time.Time

		go func() {
			startTime = time.Now()
			process()
			close(done)
		}()

		select {
		case <-done:
			u.TimeUsed += int64(time.Since(startTime))
			return true
		case <-time.After(allowedDuration - time.Duration(u.TimeUsed)):
			u.TimeUsed += int64(time.Since(startTime))
			return false
		}
	}
	return false
}

func main() {
	RunMockServer()
}

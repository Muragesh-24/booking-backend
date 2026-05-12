package scripts

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)


type Visitor struct {
	LastSeen time.Time
	Requests int
}

var visitors = make(map[string]*Visitor)
var mu sync.Mutex


func LimitPerIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists || time.Since(v.LastSeen) > time.Minute {
	
			visitors[ip] = &Visitor{LastSeen: time.Now(), Requests: 1}
		} else {
			v.Requests++
			v.LastSeen = time.Now()
			if v.Requests > 150 { 
				mu.Unlock()
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Too many requests, please try again later",
				})
				c.Abort()
				return
			}
		}
		mu.Unlock()

		c.Next()
	}
}

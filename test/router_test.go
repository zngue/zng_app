package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestRouter(t *testing.T) {
	engine := gin.Default()
	fmt.Println(engine)

}

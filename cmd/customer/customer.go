package main

import (
	"github.com/karincake/apem"

	"github.com/munaja/blog-be-in-go/internal/handler/customer"
)

func main() {
	apem.Run("skgobo/customer", customer.SetRoutes())
}

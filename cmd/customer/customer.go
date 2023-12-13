package main

import (
	"github.com/karincake/apem"

	"github.com/munaja/blog-practice-be-using-go/internal/handler/customer"
)

func main() {
	apem.Run("skgobo/customer", customer.SetRoutes())
}

# NativeBPM Go Client SDK

> Cloud-native BPMN 2.0 / DMN 1.3 client and Fluent workflow builder for Go.

This branch (`go`) is the dedicated home for the NativeBPM Go Client SDK (`module gitlab.com/nativebpm/sdk/go`).

---

## 📦 Installation

```bash
go get gitlab.com/nativebpm/sdk/go@latest
```

---

## 🚀 Quick Start

```go
package main

import (
	"context"
	"fmt"
	"gitlab.com/nativebpm/sdk/go"
)

type OrderInput struct {
	OrderID string  `json:"orderId" validate:"required"`
	Amount  float64 `json:"amount" validate:"gt=0"`
	Tier    string  `json:"tier" validate:"oneof=standard vip"`
}

func main() {
	client := nativebpm.NewClient("http://localhost:8080", "secret-token")

	wf := nativebpm.NewWorkflow("order-process", "Order Fulfillment").
		Variables(OrderInput{}).
		Start("start").
		ExclusiveGateway("check_tier", "Check Tier").
		When("tier == 'vip'").
		Then(func(b *nativebpm.Branch) {
			b.Service("vip_task", "VIP Processing", "vip_topic")
		}).
		Otherwise(func(b *nativebpm.Branch) {
			b.Service("std_task", "Standard Processing", "std_topic")
		}).
		End("end", "End")

	resp, err := client.Deploy(wf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deployed workflow: %s\n", resp.ID)
}
```

---

## 🧪 Testing

```bash
make test
```

# NativeBPM Python Client SDK

> Python client, Pydantic v2 data models, and Fluent workflow builder for NativeBPM 2.0 / DMN 1.3 execution engine.

This branch (`python`) is the dedicated home for the NativeBPM Python Client SDK.

---

## 📦 Installation

```bash
pip install git+https://gitlab.com/nativebpm/sdk.git@python#subdirectory=python
```

---

## 🚀 Quick Start (Pydantic v2 + Fluent Builder)

```python
from nativebpm import Client, Workflow
from pydantic import BaseModel, Field
from typing import Literal

# 1. Declare process input contract using Pydantic v2
class OrderInput(BaseModel):
    order_id: str = Field(min_length=1)
    amount: float = Field(gt=0)
    tier: Literal["standard", "vip"]

# 2. Build the workflow graph with embedded JSON Schema and branching
workflow = Workflow("order-process", "Order Fulfillment")\
    .variables(OrderInput)\
    .start("start")\
    .exclusive_gateway("check_tier", "Check Tier")\
    .when("tier == 'vip'").then("vip_task").service_task("vip_task", "VIP Processing", "vip_topic")\
    .otherwise("std_task").service_task("std_task", "Standard Processing", "std_topic")\
    .end("end", "End")

client = Client("http://localhost:8080", "secret-token")
client.deploy(workflow)

# 3. Fail-fast validation before network call
raw_input = {"order_id": "ORD-123", "amount": 250.0, "tier": "vip"}
validated = OrderInput(**raw_input)

instance = client.instances().start("order-process")\
    .with_variables(validated.model_dump())\
    .send()

print(f"Started instance: {instance.id}")
```

---

## 🧪 Testing

```bash
make test
```

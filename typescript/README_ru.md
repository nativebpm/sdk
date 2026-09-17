[← Назад к корню платформы](../../README_ru.md)

# NativeBPM TypeScript SDK

Официальный TypeScript Client SDK для облачного BPMN 2.0 / DMN 1.3 движка NativeBPM.

Разработан для **Node.js 22+ (с поддержкой прямого исполнения TypeScript без шага сборки через `--experimental-strip-types`)**, Edge CDN сред (Cloudflare Workers / Pages) и современных браузеров.

---

## Ключевые возможности

- **Native Zod 4 «из коробки»**: Реэкспорт `z`, `ZodError` и `ZodType`. Описание схем форм и графа процесса в одном файле (`*.flow.ts`).
- **Автоматическая компиляция в JSON Schema Draft 2020-12**: Нативная компиляция Zod-схем в Draft 2020-12 для Server-Driven UI (`<nativebpm-trigger>`, Vue, React, Capacitor).
- **Трехуровневая архитектура**:
  1. *Уровень 1 (SSOT)*: Спецификация OpenAPI 3.0 (`openapi.yaml`).
  2. *Уровень 2 (Generated Transport)*: Строго типизированный низкоуровневый клиент.
  3. *Уровень 3 (Fluent API Façade & Worker)*: Высокоуровневый фасад с цепочками вызовов для деплоя, запуска процессов, управления задачами и воркерами.
- **Clean Code (Zero Parameter Properties)**: 100% совместимость со strip-типами Node 22+ (`node script.ts`).
- **Экспорт в стандартный OMG BPMN 2.0 XML**: Метод `workflow.toBPMN()`.
- **Таблицы решений DMN 1.3**: Поддержка `businessRuleTask` и локального вычисления правил.

---

## Установка

```bash
npm install @nativebpm/sdk
```

---

## 1. Single-File Workflow-as-Code и Schema-as-Code на Zod

```typescript
import { z, WorkflowBuilder } from "@nativebpm/sdk";

// 1. Описание схемы формы через встроенный Zod 4
export const OrderFormSchema = z.object({
  order_id: z.string().min(1, "ID заказа обязателен"),
  amount: z.number().positive("Сумма должна быть положительной"),
  promo_code: z.string().optional(),
});

// 2. Описание процесса с использованием Fluent API
export const orderWorkflow = new WorkflowBuilder("order_flow", "Order Process")
  .start("start")
  .userTask("fill_order", "Ввод заказа", {
    form: OrderFormSchema,
    formId: "order_form_v1",
    candidateGroups: "sales",
  })
  .exclusiveGateway("check_amount", "Проверка суммы")
    .when("amount > 100000").then("vip_review").userTask("vip_review", "VIP Проверка", { candidateGroups: "vip_managers" })
    .else("auto_approve").serviceTask("auto_approve", "Автоодобрение", "auto_approver")
  .end("end", "Процесс завершен");

// 3. Извлечение скомпилированных JSON Schema Draft 2020-12 для Server-Driven UI (BDUI)
const forms = orderWorkflow.extractForms();
console.log(forms.order_form_v1);
// Результат: { $schema: "https://json-schema.org/draft/2020-12/schema", type: "object", ... }

// 4. Экспорт стандартного BPMN 2.0 XML
const bpmnXml = orderWorkflow.toBPMN();
```

---

## 2. Fluent Client API

Клиент NativeBPM предоставляет модульный цепочечный интерфейс:

```typescript
import { Client } from "@nativebpm/sdk";
import { orderWorkflow } from "./order.flow.js";

const client = new Client("http://localhost:8080", "your-api-token");

// 1. Деплой дефиниции процесса
const definition = await client.definitions().deploy()
  .withWorkflow(orderWorkflow)
  .send();

// 2. Запуск инстанса процесса
const instance = await client.instances()
  .start("order_flow")
  .withBusinessKey("ORD-2026-001")
  .withVariables({ amount: 1500, promo_code: "DISCOUNT10" })
  .send();

// 3. Запрос и взятие пользовательских задач
const tasks = await client.tasks().list()
  .withAssignee("sales_manager")
  .withStatus("CREATED")
  .send();

if (tasks.length > 0) {
  const task = tasks[0];
  await client.tasks().claim(task.id).withAssignee("sales_manager").send();

  // 4. Завершение задачи с передачей выходных переменных
  await client.tasks().complete(task.id)
    .withVariables({ approved: true, comment: "Одобрено" })
    .send();
}
```

---

## 3. Воркеры фоновых сервисных задач (Worker)

Класс `Worker` с зеркальным Go-синтаксисом для выполнения сервисных задач:

```typescript
import { Worker } from "@nativebpm/sdk";

const worker = new Worker("http://localhost:8080", "your-api-token")
  .withWorkerId("billing-worker-01")
  .withMaxConcurrency(20)
  .withTopic("payment_gateway", async (task) => {
    console.log(`Обработка задачи ${task.id} инстанса ${task.instanceId}`);
    console.log("Входные переменные:", task.variables);

    if (task.variables.amount > 500000) {
      // Выброс исключения автоматически переводит инстанс в статус Incident
      throw new Error("Превышен кредитный лимит");
    }

    // Возврат переменных для продвижения токена в движке
    return {
      payment_status: "PAID",
      tx_id: `tx_${Date.now()}`,
    };
  });

// Запуск цикла обработки
await worker.start();
```

---

## Запуск тестов

Пакет включает сквозной тестовый сьют (TDD), покрывающий 6 фундаментальных паттернов бизнес-логики:
1. Линейный процесс с валидацией Zod (UserTask $\rightarrow$ ServiceTask).
2. Таблицы решений DMN (BusinessRuleTask).
3. Условное ветвление (ExclusiveGateway When/Then/Else).
4. Сервисные воркеры и инциденты (Service Tasks & Workers).
5. Таймеры и SLA-эскалация (`PT15M`).
6. Экспорт схем под Server-Driven UI (BDUI).

```bash
# Сборка и запуск всех тестов
npm test
```

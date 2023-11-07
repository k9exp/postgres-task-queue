# postgres-task-queue

Using postgres to make a task queue

Blog: https://kunalsin9h.com/blog/potgres-task-queue

Run the following commands to get started:

```bash
export DATABASE_URL=...

go run main.go
```

### API Docs

POST `/producer` - Create a new task

With body

```json
{
  "text": "The text to be printed",
  "time": 2
}
```

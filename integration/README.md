# Integration Reference

Reference implementation for embedding Skill-Plus Engine into a host system.

## V-Box Integration (Reference)

The engine integrates into V-Box's posting pipeline as follows:

```
POST /v1/posts (A-face content)
    │
    ▼
backend (Go)
    │ ① moderation
    │ ② media verification
    │ ③ persist post (A-face), state = processing
    │ ④ trigger Skill-Plus pipeline (async)
    │ ⑤ return to client (optimistic published)
    ▼
skillplus-engine
    │ ⑥ read post A-face content
    │ ⑦ dispatch matching Skills
    │ ⑧ concurrent execution in sandboxes
    │ ⑨ aggregate → B-face JSON
    ▼
backend
    │ ⑩ persist B-face
    │ ⑪ post.state → published
    │ ⑫ emit feed_distribution event
```

See `integration_example.go` for a minimal integration skeleton.

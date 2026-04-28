# Food Exploration Example

This example shows how OpenMelon uses a compiled Skill-Plus package during a food exploration workflow.

User intent:

```text
下班后在老小区楼下吃一碗牛肉面，想发一条真实的探店帖。
```

Workflow:

1. Create project context for a social food account.
2. Select the food exploration workflow.
3. Enter `visual_concretization` stage.
4. Compile `food-street-realism.skillplus` with:
   - target: `openmelon`
   - model profile: `gpt-image-family`
   - locale: `zh-CN`
   - runtime var: `realism_level=high`
5. Produce an image prompt artifact.
6. Route the prompt artifact to the configured image generation model.
7. Produce an image artifact.
8. Attach labels and provenance to both the prompt and image.
9. Return review checklist and failure modes.

Run without final image generation:

```bash
go run ./cmd/openmelon \
  --example examples/food-exploration/beef-noodles.json \
  --compiler ../skillplus/compiler/reference \
  --config config/openmelon.example.json
```

Run with final image generation:

```bash
go run ./cmd/openmelon \
  --example examples/food-exploration/beef-noodles.json \
  --compiler ../skillplus/compiler/reference \
  --config config/openmelon.example.json \
  --generate \
  --artifact-dir examples/food-exploration/artifacts
```

The important part is not the prompt alone. The important part is the trace:

```text
source skill package -> compiled skill -> prompt artifact -> generation model -> image artifact -> labels -> provenance -> review
```

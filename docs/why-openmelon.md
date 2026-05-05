# Why OpenMelon

OpenMelon exists because content creation needs more than one-shot prompting.

A serious content workflow has to keep context across many turns: which project the work belongs to, which character must keep looking like the same person across posts, which scenes are reusable, which generations were good enough to ship vs which were drafts. None of that fits into a chat window with copy-paste.

OpenMelon is a terminal agent for content production. Each project is a directory on disk with persistent character / reference libraries, sessions that record every turn, and a tool-using LLM that can pull from them while drafting. Generated images are anchored to the same character portraits run after run. The session log makes the production process inspectable and resumable.

It runs locally — no SaaS account, no usage tracking — against any LLM + image model you have an API key for (OpenRouter / OpenAI / Anthropic, plus most things they route to). Skills (the prompt + output-schema bundles that turn intent into a generation prompt) live in [skillplus](https://github.com/eight-acres-lab/skillplus), versioned and shareable.

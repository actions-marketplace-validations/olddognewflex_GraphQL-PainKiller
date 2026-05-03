# GraphQL Painkiller — V1

## TL;DR

GraphQL Painkiller is a zero-runtime static analysis tool that detects risky GraphQL query patterns in pull requests and flags them before they are merged.

It works as:
- A CLI tool
- A GitHub Action that leaves PR review comments

---

## Core Problem

GraphQL performance issues are:
- Hard to detect before runtime
- Often caused by query shape (not code)
- Invisible during code review

AI tools (Codex, GPT) can suggest improvements, but they:
- Do not analyze actual query structure deterministically
- Do not understand runtime fan-out risks
- Cannot enforce standards in CI

---

## Core Value Proposition

> Catch expensive GraphQL queries before they hit production.

---

## Product Shape

### Mode 1 — CLI

```bash
gql-painkiller analyze ./queries
```

### Mode 2 - GitHub Action

```yaml
- uses: olddognewflex/graphql-painkiller-action@v1
```

## V1 Capabilities

### Input Sources

- .graphql, .gql files
- Tagged template literals:

```typescript
gql`
  query GetPosts {
    posts {
      comments {
        id
      }
    }
  }`
```

- Comment based template
  
```typescript
const query = /* GraphQL */ `
  query GetPosts {
    posts {
      comments {
        id
      }
    }
  }
`;
```

### Output

#### CLI
```
GraphQL Painkiller Report

Operation: GetPosts
Risk Score: 8/10 — High

Findings:
- Potential N+1 path: posts.comments
- Missing pagination on posts
- Large selection set under posts
```

#### GitHub PR Comment
```
⚠️ Potential N+1 risk

posts → comments may trigger resolver fan-out.

Suggestions:
- Add pagination
- Batch resolver
- Avoid nested selection if unnecessary
```

## Rule Engine (V1)

### 1. Max Depth

Detect deeply nested queries.

**Config**:
```json
{
  "maxDepth": 5
}
```
### 2. Collection Detection

Identify likely list fields using:

- naming (items, nodes, plural fields)
- schema (if provided)

### 3. Missing Pagination

Flag collection-like fields without:

- first
- last
- limit
- take
- pageSize
- after / before  

### 4. Nested Collection Risk (N+1 Heuristic)

Detect patterns like:
```gql
posts {
  comments {
    author {
      name
    }
  }
}
```

Flag:

	Potential resolver fan-out (N+1 risk)

### 5. Large Selection Set

  Too many fields under a collection:

```json
{
  "maxCollectionSelectionFields": 8
}
```

### 6. Expensive Field Patterns
Configurable
```json
[
  "comments",
  "history",
  "events",
  "logs",
  "charges",
  "payments",
  "inspections",
  "accounts"
]
```

### 7. Known Resolver Risk (🔥 Differentiator)
Team-defined risk map:

```json
{
  "knownResolvers": {
    "posts.comments": {
      "risk": "high",
      "reason": "Resolver runs per post"
    }
  }
}
```

Enables **deterministic warnings** instead of heuristics.

#### Risk Scoring

Each rule contributes to a score:
```
riskScore = sum(ruleImpacts)
```

Severity:

- 0–3 → Low
- 4–6 → Medium
- 7–8 → High
- 9–10 → Critical

### Config File
```json
{
  "rules": {
    "maxDepth": 5,
    "maxCollectionSelectionFields": 8,
    "requirePagination": true,
    "failOnSeverity": "high"
  },
  "paginationArgs": ["first", "last", "limit", "take", "pageSize"],
  "collectionFieldPatterns": ["items", "nodes", "edges"],
  "expensiveFieldPatterns": [
    "comments",
    "history",
    "events",
    "logs",
    "charges",
    "payments",
    "inspections",
    "accounts"
  ],
  "knownResolvers": {}
}
```

## CLI Commands
```
gql-painkiller init
gql-painkiller analyze ./queries
gql-painkiller analyze . --changed-only
gql-painkiller analyze . --json
gql-painkiller analyze . --fail-on high
```

### GitHub Action
```yaml
name: GraphQL Painkiller

on:
  pull_request:
    paths:
      - "**/*.graphql"
      - "**/*.ts"
      - "**/*.tsx"

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: olddognewflex/graphql-painkiller-action@v1
```

## Architecture (V1)
```
src/
├── cli.ts
├── action.ts
├── extractors/
├── analyzer/
├── rules/
├── reporters/
├── github/
└── config/
```

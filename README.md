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


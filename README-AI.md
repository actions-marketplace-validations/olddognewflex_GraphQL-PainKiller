# GraphQL Painkiller — V1 Design Notes

## Product Definition

GraphQL Painkiller is a zero-runtime static analysis tool that reviews GraphQL operations before merge.

It works as:

1. A local CLI
2. A future GitHub Action / PR reviewer

## Core Promise

> Catch expensive GraphQL queries before they hit production.

## V1 Scope

Analyze:

- `.graphql`
- `.gql`
- `gql` tagged template literals
- `/* GraphQL */` template literals

Detect:

- deep queries
- missing pagination
- potential N+1 resolver fan-out
- large selection sets under collections
- expensive field names
- known resolver risks from config

## Why Go

Go is a strong fit for this product because:

- single binary distribution
- fast startup in CI
- no Node install required for consumers
- clean path to GitHub Action binary usage
- simple CLI ergonomics

## Non-Goals

V1 does **not** include:

- runtime tracing
- SaaS dashboard
- billing
- GitHub App auth
- AI-generated fixes
- precise latency prediction

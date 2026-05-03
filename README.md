# GraphQL Painkiller

Zero-runtime static analysis for risky GraphQL query patterns.

V1 is a Go CLI that scans:

- `.graphql`
- `.gql`
- `gql` tagged template literals
- `/* GraphQL */` template literals

It flags likely:

- deep queries
- missing pagination
- nested collection / N+1 risks
- expensive field names
- large selection sets under collection-like fields
- team-defined known resolver risks

## Build

```bash
go mod tidy
go build -o gql-painkiller ./cmd/gql-painkiller
```

## Run

```bash
./gql-painkiller analyze ./examples
```

JSON output:

```bash
./gql-painkiller analyze ./examples --json
```

Fail CI on high findings:

```bash
./gql-painkiller analyze ./examples --fail-on high
```

Create default config:

```bash
./gql-painkiller init
```

## Philosophy

GraphQL Painkiller does not pretend to know runtime truth from static analysis.

It says things like:

- likely
- potential
- may cause

Unless the path is configured in `knownResolvers`, which represents team-owned knowledge.

That distinction matters. Otherwise we’re just inventing performance astrology with a compiler badge.

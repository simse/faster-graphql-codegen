[![Build and Release](https://github.com/simse/faster-graphql-codegen/actions/workflows/release.yml/badge.svg)](https://github.com/simse/faster-graphql-codegen/actions/workflows/release.yml)
# 🚀 faster-graphql-codegen 🚀
A reimplementation of graphql-codegen designed to be (much) faster than the original. It aims to be drop-in compatible where possible.

[Docs](https://faster-graphql-codegen.simse.io/) | [Quick Start](https://faster-graphql-codegen.simse.io/quick-start)

## Current functionality
- [x] Parse single graphql schema from file
- [ ] Parse and merge multiple graphql schemas
- [ ] Load schema from URL via introspection
- [x] Convert schema to types
  - [x] Enums
  - [x] Comments
  - [x] InputObjects
  - [x] Objects
  - [x] Implements
  - [x] Custom Scalars
  - [ ] Lots of other things
- [x] Load .yaml config
- [ ] Load .js/.ts config
- [ ] Extract and generate types for queries
- [ ] Extract and generate types for mutations
- [x] Monorepo support

## Yo!
This is still just an experiment.
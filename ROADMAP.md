# SpiceDB Go Library Roadmap

This is the public roadmap for the SpiceDB Go Library. It outlines planned features and improvements.

## Vision

To provide the most complete, well-documented, and user-friendly Go client library for SpiceDB, enabling developers to build authorization into their applications with confidence and ease.

## Current Status: Beta

The library is in active development and production usage. API stability is not yet guaranteed, but we aim to minimize breaking changes.

## Prioritized Roadmap

### Phase 1: Core Operations (In Progress)

- [x] CheckPermission operation
- [x] Builder pattern support
- [x] Caching abstraction and in-memory implementation
- [x] Error handling and typed errors
- [x] Full documentation and examples
- [ ] ReadRelationships
- [ ] WriteRelationships
- [ ] DeleteRelationships
- [ ] LookupResources
- [ ] LookupSubjects

### Phase 2: Advanced Features

- [ ] Streaming operation support
- [ ] Batch operation helpers
- [ ] Middleware support (logging, tracing, metrics)
- [ ] OpenTelemetry integration
- [ ] gRPC load balancing support
- [ ] Health check support
- [ ] Schema validation helpers

### Phase 3: Ecosystem Integration

- [ ] OpenTelemetry exporters
- [ ] Prometheus metrics
- [ ] Structured logging support
- [ ] Tracing integration examples
- [ ] Common middleware patterns

### Phase 4: Performance & Scalability

- [ ] Advanced caching strategies (LRU, distributed)
- [ ] Connection pooling optimization
- [ ] Request batching utilities
- [ ] Performance benchmarking suite
- [ ] Optimization guides and best practices

### Phase 5: Enterprise Features

- [ ] Custom authentication strategies
- [ ] Request signing
- [ ] API key rotation support
- [ ] Audit logging
- [ ] Circuit breaker pattern
- [ ] Retry policies with backoff

## Under Consideration

- GraphQL integration layer
- Async/streaming query builder patterns
- Interceptor middleware system
- Multi-region/failover support
- Custom cache implementations guide
- Protocol buffer code generation helpers

## Community Input

We value community feedback! Please share:
- Feature requests via GitHub Issues
- Use case descriptions
- Pain points in current implementation
- Integration opportunities

## Timeline

- **Q1 2024**: Core operations (Phase 1)
- **Q2 2024**: Advanced features start (Phase 2)
- **Q3 2024**: Ecosystem integration (Phase 3)
- **Q4 2024**: Performance focus (Phase 4)
- **2025+**: Enterprise features and expansions

## How to Contribute

Interested in working on roadmap items? We welcome contributions!

1. **Pick an item** from the roadmap
2. **Discuss** your approach in GitHub Issues
3. **Submit a PR** with your implementation
4. **Get reviewed** by maintainers
5. **Merge** and celebrate! 🎉

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## Feedback

Your feedback shapes this roadmap! Please:
- Comment on issues with your use cases
- Suggest features in GitHub Discussions
- Share performance insights
- Report gaps in documentation

## Related Initiatives

- [SpiceDB Project](https://github.com/authzed/spicedb)
- [AuthZed SDK Roadmaps](https://authzed.com/docs)
- [Community Requests](https://authzed.com/community)

---

Last updated: January 2024

*This roadmap is subject to change based on community needs and priorities.*


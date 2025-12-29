# JIRA Ticket: Create Test Repository for Relationship-Aware File Batching

## Summary
Create a Node.js test repository with an inventory management system to validate relationship-aware file batching for LLM analysis

## Description
Create a comprehensive test repository that demonstrates clear file relationships (imports, layers, tests, migrations) to test the relationship-aware batching system. The repository should include multiple domains with cross-dependencies to validate clustering and batching algorithms.

## Acceptance Criteria
- [ ] Node.js project structure with package.json and dependencies
- [ ] Three distinct domains: Product, Order, User
- [ ] Clear layer separation: Models → Repositories → Services → Controllers → Routes
- [ ] Test files paired with corresponding source files
- [ ] Database migrations linked to models
- [ ] Cross-domain dependencies (Order depends on Product and User)
- [ ] Documentation (README.md and RELATIONSHIPS.md)
- [ ] All files have meaningful relationships that can be detected

## Technical Details

### Structure
- **Models**: Product, Order, User, Category
- **Repositories**: Database access layer
- **Services**: Business logic layer
- **Controllers**: HTTP request handlers
- **Routes**: Express route definitions
- **Tests**: Unit tests for models, services, controllers
- **Migrations**: SQL schema definitions

### Relationship Types
- **Strong**: Test pairing, migration links
- **Medium**: Layer dependencies, cross-domain imports
- **Weak**: Configuration, application setup

### Expected Clusters
1. Product Domain (11 files)
2. Order Domain (8 files + dependencies)
3. User Domain (8 files)

## Labels
- `testing`
- `batching`
- `nodejs`
- `test-repository`

## Priority
Medium

## Story Points
5


# Contributing to mautrix-viber

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## ⚠️ Important: Read Coding Standards First

**Before making any changes, please read [`.cursorrules`](.cursorrules)**. This file defines our comprehensive coding standards, best practices, and architectural patterns. All code must follow these guidelines.

Key areas covered:
- Go idioms and best practices
- Error handling patterns
- Testing requirements (80%+ coverage target)
- Security checklist
- Code quality gates (go vet, golangci-lint)
- Performance best practices
- And much more...

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/mautrix-viber.git`
3. Create a feature branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Run tests: `make test`
6. Run linters: `make lint`
7. Commit your changes: `git commit -m "Add feature: your feature"`
8. Push to your fork: `git push origin feature/your-feature`
9. Open a Pull Request

## Code Style

- Follow Go conventions and use `gofmt`
- Write clear, descriptive commit messages
- Add comments for exported functions and types
- Write unit tests for new functionality

## Testing

- Add tests for new features
- Ensure all existing tests pass
- Aim for high test coverage

## Documentation

- Update README.md if adding new features
- Add code comments for complex logic
- Update API documentation if needed

## Pull Requests

- Keep PRs focused on a single feature or bugfix
- Ensure all CI checks pass
- Request review from maintainers

## Questions?

Open an issue or reach out to maintainers.


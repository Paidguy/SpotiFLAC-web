# Contributing to SpotiFLAC

First off, thank you for considering contributing to SpotiFLAC! It's people like you that make SpotiFLAC such a great tool.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Commit Messages](#commit-messages)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

This project and everyone participating in it is governed by respect and professionalism. By participating, you are expected to uphold this code.

### Our Standards

- Be respectful and inclusive
- Accept constructive criticism gracefully
- Focus on what is best for the community
- Show empathy towards other community members

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples**
- **Describe the behavior you observed and what you expected**
- **Include screenshots if relevant**
- **Include your environment details** (OS, browser, Docker version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description of the proposed feature**
- **Explain why this enhancement would be useful**
- **List any alternatives you've considered**

### Your First Code Contribution

Unsure where to begin? You can start by looking through these issues:

- Issues labeled `good first issue` - simple issues that are good for beginners
- Issues labeled `help wanted` - issues that need attention

### Pull Requests

1. Fork the repository
2. Create a new branch from `main`
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.22 or later
- Node.js 20 or later
- pnpm package manager
- ffmpeg
- Git

### Setting Up Your Development Environment

1. **Fork and clone the repository**

```bash
git clone https://github.com/YOUR_USERNAME/SpotiFLAC-web.git
cd SpotiFLAC-web
```

2. **Install dependencies**

```bash
# Frontend dependencies
cd frontend
pnpm install
cd ..

# Backend dependencies
go mod download
```

3. **Start development servers**

```bash
# Terminal 1: Frontend (with hot reload)
cd frontend
pnpm run dev

# Terminal 2: Backend
ENV=development DOWNLOAD_PATH=./test-downloads go run .
```

4. **Access the application**

- Frontend dev server: http://localhost:5173
- Backend server: http://localhost:8080

## Pull Request Process

1. **Create a feature branch**

```bash
git checkout -b feature/amazing-feature
```

2. **Make your changes**
   - Write clean, readable code
   - Follow the coding standards
   - Add tests if applicable
   - Update documentation

3. **Test your changes**

```bash
# Frontend tests
cd frontend
pnpm test
pnpm lint

# Backend tests
go test ./...
go fmt ./...
```

4. **Commit your changes**

```bash
git add .
git commit -m "Add some amazing feature"
```

5. **Push to your fork**

```bash
git push origin feature/amazing-feature
```

6. **Open a Pull Request**
   - Provide a clear title and description
   - Reference any related issues
   - Wait for review and address feedback

### PR Review Process

- At least one maintainer will review your PR
- Address any requested changes
- Once approved, your PR will be merged
- Your contribution will be recognized in the project

## Coding Standards

### Go (Backend)

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused
- Handle errors explicitly

**Example:**

```go
// Good
func downloadTrack(ctx context.Context, trackID string) error {
    if trackID == "" {
        return errors.New("track ID cannot be empty")
    }
    // Implementation...
}

// Bad
func dt(t string) {
    // Implementation without error handling...
}
```

### TypeScript/React (Frontend)

- Use TypeScript types consistently
- Follow React best practices
- Use functional components with hooks
- Keep components small and reusable
- Use meaningful component and variable names

**Example:**

```typescript
// Good
interface TrackProps {
  trackName: string;
  artistName: string;
  onDownload: (trackId: string) => void;
}

export function Track({ trackName, artistName, onDownload }: TrackProps) {
  return (
    <div className="track">
      <h3>{trackName}</h3>
      <p>{artistName}</p>
    </div>
  );
}

// Bad
function T(p: any) {
  return <div>{p.n}</div>;
}
```

### Code Style

- **Indentation**: 2 spaces for TypeScript/React, tabs for Go
- **Line length**: Try to keep under 100 characters
- **Naming**:
  - Go: camelCase for private, PascalCase for public
  - TypeScript: camelCase for variables, PascalCase for components
- **Comments**: Write clear, concise comments for complex logic

## Commit Messages

Write clear, descriptive commit messages following this format:

```
<type>: <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat: Add batch download for playlists

Implements functionality to download all tracks in a playlist
with a single click. Includes progress tracking and error handling.

Closes #123
```

```
fix: Resolve crash on empty playlist

Fixed null pointer dereference when fetching playlists with
no tracks. Added validation and error messages.

Fixes #456
```

## Testing

### Frontend Tests

```bash
cd frontend
pnpm test
```

### Backend Tests

```bash
go test ./...
```

### Manual Testing

Before submitting a PR, manually test:

1. Basic functionality (download single track, album, playlist)
2. Error cases (invalid URLs, network errors)
3. UI interactions (buttons, forms, navigation)
4. Responsive design (mobile, tablet, desktop)

## Documentation

- Update README.md if you add new features
- Add JSDoc/GoDoc comments for new functions
- Update API documentation if you change endpoints
- Include examples for new functionality

### Documentation Standards

- **Clear and concise**: Get to the point quickly
- **Examples**: Include code examples where helpful
- **Keep it updated**: Ensure docs match the current code
- **Screenshots**: Add images for UI changes

## Questions?

If you have questions about contributing:

- Open an issue with the `question` label
- Join the discussions on GitHub
- Check existing documentation and issues

## Recognition

Contributors will be recognized in:

- GitHub contributors list
- Release notes (for significant contributions)
- Project documentation (for major features)

## License

By contributing to SpotiFLAC, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to SpotiFLAC! 🎵

**Maintained with ❤️ by [@Paidguy](https://github.com/Paidguy)**

**Built upon the original work by [@afkarxyz](https://github.com/afkarxyz)**

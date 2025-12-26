# Git Setup Instructions for Comply360

## 📦 Initial Repository Setup

Follow these exact steps to push your Comply360 project to GitHub:

### Step 1: Navigate to Project Directory

```bash
cd comply360
```

### Step 2: Initialize Git Repository

```bash
git init
```

### Step 3: Configure Git (if not already done)

```bash
git config user.name "Rodrick Makore"
git config user.email "your-email@example.com"
```

### Step 4: Create .gitignore

The `.gitignore` file is already created, but verify it contains:
```
node_modules/
.env
.env.local
dist/
build/
.next/
*.log
.DS_Store
```

### Step 5: Create GitHub Repository

1. Go to https://github.com/retrocraftdevops
2. Click "New Repository"
3. Name: `comply360`
4. Description: "SADC Corporate Gateway - Multi-tenant SaaS for company registration"
5. **Keep it Private** (this is proprietary software)
6. **DO NOT** initialize with README (we have our own)
7. Click "Create repository"

### Step 6: Add All Files to Git

```bash
git add .
```

### Step 7: Create Initial Commit

```bash
git commit -m "feat: initial commit - Comply360 enterprise setup

- Complete documentation suite (PRD, TDD, Vision)
- Project structure with Next.js + Go
- Multi-tenant PostgreSQL schema
- Docker compose development environment
- AI code validation system
- CI/CD pipeline configuration
- Comprehensive README and setup guides"
```

### Step 8: Set Main Branch

```bash
git branch -M main
```

### Step 9: Add Remote Repository

```bash
git remote add origin https://github.com/retrocraftdevops/comply360.git
```

### Step 10: Push to GitHub

```bash
git push -u origin main
```

---

## 🔄 Subsequent Commits

After initial setup, use these commands for regular development:

```bash
# Check status
git status

# Add specific files
git add <file1> <file2>

# Or add all changes
git add .

# Commit with conventional commit message
git commit -m "feat: add name reservation wizard"
git commit -m "fix: resolve multi-tenant isolation bug"
git commit -m "docs: update API documentation"

# Push to GitHub
git push
```

---

## 🌿 Branch Strategy

### Main Branches

- `main` - Production-ready code
- `develop` - Development branch for integration

### Feature Branches

```bash
# Create feature branch
git checkout -b feature/name-reservation-wizard

# Work on feature...
git add .
git commit -m "feat: implement name search functionality"

# Push feature branch
git push -u origin feature/name-reservation-wizard

# Create Pull Request on GitHub
# After PR approved and merged, delete local branch
git checkout main
git pull
git branch -d feature/name-reservation-wizard
```

### Branch Naming Convention

- `feature/` - New features (e.g., `feature/payment-integration`)
- `fix/` - Bug fixes (e.g., `fix/authentication-error`)
- `docs/` - Documentation (e.g., `docs/api-specification`)
- `refactor/` - Code refactoring (e.g., `refactor/database-queries`)
- `test/` - Test additions (e.g., `test/registration-workflow`)

---

## 📝 Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style (formatting, no logic change)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**

```bash
git commit -m "feat(auth): add two-factor authentication"

git commit -m "fix(api): resolve tenant isolation bug in registration query"

git commit -m "docs(readme): add deployment instructions"

git commit -m "refactor(database): optimize registration queries with indexes"

git commit -m "test(registration): add unit tests for name validation"

git commit -m "chore(deps): update Next.js to v14.1.0"
```

---

## 🔐 SSH Setup (Recommended)

For easier authentication, setup SSH keys:

```bash
# Generate SSH key (if you don't have one)
ssh-keygen -t ed25519 -C "your-email@example.com"

# Start ssh-agent
eval "$(ssh-agent -s)"

# Add your SSH key
ssh-add ~/.ssh/id_ed25519

# Copy public key
cat ~/.ssh/id_ed25519.pub
# Copy the output

# Add to GitHub:
# 1. Go to GitHub Settings > SSH and GPG keys
# 2. Click "New SSH key"
# 3. Paste your public key
# 4. Save

# Update remote to use SSH
git remote set-url origin git@github.com:retrocraftdevops/comply360.git
```

---

## 🛠️ Useful Git Commands

```bash
# View commit history
git log --oneline --graph

# View changes
git diff

# Undo last commit (keep changes)
git reset --soft HEAD~1

# Undo last commit (discard changes)
git reset --hard HEAD~1

# Stash changes temporarily
git stash
git stash pop

# Create tag for release
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0

# View branches
git branch -a

# Delete remote branch
git push origin --delete feature/old-feature

# Pull latest changes
git pull

# Merge develop into main
git checkout main
git merge develop
git push
```

---

## 🚨 Important Reminders

1. **Never commit `.env` files** - They contain secrets
2. **Never commit `node_modules/`** - They're in .gitignore
3. **Always create feature branches** - Don't work directly on main
4. **Write meaningful commit messages** - Follow conventional commits
5. **Pull before you push** - Avoid merge conflicts
6. **Review changes before committing** - Use `git status` and `git diff`

---

## ✅ Verification

After pushing, verify on GitHub:

1. Go to https://github.com/retrocraftdevops/comply360
2. Check that all files are there
3. Check that README displays correctly
4. Verify documentation folder is present
5. Confirm .env is NOT committed (should be in .gitignore)

---

## 📞 Need Help?

If you encounter any issues:
1. Check error messages carefully
2. Google the error message
3. Ask team members on Slack: #comply360-dev
4. Email: dev@comply360.com

---

**Happy Coding! 🚀**

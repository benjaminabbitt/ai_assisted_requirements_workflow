# CI Configuration

This document provides examples and guidance for configuring continuous integration to validate Gherkin specifications and enforce workflow rules.

---

## Overview

CI validates that:
1. All scenarios have required `@story-{id}` tags
2. Implementation PRs don't merge with `@pending` tags for that story
3. Context files stay current (freshness monitoring)
4. Implemented scenarios pass (not `@pending`)

---

## GitHub Actions Example

Example workflow for Gherkin validation. **Adjust paths to match your directory structure.**

```yaml
name: Gherkin Validation

on:
  pull_request:
    paths: ['features/**', 'context/**']  # Adjust to your structure

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Require @story tag
        run: |
          # Adjust path to your features directory
          for f in features/**/*.feature; do
            grep -q "@story-" "$f" || { echo "Missing @story: $f"; exit 1; }
          done

      - name: Block impl PR with @pending for this story
        if: startsWith(github.head_ref, 'impl/')
        run: |
          # Extract story ID from branch name (e.g., impl/PROJ-1234 -> PROJ-1234)
          STORY_ID=$(echo "${{ github.head_ref }}" | sed 's/^impl\///')

          # Check if any scenarios have @pending for this story
          if grep -r "@pending.*@story-$STORY_ID\|@story-$STORY_ID.*@pending" features/; then
            echo "ERROR: Implementation PR contains @pending scenarios for $STORY_ID"
            echo ""
            echo "The following scenarios still have @pending tags:"
            grep -rn "@pending.*@story-$STORY_ID\|@story-$STORY_ID.*@pending" features/
            echo ""
            echo "Remove @pending tags for this story before merging (keep @story-$STORY_ID)"
            echo ""
            echo "Note: @pending and @story-{id} are two separate tags."
            echo "After implementation, only @story-$STORY_ID should remain."
            exit 1
          fi

          echo "✓ No @pending tags found for story $STORY_ID"

      - name: Validate context file freshness
        run: |
          for f in context/*.md; do
            days=$(( ($(date +%s) - $(stat -c %Y "$f")) / 86400 ))
            if [ $days -gt 60 ]; then
              echo "WARNING: $f not updated in $days days"
            fi
          done

      - name: Run implemented scenarios (not @pending)
        run: cucumber --tags "not @pending"
```

---

## Gherkin/Cucumber PR Automation

Example showing requirements-analyst and bo-review integration:

```yaml
name: Gherkin/Cucumber PR Automation

on:
  issues:
    types: [labeled]
  pull_request:
    paths:
      - 'features/**/*.feature'

jobs:
  # requirements-analyst: Runs when ticket labeled "ready-for-spec"
  draft-spec:
    if: github.event.label.name == 'ready-for-spec'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run requirements-analyst
        run: |
          claude --agent requirements-analyst \
                 --input "Ticket: ${{ github.event.issue.html_url }}"

      - name: Create Gherkin/Cucumber PR
        run: |
          # Agent creates spec/PROJ-1234 branch with .feature files
          # Opens PR for BO review (triggers bo-review below)

  # bo-review: Runs on Gherkin/Cucumber PRs (spec/ branches)
  bo-review-assistant:
    if: startsWith(github.head_ref, 'spec/')
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run bo-review agent
        run: |
          claude --agent bo-review \
                 --input "PR #${{ github.event.pull_request.number }}" \
                 --context business.md

      - name: Post review report as PR comment
        run: |
          gh pr comment ${{ github.event.pull_request.number }} \
             --body-file review-report.md

  block-pending-impl:
    runs-on: ubuntu-latest
    if: startsWith(github.head_ref, 'impl/')
    steps:
      - uses: actions/checkout@v4

      - name: Check for @pending tags
        run: |
          STORY_ID=$(echo "${{ github.head_ref }}" | sed 's/^impl\///')
          if grep -r "@pending.*@story-$STORY_ID\|@story-$STORY_ID.*@pending" features/; then
            echo "ERROR: Implementation PR contains @pending for $STORY_ID"
            exit 1
          fi
```

---

## GitLab CI Example

```yaml
# .gitlab-ci.yml

stages:
  - validate
  - test

validate-gherkin:
  stage: validate
  script:
    - |
      for f in features/**/*.feature; do
        grep -q "@story-" "$f" || { echo "Missing @story: $f"; exit 1; }
      done
  only:
    changes:
      - features/**/*.feature

block-pending-impl:
  stage: validate
  script:
    - |
      if [[ $CI_COMMIT_REF_NAME == impl/* ]]; then
        STORY_ID="${CI_COMMIT_REF_NAME#impl/}"
        if grep -r "@pending.*@story-$STORY_ID\|@story-$STORY_ID.*@pending" features/; then
          echo "ERROR: Implementation PR contains @pending for $STORY_ID"
          exit 1
        fi
      fi
  only:
    - merge_requests

run-scenarios:
  stage: test
  script:
    - cucumber --tags "not @pending"
  only:
    changes:
      - features/**/*.feature
```

---

## Azure Pipelines Example

```yaml
# azure-pipelines.yml

trigger:
  branches:
    include:
      - main
  paths:
    include:
      - features/**

pr:
  branches:
    include:
      - main
  paths:
    include:
      - features/**

jobs:
- job: ValidateGherkin
  pool:
    vmImage: 'ubuntu-latest'
  steps:
  - bash: |
      for f in features/**/*.feature; do
        grep -q "@story-" "$f" || { echo "Missing @story: $f"; exit 1; }
      done
    displayName: 'Require @story tags'

  - bash: |
      if [[ $(Build.SourceBranch) == refs/heads/impl/* ]]; then
        STORY_ID="${Build.SourceBranch#refs/heads/impl/}"
        if grep -r "@pending.*@story-$STORY_ID\|@story-$STORY_ID.*@pending" features/; then
          echo "ERROR: Implementation PR contains @pending for $STORY_ID"
          exit 1
        fi
      fi
    displayName: 'Block @pending in impl PRs'
    condition: startsWith(variables['Build.SourceBranch'], 'refs/heads/impl/')

  - bash: cucumber --tags "not @pending"
    displayName: 'Run implemented scenarios'
```

---

## Validation Rules

### Required Tags

**Rule:** All scenarios must have `@story-{id}` tag

**Rationale:** Traceability to business request

**Enforcement:**
```bash
for f in features/**/*.feature; do
  grep -q "@story-" "$f" || { echo "Missing @story: $f"; exit 1; }
done
```

### Block @pending in Implementation PRs

**Rule:** Implementation PRs (`impl/{ticket-id}`) cannot merge if scenarios for that story still have `@pending` tags

**Rationale:** Prevents incomplete implementation from merging

**Enforcement:**
```bash
STORY_ID=$(echo "$BRANCH_NAME" | sed 's/^impl\///')
if grep -r "@pending.*@story-$STORY_ID\|@story-$STORY_ID.*@pending" features/; then
  echo "ERROR: Remove @pending tags for story $STORY_ID"
  exit 1
fi
```

### BO Approval for .feature Files

**Rule:** All `.feature` file changes require Business Owner approval

**Rationale:** Business alignment guaranteed

**Enforcement:** CODEOWNERS + branch protection
```
# CODEOWNERS
features/**/*.feature @business-owner
```

### Context File Freshness

**Rule:** Warning if context file not updated in 30 days, alert at 60 days

**Rationale:** Prevent stale documentation

**Enforcement:**
```bash
for f in context/*.md; do
  days=$(( ($(date +%s) - $(stat -c %Y "$f")) / 86400 ))
  if [ $days -gt 60 ]; then
    echo "WARNING: $f not updated in $days days"
  fi
done
```

---

## Running Scenarios

### Skip @pending Scenarios

Only run implemented scenarios:

```bash
# Cucumber
cucumber --tags "not @pending"

# Godog (Go)
godog run --tags "not @pending"

# Behave (Python)
behave --tags=-pending

# Cucumber-JVM (Java)
mvn test -Dcucumber.filter.tags="not @pending"
```

### Run Specific Story

Run scenarios for a specific story only:

```bash
cucumber --tags "@story-PROJ-1234"
```

### Run Multiple Stories

```bash
cucumber --tags "@story-PROJ-1234 or @story-PROJ-1235"
```

---

## Standards Compliance Check

**Automated standards-compliance agent** runs on implementation PRs to check IoC patterns and technical standards adherence.

```yaml
standards-check:
  runs-on: ubuntu-latest
  if: startsWith(github.head_ref, 'impl/')
  steps:
    - uses: actions/checkout@v4

    - name: Run standards-compliance agent
      run: |
        claude --agent standards-compliance \
               --input "$(git diff --name-only origin/main... | grep '\.go$')" \
               --context tech_standards.md

    - name: Comment on PR if violations
      if: failure()
      run: |
        gh pr comment ${{ github.event.pull_request.number }} \
           --body-file compliance-report.md
```

---

## Branch Protection

Configure in your Git provider UI:

### GitHub

**Settings → Branches → Branch protection rules → `main`:**
- ✅ Require pull request reviews before merging
- ✅ Require status checks to pass before merging
  - Select: `validate-gherkin`, `block-pending-impl`, `run-scenarios`
- ✅ Require conversation resolution before merging
- ✅ Require linear history (squash or rebase merges only)
- ✅ Do not allow bypassing the above settings
- ✅ Require CODEOWNERS review

### GitLab

**Settings → Repository → Protected Branches → `main`:**
- Allowed to merge: Maintainers + Developers
- Allowed to push: No one
- Require approval from code owners: Yes
- Require pipelines to succeed: Yes

### Azure DevOps

**Project Settings → Repos → Policies → `main`:**
- Require pull request: Yes
- Require minimum number of reviewers: 1
- Check for linked work items: Yes
- Require successful build: Yes
- Automatically include code reviewers: From CODEOWNERS

---

## Monitoring CI Failures

### Common Failures

| Failure | Cause | Resolution |
|---------|-------|------------|
| Missing @story tag | Scenario doesn't link to ticket | Add `@story-{id}` tag |
| @pending in impl PR | Developer forgot to remove @pending | Remove `@pending` tags for that story |
| Context file stale | Not updated in 60+ days | Review and update or confirm current |
| Scenario fails | Implementation bug | Fix implementation, rerun tests |

### Alert Thresholds

- **High CI failure rate:** > 20% of PRs failing validation
- **Stuck PRs:** No review after 24 hours
- **Stale context files:** > 3 files not updated in 60+ days

---

## Next Steps

- **For workflow details:** See [Workflow Guide](workflow.md)
- **For source control strategy:** See [Source Control Guide](source-control.md)
- **For integration architecture:** See [Integration Guide](integration.md)
- **For troubleshooting:** See [Troubleshooting Guide](troubleshooting.md)

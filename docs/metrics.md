# Metrics & SLAs

**The goal is speed parity:** Requirements should not be the bottleneck when development is AI-assisted.

This document defines the metrics to track, target SLAs, and enforcement mechanisms to ensure the AI-augmented workflow delivers on its promise.

---

## Primary: AI Confidence Rate (Measures Context File Quality)

| Metric | Target | If Below Target |
|--------|--------|-----------------|
| AI confident (High + Medium) | Majority of stories | Context files incomplete—every miss is a story requiring human discussion |
| AI escalates (Low) | Minority of stories | Document missing systems/patterns |

**Why this is primary:** AI confidence directly determines how many stories can be processed without human wait time. Low confidence = slow requirements = development bottleneck. Track this metric over time—it should improve as context files mature.

**How to measure:**
- Count AI executions by confidence level (High, Medium, Low)
- Calculate percentage: `(High + Medium) / Total * 100`
- Track weekly trends

**Action triggers:**
- < 70% confident for 1 week → Review context files for gaps
- < 50% confident for 1 week → Pause and fix context files before continuing
- Escalation rate > 20% → Urgent context file improvement needed

---

## Secondary: Speed (Measures Process Efficiency)

| Metric | Target | Why It Matters |
|--------|--------|----------------|
| Story → PR ready (confident) | Same day | Requirements shouldn't wait overnight |
| PR → BO Approved | < 24 hours | Review backlog = development blocked |
| Escalation resolution | 2-3 days | Complex cases still need to be fast |

**How to measure:**
- Timestamp when ticket labeled "ready-for-spec"
- Timestamp when spec PR created
- Timestamp when BO approves PR
- Calculate duration for each stage

**Dashboard example:**
```
┌─────────────────────────────────────────┐
│ Requirements Velocity (Last 7 Days)    │
├─────────────────────────────────────────┤
│ Avg Story → PR:        4.2 hours ✓     │
│ Avg PR → Approved:     8.5 hours ✓     │
│ Avg Escalation Time:   1.8 days  ✓     │
│                                         │
│ Stories Processed:     47               │
│ AI Confident:          38 (81%) ✓      │
│ AI Escalated:           9 (19%) ⚠      │
└─────────────────────────────────────────┘
```

---

## Quality (Measures Accuracy)

| Metric | Target | If Above Target |
|--------|--------|-----------------|
| BO overrides AI | < 10% | AI missing context or miscalibrated |
| Post-merge reverts | < 5% | Review process not catching issues |
| Implementation rejects spec | < 3% | Feasibility assessment weak |

**Why these matter:**
- **BO overrides:** AI should draft acceptable specs most of the time
- **Post-merge reverts:** Specs should be correct when merged
- **Implementation rejects:** Specs should be implementable

**How to measure:**
- Track PRs where BO requested significant changes
- Track spec amendments or reverts after merge
- Track implementation PRs that identified spec issues

**Quality vs. Speed tradeoff:** This workflow optimizes for speed without sacrificing quality. If quality metrics slip, the problem is usually incomplete context files, not the workflow itself.

---

## Review Timeline Targets (SLAs)

| AI Confidence | Review SLA | Escalation After |
|---------------|------------|------------------|
| High | 4 hours | 8 hours |
| Medium | Same business day | 24 hours |
| Post-escalation | 1 business day | 48 hours |

**Purpose:** Ensure specs don't sit in review queue, blocking development.

---

## Enforcement Mechanisms

### Automated Reminders

- **At 50% of SLA:** Slack/Teams notification to assigned reviewers
- **At 75% of SLA:** Escalation to team lead
- **At 100% of SLA:** Visible in team dashboard, blocks new spec PRs for that reviewer

### Auto-Merge Option

For HIGH confidence PRs only:

- If CI passes and no review within 24 hours, auto-merge with notification
- Team can disable auto-merge if preferred
- Auto-merged PRs flagged for retrospective review

---

## Escalation Response SLAs

| Escalation Type | Response SLA | Resolver |
|-----------------|--------------|----------|
| Missing architecture info | 1 business day | Architect/Tech Lead |
| Business rule clarification | 4 hours | Product Owner |
| Compliance question | 2 business days | Compliance team |
| Security concern | 1 business day | Security team |

**How to measure:**
- Timestamp when AI posts escalation to ticket
- Timestamp when expert responds
- Calculate duration

**Action triggers:**
- Escalation SLA missed > 3 times/week → Add capacity or simplify process
- Same escalation type repeatedly → Update context file to prevent future escalations

---

## Context File Freshness

| Metric | Warning Threshold | Alert Threshold |
|--------|-------------------|-----------------|
| Days since last update | 30 days | 60 days |

**Why this matters:** Stale context files lead to higher escalation rates and incorrect specs.

**Enforcement:**
- **Warning (30 days):** CI job creates ticket for owner to review
- **Alert (60 days):** Blocks new spec PRs until context file reviewed
- **Mechanism:** Weekly scheduled CI job checks last-modified dates

**How to implement:**
```bash
for f in context/*.md; do
  days=$(( ($(date +%s) - $(stat -c %Y "$f")) / 86400 ))
  if [ $days -gt 60 ]; then
    echo "ALERT: $f not updated in $days days"
    # Create blocking issue
  elif [ $days -gt 30 ]; then
    echo "WARNING: $f not updated in $days days"
    # Create reminder issue
  fi
done
```

---

## Dashboard Example

### Team Overview

```
┌──────────────────────────────────────────────────────────────┐
│  AI-Augmented Requirements Dashboard                        │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  📊 This Sprint (Week 2 of 2)                               │
│  ├─ Stories Processed:        23                            │
│  ├─ AI Confidence Rate:       78% (18 High/Medium) ✓       │
│  ├─ Escalation Rate:          22% (5 escalated)    ⚠       │
│  └─ Avg Story→Approved:       6.2 hours            ✓       │
│                                                              │
│  🎯 Quality Metrics                                         │
│  ├─ BO Override Rate:         4%  (1 of 23)        ✓       │
│  ├─ Post-Merge Reverts:       0%  (0 of 23)        ✓       │
│  └─ Impl Rejects Spec:        0%  (0 of 23)        ✓       │
│                                                              │
│  ⚡ Speed Metrics                                           │
│  ├─ Story → PR Created:       3.8 hours avg        ✓       │
│  ├─ PR → BO Approved:         2.4 hours avg        ✓       │
│  └─ Escalation Resolution:    1.2 days avg         ✓       │
│                                                              │
│  📁 Context File Health                                     │
│  ├─ business.md:              12 days ago          ✓       │
│  ├─ architecture.md:          45 days ago          ⚠       │
│  ├─ testing.md:               8 days ago           ✓       │
│  └─ tech_standards.md:        67 days ago          🔴      │
│                                                              │
│  🔔 Active Alerts                                           │
│  ├─ tech_standards.md overdue for review (67 days)         │
│  └─ Escalation rate above 20% threshold                    │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### Individual Metrics

```
┌──────────────────────────────────────────────────────────────┐
│  Story: PROJ-1234 - Password Reset Feature                  │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  📍 Timeline                                                │
│  ├─ Created:              Mon 10:23 AM                      │
│  ├─ Labeled ready:        Mon 10:25 AM                      │
│  ├─ AI draft created:     Mon 10:28 AM  (3 min) ✓          │
│  ├─ BO approved:          Mon 2:45 PM   (4h 17m) ✓         │
│  ├─ Spec merged:          Mon 2:46 PM                       │
│  ├─ Impl started:         Tue 9:00 AM                       │
│  ├─ Impl PR opened:       Tue 3:15 PM                       │
│  └─ Impl merged:          Tue 4:30 PM                       │
│                                                              │
│  🤖 AI Analysis                                             │
│  ├─ Confidence:           HIGH ✓                            │
│  ├─ Scenarios generated:  8 (6 core, 2 boundary)           │
│  ├─ Steps reused:         12 of 18 (67%)                   │
│  └─ Context versions:     business.md v1.2.0               │
│                            testing.md v3.1.0                │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

---

## Tracking Implementation

### Data Collection

**Ticket system webhook:**
- Capture timestamp when labeled "ready-for-spec"
- Capture timestamp when AI posts draft or escalation

**PR events:**
- Capture timestamp when spec PR created
- Capture timestamp when BO approves
- Capture timestamp when PR merges
- Capture timestamp when implementation PR opens
- Capture timestamp when implementation PR merges

**AI execution logs:**
- Confidence level (High, Medium, Low)
- Context versions used
- Number of scenarios generated
- Steps reused vs created

**Store in:**
- Database (PostgreSQL, MySQL)
- Time-series DB (InfluxDB, Prometheus)
- Spreadsheet (Google Sheets, Excel) for small teams

### Visualization

**Recommended tools:**
- Grafana (dashboards for time-series data)
- Tableau/PowerBI (executive dashboards)
- Custom web dashboard (React/Vue + Chart.js)

---

## Alerts

Configure alerts for:

### Critical (Page immediately)
- Escalation SLA missed > 48 hours
- CI blocking due to context file staleness
- Quality metric regression (> 20% override rate)

### Warning (Notify next business day)
- AI confidence rate < 70% for 1 week
- Review SLA missed (but < 48 hours)
- Context file approaching staleness (30-59 days)

### Info (Weekly digest)
- Speed metrics trending (faster or slower)
- Step reuse rate
- Scenarios per story trending

---

## Continuous Improvement

### Weekly Review

Review metrics in sprint retro:
- Are we meeting SLA targets?
- Is AI confidence improving?
- Are escalations decreasing?
- Which context files need updates?

### Monthly Review

Deeper analysis:
- Which story types escalate most often?
- Which context files correlate with high confidence?
- Are quality metrics stable?
- Is the team getting faster?

### Quarterly Review

Strategic assessment:
- Has AI-augmented workflow delivered promised speed?
- Are humans spending less time in meetings?
- Is the escalation → codification loop working?
- Should we adjust SLAs or targets?

---

## Comparison: Traditional vs AI-Augmented

| Metric | Traditional BDD | AI-Augmented | Improvement |
|--------|-----------------|--------------|-------------|
| Time per story | 2-4 hours | 5-30 minutes | 75-95% faster |
| Stories per sprint | 10-15 | 30-50 | 2-3x throughput |
| Meeting time | 30-60 min/story | 0 min/story (async review) | 100% reduction |
| Spec quality | Variable | Consistent (context-driven) | Higher consistency |
| Edge case coverage | Manual | Automated (boundary patterns) | More comprehensive |

---

## Next Steps

- **For workflow details:** See [Workflow Guide](workflow.md)
- **For context file maintenance:** See [Context Files Guide](context-files.md)
- **For troubleshooting:** See [Troubleshooting Guide](troubleshooting.md)

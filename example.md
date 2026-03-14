# SDD — Example Conversations

## Greenfield: Task management API

```
user: /sdd brainstorm

claude: What are you building?

user: A task management API. Users create projects, add tasks with priorities
      and due dates, assign tasks to team members. REST API, PostgreSQL.

claude: Decomposing into domains...

      Created 4 specs:
      - spec-users.md (4 requirements, 12 acceptance criteria)
      - spec-projects.md (6 requirements, 18 acceptance criteria)
      - spec-tasks.md (9 requirements, 31 acceptance criteria)
      - spec-notifications.md (3 requirements, 8 acceptance criteria)

      Dependency order: users → projects → tasks → notifications
      Next: /sdd plan

user: /sdd plan

claude: Reading 4 specs, 22 requirements...

      Generated frontier: 34 tasks, 5 tiers
      - Tier 0: DB schema, user model (3 tasks)
      - Tier 1: User CRUD, project model (5 tasks)
      - Tier 2: Project CRUD, task model (8 tasks)
      - Tier 3: Task CRUD, assignments (10 tasks)
      - Tier 4: Notifications (8 tasks)

      Next: /sdd execute

user: /sdd execute

claude: 🔄 Loop activated! 34 tasks, 20 max iterations.

      --- Iteration 1 ---
      T-001: Database Schema Setup (Tier 0, spec-users.md R1)
      Build PASS, Tests PASS
      Committed: feat(users): add User model and migration

      --- Iteration 2 ---
      T-002: User Registration Endpoint (Tier 0, spec-users.md R1)
      ...

      --- Iteration 18 ---
      All tasks done. Build passes. Tests pass.
      <promise>SPEC COMPLETE</promise>

      ✅ Done. 34 tasks, 18 iterations.
```

## Brownfield: Adding collaboration to existing Next.js app

```
user: /sdd brainstorm --from-code

claude: Exploring codebase... Next.js 14, Prisma, NextAuth.

      Found 6 domains. Created specs:
      - spec-auth.md (5 requirements — all satisfied by existing code)
      - spec-documents.md (7 requirements — all satisfied)
      - spec-editor.md (6 requirements — all satisfied)
      - spec-collaboration.md (3 requirements — all [GAP], not implemented)
      - spec-comments.md (4 requirements — all satisfied)
      - spec-notifications.md (3 requirements — 2 satisfied, 1 [GAP])

      4 requirements are gaps (not yet implemented).

user: I only want the collaboration features.
      /sdd plan --filter collaboration

claude: Reading spec-collaboration.md (3 requirements)...

      Generated frontier: 8 tasks, 3 tiers
      - Tier 0: WebSocket setup, presence types
      - Tier 1: Cursor presence, CRDT integration
      - Tier 2: Conflict resolution, live sync

user: /sdd execute --filter collaboration

claude: 🔄 Loop activated! 8 tasks.

      --- Iteration 1 ---
      T-001: WebSocket Server (Tier 0, R1)
      ...

      --- Iteration 8 ---
      All done. <promise>SPEC COMPLETE</promise>

user: Now notifications with adversarial review.
      /sdd plan --filter notifications
      /sdd execute --filter notifications --adversarial

claude: Archived previous cycle → context/impl/archive/20260314-201500/
      Configured Codex (gpt-5.4) as MCP adversary.

      --- Iteration 1 (BUILD) ---
      T-001: Notification Service (R1)
      Committed.

      --- Iteration 2 (REVIEW) ---
      Calling Codex...
      Finding: HIGH — missing retry logic for email delivery
      Fixed. Committed.

      --- Iteration 6 (REVIEW) ---
      No new findings. All tasks done.
      <promise>SPEC COMPLETE</promise>
```

## The flow

```
/sdd brainstorm  →  specs
/sdd plan        →  frontier
/sdd execute     →  code
```

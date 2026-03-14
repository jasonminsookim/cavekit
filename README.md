# SDD — Spec-Driven Development

Claude Code plugin + tmux dashboard for spec-driven development with automated iteration loops.

## Install

```bash
git clone https://github.com/JuliusBrussee/sdd-os.git ~/.sdd
~/.sdd/install.sh
```

This registers the Claude Code plugin, installs the `sdd` CLI, and makes all scripts executable.

## Terminal: tmux monitor

```bash
sdd --monitor                     # launch dashboard + claude + build loop
sdd --monitor --adversarial       # with Codex adversarial review
sdd --monitor --filter v2         # scope to v2 specs
sdd --kill                        # tear down the session
```

Opens a tmux session with three panes:
- **Left (70%)** — Claude Code running `/sdd-execute`
- **Top-right** — live progress: tasks done, tiers, progress bar
- **Bottom-right** — live activity: iteration log, git commits

## Claude commands

```
/sdd-brainstorm    →  write specs (the WHAT)
/sdd-plan          →  generate frontier (the ORDER)
/sdd-execute       →  ralph loop (the BUILD)
/sdd-review        →  gap analysis + adversarial review (the CHECK)
```

### Brainstorm — write specs

```bash
/sdd-brainstorm                    # interactive — asks what to build
/sdd-brainstorm context/refs/      # from PRDs, API docs, research
/sdd-brainstorm --from-code        # from existing codebase
```

### Plan — generate frontier

```bash
/sdd-plan                          # all specs
/sdd-plan --filter v2              # only v2 specs
```

### Execute — run the loop

```bash
/sdd-execute                       # implement everything
/sdd-execute --adversarial         # add Codex (GPT-5.4) review
/sdd-execute --max-iterations 30
```

### Review — post-loop check

```bash
/sdd-review                        # gap analysis + adversarial review
```

## All commands

| Command | Description |
|---------|-------------|
| `/sdd-brainstorm` | Write specs |
| `/sdd-plan` | Generate feature frontier |
| `/sdd-execute` | Implementation loop |
| `/sdd-review` | Gap analysis + adversarial review |
| `/sdd-progress` | Check frontier progress |
| `/sdd-gap-analysis` | Compare built vs intended |
| `/sdd-back-propagate` | Trace manual fixes to specs |
| `/sdd-help` | Show usage |

## Update

```bash
cd ~/.sdd && git pull
```

## License

MIT

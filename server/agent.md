# Commenting Guidelines for Maintainers

Read the code **as a maintainer**, not as a teacher.

**Before writing any comments**, fully understand:
- Data flow
- Invariants
- Transaction boundaries
- Concurrency assumptions
- Durability guarantees
- Failure and rollback semantics

Then add **only comments that encode intent and correctness**, not narration.

## Comment Rules

**Do NOT:**
- Explain syntax or obvious control flow
- Restate what the code literally does
- Add step numbers or tutorial explanations
- Add comments where names already explain intent

**Do:**
- Explain **why** something exists
- Document **invariants and guarantees**
- Clarify **transaction and concurrency boundaries**
- Explain **failure, rollback, and recovery behavior**
- State **assumptions that future changes must not violate**
- Note **tradeoffs** (performance vs safety) when relevant

## Scope

Focus comments on:
- Service layers
- Transaction boundaries
- Infrastructure wiring
- Durability logic
- Concurrency-sensitive code

Skip:
- Generated code
- Trivial helpers
- Self-explanatory getters/setters

## Output Requirements

- Modify **comments only**; do not change logic
- Keep comments **minimal but precise**
- Write comments as if this code will be maintained by someone else in 6 months
- Every comment must reduce future bug risk

If a comment does not add real semantic value, **do not write it**.
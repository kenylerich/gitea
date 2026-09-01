---
name: k7-trim-cot-leakage
description: "Use when auditing or fixing prose that reads like a leaked reasoning transcript in the gitea repo: dead design-session citations, change narration ('used to', 'no longer', 'this cut'), stack or review vantage ('a later PR in this stack', 'rejected in review'), reviewer-addressed justification, control-flow narration, or hedged planning residue, in Go and TS comments, docs/*.md, commit messages, PR descriptions, go-swagger annotations, and options/locale/locale_en-US.json."
---

# Trimming chain-of-thought leakage

Chain-of-thought leakage is prose whose vantage is the authoring session rather than the repository: it cites artifacts only that session could see, narrates the change instead of the state, or argues with a reviewer who has left. The fix is never deletion alone when a passage carries factual clauses — restate each so it stands at HEAD, then delete the transcript around it; a passage carrying none (an audit code, control-flow narration) is deleted outright. **REQUIRED BACKGROUND:** [k7-prose-standard](../k7-prose-standard/SKILL.md) owns the complete-proposition rule this skill applies. It is guidance, not a script.

## The one test

For every suspect passage ask: **could a reader at HEAD, with no access to any session transcript, PR thread, or uncommitted draft, resolve every reference and verify every claim?** If no, restate the surviving facts from the repository's vantage and delete the rest. If yes, it is not leakage, however historical it sounds — but resolvability only clears this skill's bar: on current-state surfaces (`docs/*.md`, code comments, `options/locale/locale_en-US.json`, go-swagger comments) a resolvable change story is still change narration, and class 3 routes it to its sanctioned home.

## Taxonomy

1. **Dead design-session citations** — `(decision N)`, `audit §N`, `phase T4`, "the design ledger", "plan §1.4". If the decision has a committed owner, cite it by path or full URL; otherwise delete the citation and restate its factual clause to stand alone.
2. **Stack and PR vantage** — "a later PR in this stack", "this PR adds", "the previous commit". State the shipped mechanism or the extension point; deferred work moves to a `TODO` marker or an issue reference by full URL.
3. **Change narration and version stamps** — "used to", "no longer", "the old X", and indexical stamps ("v1", "this cut", "today" contrasting with a past state). State the present behavior; a fixed regression becomes a present-tense counterfactual ("without X, Y happens"), never repo history ("used to Y").
4. **Review choreography** — "Rejected in review:", "the reviewer confirmed", draft ordinals ("v5 of this note"), round attributions. Keep the surviving decision and rationale as plain fact; delete who said it when.
5. **Reviewer-addressed justification** — "the cast is safe — it simply…", "this is correct because…". A comment arguing its own correctness addresses a reviewer, not a maintainer. State the invariant that makes the code safe, or delete the comment if the code shows it.
6. **Restatement and derivation transcripts** — control-flow narration ("first we X, then we Y"), test walkthroughs, proofs of obvious branches. Delete; keep only a non-obvious contract or invariant.
7. **Hedges and planning residue** — "probably fine for now", "should be enough", deferrals with no owner. Promote to `TODO`/`FIXME` or restate as the actual bound; delete the hedge.
8. **Authoring-language slips** — untranslated working-language fragments in prose whose language is otherwise English, or the reverse in a zh counterpart. Translate or delete.

## What is not leakage

Unaided citation passes fail in both directions by deleting durable references and keeping dead ones. Apply these keep rules as written; [examples](references/examples.md) calibrates each:

- **Issue references** — a full URL (`https://github.com/go-gitea/gitea/issues/1470`), `TODO(name):`, or "issue #N owns the follow-up" resolve at HEAD; keep them on any surface, including `docs/*.md`. Per `AGENTS.md`, reference issues and PRs by full URL, not by number — when the prose is permanent, prefer the full URL.
- **Merged-PR citations inside postmortems and design docs** — sanctioned evidence when the citation names the merged PR by full URL and the citation still resolves.
- **Suppression justifications** — `//nolint:revive // reason`, `//nolint:gocritic // reason`, coverage-ignore reasons, `// swagger:operation` annotations naming their scope, `// TODO(name):` with a reason. The justification clause is required prose; fix a false reason, never delete it.
- **Counterfactual-present regression pins** — "without X, Y happens", "a naive X would…".
- **Measured bounds** — "(measured: 512 nests ≈ 0.15s)" calibrating a constant; the provenance word "measured" is load-bearing.
- **Runtime old/new states** — "the old connection drains before the new one accepts" is runtime lifecycle, not change history.
- **External references that resolve outside the repo by design** — standards sections (RFC 9110 §10.1.5), Go module paths (`gitea.dev/backend/...`), and committed docs that own their §-numbering may be cited by section.
- **Project voice and genre forms** — "we" as project voice; the "Alternatives considered" section of a design doc.

## Workflow

1. Scope and exclusions per [k7-prose-standard](../k7-prose-standard/SKILL.md): require an explicit scope; never touch `vendor/` or `frontend/web_src/js/vendor/`. Generated artifacts (`templates/swagger/v1-swagger.generated.json`, `templates/swagger/v1-openapi3.generated.json`, files under `public/assets/`, `options/bindata.go`) are derivatives, not prose targets: change the owning source or scenario and regenerate them only when an authorized behavior change requires new evidence.
2. Audit read-only first: run the [recall batteries](references/recall-batteries.md), calibrating each probe against a known positive and a near-miss negative before trusting its output, then judge every hit semantically. The batteries are probes, not the definition — also read the densest prose in scope (Go doc comments on exported identifiers, the `templates/shared/*.tmpl` attribute blocks, swagger annotations on `routers/api/v1/`, `options/locale/locale_en-US.json` keys with `TODO` markers, `docs/*.md`) without a pattern in hand.
3. Fix owner-first per surface: generated swagger → trace every consumer, fix the go-swagger comment on the route, then run `make generate-swagger` and `make swagger-validate`; go-swagger JSON schema → fix the owning struct in `modules/structs/` and the option registration under `routers/api/v1/swagger/`; model- or user-visible strings → route through [k7-prose-standard](../k7-prose-standard/SKILL.md) and change only with owning behavior evidence, otherwise leave unchanged and report the deferral.
4. Before deleting anything, enumerate the passage's propositions (prose-standard) and check the [overcorrection traps](references/examples.md#overcorrection-traps): trims that flip an obligation into an endorsement, promote a hypothetical to a shipped feature, delete a true fact, or drop provenance.
5. Verify: re-run the batteries expecting only sanctioned keeps, this skill's own directory, and any quoted evidence the owning author has carried into a doc; confirm every remaining citation resolves at HEAD; run the gates for touched surfaces (`make lint-md` for `docs/*.md`, `make lint-go` for Go comments, `make generate-swagger` and `make swagger-validate` for API changes).
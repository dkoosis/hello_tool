# End-of-Session Roadblock Review

Date: {{YYYY-MM-DD}}
Session Focus: {{Briefly describe the main goals or tasks of the work session}}

---

## 1. Roadblocks & Time Sinks

* **Significant Roadblocks:**
    * Were there any issues, bugs, or unexpected behaviors that took more than 5-10 minutes to resolve or work around during this session?
    * List them here:
        * Roadblock 1:
        * Roadblock 2:
        * ...
* **Time Sinks:**
    * Were there any tasks that felt like they took longer than they should have, even if not a complete "blocker"?
    * List them here:
        * Time Sink 1:
        * ...

---

## 2. Analysis & Prevention (For each significant roadblock identified above)

### Roadblock: {{Name of Roadblock 1}}

* **Symptoms:**
    * How did this issue manifest?
* **Root Cause Investigation:**
    * What was the underlying cause of this issue? (Be as specific as possible)
* **Time to Resolve:**
    * Approximately how long did it take to debug and resolve/workaround?
* **Proactive Prevention Ideas:**
    * **Documentation:**
        * Could better comments in code have helped? Where?
        * Is there a missing piece in `README.md`, an ADR, or other project documentation? What should be added/clarified?
        * Could a wiki page or troubleshooting note be useful?
    * **Automation/Tooling:**
        * Could a `Makefile` target be added or improved to catch or prevent this?
        * Could a shell script or other helper tool assist?
        * Would a new linter rule or static analysis check be relevant?
    * **Testing:**
        * Could a new unit test, integration test, or end-to-end test have caught this earlier or made the root cause more obvious?
        * Describe the test case.
    * **Code/Configuration Changes:**
        * Could the code be refactored for clarity or robustness in this area?
        * Could configuration be simplified or validated better?
    * **Process/Workflow:**
        * Could a change in our development process or workflow help avoid this type of issue?
* **Simplest Next Step for Prevention:**
    * What is the one most impactful but relatively simple thing we could do (e.g., add a specific comment, log an issue for a new test, draft a small doc update) to reduce the chance of this recurring?

---

*(Repeat section 2 for each identified roadblock)*

---

## 3. General Observations & Workflow Improvements

* Were there any parts of the workflow today that felt clunky or inefficient?
* Any "aha!" moments or learnings that should be shared or documented more broadly?
* Are there any tools or techniques we discussed that we should explore further?

---

## 4. Action Items

*List any concrete action items resulting from this review, assign owners, and estimate priority/effort if possible.*

* [ ] Action Item 1 (e.g., "Draft ADR for X", "Add test case for Y to `tests/`", "Update README section on Z") - Owner: - Priority:
* [ ] Action Item 2 - Owner: - Priority:
* ...

---
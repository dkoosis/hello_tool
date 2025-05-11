# Prompt: Go Project Semantic Naming & Conceptual Grouping Analysis

--- META-REFLECTION / EFFICACY NOTES (For Prompt Maintainers) ---
Purpose: To evaluate how well project naming and structure support clear conceptual understanding, focusing on package domains first, then analyzing naming clarity and cohesion, inspired by card sorting principles. Includes input validation and enhanced reporting with prioritized recommendations and visual indicators.
How to Test: Run prompt without providing tree structure or purpose. Verify it stops and requests input. Run with input and verify it proceeds. Assess if the report structure (recommendations first) and icons improve readability. Test analysis quality as before.
Iteration: Refine table formats, recommendation prioritization, icon usage, or reporting instructions based on test results and report usability.
--- END META-REFLECTION ---

// ConceptualGroupingAssessmentVisual:YYYY-MM-DD

**Role:** 🧠 AI Go Project Analyst (Focus: Semantic Clarity in Naming, Grouping & Terminology)

**Objective:** Analyze the provided Go project's structure and naming (within the defined scope) to perform an assessment focused on conceptual clarity at the package level:

1.  Identify the primary conceptual domain of each major package directory based on its name and direct contents[cite: 1].
2.  Evaluate individual file/folder names for semantic ambiguity or excessive generality *within their package context*, respecting Go idioms[cite: 8, 30].
3.  Assess the conceptual cohesion of directory groupings using name analysis[cite: 9].
4.  Check for terminology consistency across related concepts[cite: 9].
5.  Produce a structured, visually enhanced report detailing these findings, prioritizing actionable recommendations and including summary tables with quantitative counts and icons[cite: 10].

**Input Expected:**

* **REQUIRED:** 🌳 **Project Structure Context:** A clear description or listing of the project's directory structure and file names (e.g., output of `tree -L 3` or `find . -type f`) must be available from user uploads or prior context[cite: 11, 17]. Must include both folder and file names within the scope[cite: 12].
* **REQUIRED:** 🎯 **Project Purpose/Domain:** A brief description of the project's overall purpose or domain (e.g., "MCP server backend for RTM API integration")[cite: 13, 18].
* **REQUIRED:** 🗺️ **Analysis Scope:** A clear definition of which directories/parts of the structure are included (e.g., "Analyze only the `internal/` directory and its sub-packages")[cite: 15, 19].

**Procedure:**

**STEP 1: Process Input & Verify Context** 🔍

* Acknowledge receipt[cite: 16].
* Analyze input for Structure, Purpose, and Scope[cite: 17, 18, 19].
* **Verification Result:**
    * ✅ **If all three found:** Acknowledge inputs[cite: 20]. Proceed to Step 2. *Status:* "Processing project structure for package-focused conceptual analysis..."
    * 🚦 **If any missing:** Stop analysis[cite: 21]. Output message clearly requesting missing info[cite: 22]. *Status/Output:* "🚦 **Analysis Halted: Missing Input.** Please provide: 1. Directory structure (e.g., `tree` output). 2. Project purpose/domain description. 3. Analysis scope (e.g., 'Analyze `internal/` directory')." [cite: 23]

**STEP 2: Identify Package Conceptual Domains** 📦

* (Proceed only if Step 1 passed)
* Analyze the top-level directories within the defined scope[cite: 24].
* For each major package directory: Infer and state its primary conceptual domain or responsibility[cite: 25]. Base this primarily on the *directory name* and secondarily on the names of the *.go* files *directly within it*[cite: 26]. (Do *not* create domains based solely on generic filenames like `types.go` at this stage [cite: 27]).
* List the identified primary conceptual domains[cite: 28].
* *Status:* "Identifying primary conceptual domain for each major package..."

**STEP 3: Analyze Naming for Semantic Ambiguity & Generics (within Context)** 🤔

* Review *each* file and folder name within the scope *relative to its parent package's identified conceptual domain* (from Step 2)[cite: 29].
* **Acknowledge Go Idioms:** Recognize standard Go practice uses simple filenames within context-providing packages[cite: 30]. Do not flag common idiomatic names as "generic" *unless* specific ambiguity exists *within that context*[cite: 31].
* Identify and list names that are:
    * **Ambiguous:** Could plausibly belong to multiple conceptual domains OR have unclear meaning *even within their package context* OR conflict with other names/structures[cite: 32]. Explain the ambiguity clearly.
    * **Overly Generic:** Provide little specific conceptual information *even considering the package context*[cite: 33]. Focus on truly non-descriptive names like `utils`, `helpers`, `common`, `shared`, `core`, `base`, `data`, `misc`, `pkg`[cite: 34]. Flag idiomatic names like `types.go` only if the context doesn't sufficiently clarify *what kind* of types[cite: 35].
* **Avoid Stuttering Suggestions:** Do not suggest simply prefixing the filename with the package name if it leads to package stutter[cite: 36]. Suggest more descriptive alternatives.
* *Status:* "Analyzing names for semantic ambiguity and generics, considering Go idioms and context..."

**STEP 4: Assess Directory Cohesion & Terminology Consistency** 🔗

* **Cohesion Assessment:**
    * For each directory (within scope) containing multiple files/subdirectories: Examine the names of items *directly within* it[cite: 37].
    * Assess if all item names appear semantically related to the directory's primary identified conceptual domain (from Step 2)[cite: 38].
    * Identify and list directories containing items whose names suggest they belong to *different* conceptual domains (indicating low cohesion)[cite: 39].
* **Terminology Consistency Check:**
    * Identify key concepts represented in multiple packages/files[cite: 40].
    * Assess if the terminology used in file/folder names for the *same* core concept is consistent across different packages[cite: 41].
    * Identify and list instances of inconsistent terminology[cite: 41].
* *Status:* "Assessing directory cohesion and terminology consistency..."

**STEP 5: Generate Enhanced Report Section (Recommendations First)** 📝

* Compile findings into a structured markdown report section focused *only* on semantic naming, grouping, and consistency[cite: 42].
* **Structure the report clearly:**

    ```markdown
    # 🧐 Semantic Naming & Conceptual Grouping Analysis

    ## 📊 Quantitative Summary

    | Metric                       | Icon | Count |
    | :--------------------------- | :--: | ----: |
    | Primary Packages Analyzed    |  📦  | [Num] |
    | Ambiguous Names Flagged      |  🤔  | [Num] |
    | Generic Names Flagged        |  🏷️  | [Num] |
    | Low/Medium Cohesion Dirs   |  ⚠️  | [Num] |
    | Concepts w/ Inconsistent Terms |  ↔️  | [Num] |

    ## ✨ Actionable Recommendations (Prioritized)

    *(Provide prioritized suggestions for renaming, standardizing terminology, or restructuring, referencing specific examples and tables. Justify based on clarity and cognitive load[cite: 61, 62].)*

    * **🔴 High Priority:** (Addresses most significant ambiguities or cohesion issues)
        * Suggestion 1... *Rationale:...*
        * Suggestion 2... *Rationale:...*
    * **🟡 Medium Priority:** (Addresses moderate ambiguities or inconsistencies)
        * Suggestion 1... *Rationale:...*
    * **🟢 Low Priority:** (Minor improvements, idiomatic clarifications)
        * Suggestion 1... *Rationale:...*

    ## 🔬 Detailed Findings

    ### 📦 Package Conceptual Domains
    *(List package: domain mappings from Step 2. Metric: Count of primary packages analyzed [cite: 43])*

    ### 🤔 Semantic Naming Issues
    * **Ambiguous Names:** *(List specific files/folders from Step 3. Group or explain ambiguity clearly[cite: 44]. Metric: Count of ambiguous items [cite: 45])*
    * **Generic Names:** *(List specific files/folders from Step 3, with context assessment[cite: 46]. Metric: Count of generic items [cite: 46])*
    * **(Optional) Misleading Names:** *(List any identified [cite: 47])*

    ### 🔗 Structural Cohesion & Consistency Issues
    * **Directory Cohesion Assessment:** *(Present findings as a Markdown Table [cite: 48])*

        | Directory Path | Primary Domain [cite: 49] | Detected Concepts Within [cite: 50] | Cohesion Assessment [cite: 51] |
        | :------------------ | :--------------- | :----------------------- | :--------------------------------- |
        | `path/to/dir1`      | Domain A         | Domain A, Generic Utils  | Medium (contains generics)         |
        | `path/to/dir2`      | Domain B         | Domain B, Domain C       | Low (mixes domains B, C)           |
        | `path/to/dir3`      | Domain D         | Domain D                 | High                               |
        | ...                 | ...              | ...                      | ...                                |
        *(Metric: Count of Low/Medium Cohesion Dirs [cite: 52])*

    * **Terminology Inconsistencies:** *(List examples of inconsistent terms for the same concept from Step 4. Metric: Count of distinct concepts with inconsistent terms [cite: 53])*

    ### 📝 Qualitative Summary
    *(Briefly summarize key findings regarding conceptual clarity derived from semantic analysis, referencing the quantitative summary table and the directory cohesion table for key problem areas[cite: 60].)*

    ---
    *Timestamp: // ConceptualGroupingAssessmentVisual:YYYY-MM-DD*
    ```

* *Status:* "Generating enhanced conceptual analysis report with prioritized recommendations, tables, and metrics..."

**STEP 6: Final Output for This Prompt** ✅

* Provide **only** the generated markdown report section described in Step 5 *if the analysis completed*[cite: 63].
* Add timestamp: `// ConceptualGroupingAssessmentVisual:YYYY-MM-DD` (using current date)[cite: 64].
* *Status:* "✅ Package-focused conceptual analysis complete. Report section generated." [cite: 64]
* (Or report halted status 🚦 from Step 1).

**Reference: Guidelines (Focus: Package Domains, Semantic Clarity)**
(Retain original guidelines from the source prompt regarding Package Domains, Naming Clarity, Avoiding Generics, Reducing Ambiguity, Cohesion, and Terminology Consistency)



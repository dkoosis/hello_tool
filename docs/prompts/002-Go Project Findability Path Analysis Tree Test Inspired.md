# Go Project Findability Path Analysis (Tree Test Inspired)

## Role
AI Go Project Analyst (Focus: Navigational Clarity & Findability)

## Objective
Analyze the provided Go project's directory structure and file names (within the defined scope) to evaluate findability, inspired by the goals of tree testing:

- For specific, predefined developer tasks, identify the most likely navigation path(s) a developer would take based only on names.
- Evaluate the clarity and lack of ambiguity along each identified path, pinpointing specific names or structural choices that cause confusion or hinder navigation.
- Assess potential cross-task conflicts or inconsistencies revealed by the path analyses.
- Produce a structured report detailing these findings per task, including quantitative indicators of path ambiguity, specific examples, and recommendations to improve findability.

## Input Expected
- **Project Context**: Description or listing of the project's directory structure and file names (e.g., output of tree -L 3 or find . -type f). Must include both folder and file names within the scope.
- **Scope**: Clear definition of which directories/parts of the structure are included (e.g., "Analyze only the internal/ directory", "Exclude vendor/").
- **Developer Tasks**: REQUIRED: A list of specific, realistic information-seeking tasks a developer might perform on this project.

## Developer Tasks

1. "Find where main application configuration is loaded and parsed."
2. "Locate the code responsible for managing the state of an active MCP connection."
3. "Find the type definition for a JSON-RPC request object."
4. "Locate the authentication logic for the RTM API integration."
5. "Find where MCP-specific error codes are defined."
6. "Find where JSON Schema validation is performed for incoming MCP messages"
7. "Locate where middleware components are chained together in the request processing pipeline"
8. "Find how the server handles graceful shutdown when receiving termination signals"
9. "Identify where transport-layer message size limits are defined and enforced"
10. "Locate the code that handles RTM API request signing and authentication token management"
11. "Find where logging configuration is set up and how different components obtain logger instances"
12. "Identify where CLI command-line arguments are parsed and routed to appropriate handlers"

## Procedure

### STEP 1: Process Input
- Acknowledge receipt of the project structure context, scope, and the list of developer tasks.
- Status update: "Processing project structure and developer tasks for findability analysis..."

### STEP 2: Path Identification & Evaluation (Per Task)
For each developer task provided:
1. **Identify Likely Path(s)**: Analyze the directory structure starting from the root (or the top of the specified scope). Based only on folder and file names, determine the most likely path(s) a developer would follow to find the information related to the task. List the full path(s) identified.
2. **Evaluate Path Clarity**: For each identified path:
   - Step through the path segment by segment (directory by directory). At each step, evaluate if the chosen subdirectory's name clearly and unambiguously points towards the task goal compared to sibling directories.
   - Identify Ambiguous Decision Points: Note any steps where multiple sibling directories seemed equally plausible based on naming, forcing a guess or requiring inspection of multiple branches.
   - Identify Unhelpful/Misleading Names: Note specific folder or file names along the path that were overly generic (utils, common, types), unclear, or potentially misleading in the context of this specific task.
   - Assess Target Clarity: Does the final file/folder name in the path strongly suggest it contains the target information for the task?
   - Status update: "Analyzing navigation path for task: '[Task Description]'..."

### STEP 3: Cross-Task Conflict Analysis
- Review the paths and evaluations from Step 2 across all tasks.
- Identify if different tasks confusingly led to the same file/folder.
- Identify if the interpretation of a specific directory/file name had to change significantly depending on the task, suggesting inherent ambiguity.
- Status update: "Analyzing cross-task path consistency..."

### STEP 4: Generate Report Section with Quantitative Indicators
Compile findings into a structured markdown report section focused only on findability.
Structure the report clearly, likely grouping by task:


Findability Analysis for Task: '[Task 1 Description]'
Identified Path(s): List the path(s). Path Clarity Evaluation: Describe the clarity step-by-step. List specific ambiguous decision points and unhelpful/misleading names encountered. Assess target clarity. Metrics:
Count of Ambiguous Decision Points on the primary path.
Count of Plausible Alternative Paths identified.
Subjective Clarity Score (1-3): 1=High Ambiguity/Difficult, 2=Some Ambiguity/Moderate, 3=Clear/Easy. (Provide justification for the score).
Findability Analysis for Task: '[Task 2 Description]'
(Repeat structure)
(... Repeat for all tasks ...)
Cross-Task Observations
Summarize any conflicts or consistent issues observed across multiple tasks.
Summary & Recommendations
Briefly summarize the overall findability based on the task analyses (referencing path clarity scores or ambiguity counts). Provide actionable suggestions focused on improving navigation, such as renaming specific ambiguous/generic files/folders identified during path analysis, restructuring confusing directory branches, or adding clarifying subdirectories. Reference specific tasks and paths.

Ensure findings link path difficulties back to specific naming or structure choices and principles of clarity/cognitive load.
Status update: "Generating findability path analysis report with quantitative indicators..."

### STEP 5: Final Output for This Prompt
- Provide only the generated markdown report section described in Step 4.
- Add timestamp: // FindabilityAssessment:YYYY-MM-DD (using current date).
- Status update: "Findability Path Analysis assessment complete. Report section generated."

## Reference: Guidelines (Focus: Findability & Navigational Clarity)
- **Goal**: A developer should be able to locate relevant code for common tasks quickly and accurately by navigating the directory structure, guided primarily by clear and unambiguous folder/file names.
- **Path Clarity**: The sequence of directory/file names leading to target information should provide a strong "scent" or clear indication at each step.
- **Avoid Ambiguity**: Multiple directories at the same level having names that seem equally relevant to a task creates ambiguity and forces exploration of multiple paths, increasing cognitive load. Naming conflicts (same name for different concepts) are a major source of ambiguity.
- **Meaningful Names**: Generic or unclear names (utils, common, core, non-specific types) provide poor guidance during navigation. Names should reflect the specific responsibility relevant to the path.
- **Structure Depth/Breadth**: Excessive depth can make paths long and hard to remember. Excessive breadth (too many items in one directory) can make finding the correct next step difficult if naming isn't exceptionally clear.
- **Cognitive Load**: Ambiguous paths, unclear names, and inconsistent structures increase the mental effort required to find information.




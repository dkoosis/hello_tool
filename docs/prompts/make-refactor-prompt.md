Utility Script Development Plan for Makefile Integration
Global Requirements for All Scripts

Robust Shell Selection

Use #!/usr/bin/env bash shebang
Include set -euo pipefail for strict error handling
Set LC_ALL=C for consistent output


Self-contained Environment Checks

Verify required dependencies at runtime
Provide helpful error messages for missing tools
Handle GNU/BSD variants of common utilities


Consistent Formatting

Import shared formatting functions
Maintain visual consistency across all output
Support both verbose and quiet modes


Error Handling

Return meaningful exit codes (0=success, 1=error, 2=config error)
Capture and format stderr from external commands
Log errors in a structured format


Idempotent Execution

Scripts should be safely re-runnable
Implement proper cleanup of temporary files



Script Development Specifications
1. check_environment.sh
Purpose: Validate environment prerequisites

Detect shell type (bash/zsh/sh)
Check for all required utilities
Identify GNU vs BSD variants
Verify minimum versions of critical tools (git, go, etc.)
Generate summary report of environment status

2. print_utils.sh
Purpose: Centralize all output formatting

Define color constants and icon variables
Implement functions for step/success/warning/error/info messages
Create external command output boundaries
Support indentation levels
Include progress indicators for long-running operations

3. version_utils.sh
Purpose: Extract and manage version information

Retrieve git version info (tag, commit hash, dirty status)
Format build timestamps consistently
Generate LDFLAGS for Go builds
Handle version comparison for compatibility checks

4. build.sh
Purpose: Handle binary compilation

Configure build flags and environment variables
Support multiple targets/architectures
Implement clean build option
Generate metadata files for build artifacts
Parameters: binary_name, module_path, version, commit, output_dir

5. dependency_manager.sh
Purpose: Manage project dependencies

Execute go mod tidy with proper flags
Download and verify dependencies
Validate go.mod module path
Check for vulnerable dependencies
Parameters: module_path, expected_version

6. code_quality.sh
Purpose: Run linting and formatting

Execute golangci-lint with proper configuration
Format code according to standards
Check line lengths in source files
Validate YAML files
Parameters: warn_threshold, error_threshold, paths_to_check

7. test_runner.sh
Purpose: Execute and report on tests

Run tests with proper flags (race detection, timeouts)
Generate coverage reports
Format test output for readability
Support test filtering and debugging modes
Parameters: test_mode, coverage_path, package_pattern

8. deployment.sh
Purpose: Handle cloud deployment

Build proper substitution strings for Cloud Build
Submit to Google Cloud Build
Format build output
Monitor deployment progress
Retrieve service URL
Parameters: service_name, project_id, region, artifact_repo, version, commit

9. health_checker.sh
Purpose: Verify deployed service health

Make HTTP requests to health endpoints
Validate response content and status codes
Verify version information matches expected values
Implement retry logic with timeouts
Parameters: endpoint_url, expected_version, expected_commit, max_retries

Implementation Notes

Begin implementation with print_utils.sh and check_environment.sh as foundations
Each script should have a usage function triggered by --help
Support both direct invocation and sourcing by other scripts
Include version stamps in script headers
Maintain comprehensive error messages
Document all arguments and return values
Use consistent parameter naming across scripts
Keep logging verbosity controllable via environment variables

This plan provides a framework for developing modular, robust utility scripts that integrate seamlessly with Make while maintaining excellent user experience through consistent, visually clear output.
# Deployment targets
# makefiles/deploy.mk

.PHONY: deploy health-check

# --- Cloud Deployment ---
deploy: all check-vulns lint-yaml
	@printf "$(ICON_START) $(BOLD)$(BLUE)Deploying $(SERVICE_NAME) to Google Cloud...$(NC)\n"
	@if [ -z "$(PROJECT_ID)" ]; then \
		printf "  $(ICON_FAIL) $(RED)Error: Google Cloud Project ID not found. Set via 'gcloud config set project YOUR_PROJECT_ID' or ensure gcloud is configured.$(NC)\n"; \
		exit 1; \
	fi
	@printf "  $(ICON_INFO) Project: $(PROJECT_ID)\n"
	@# Determine the path to cloudbuild.yaml (prefer build/cloudbuild/cloudbuild.yaml)
	@CONFIG_PATH=""; \
	if [ -f "build/cloudbuild/cloudbuild.yaml" ]; then \
		CONFIG_PATH="./build/cloudbuild/cloudbuild.yaml"; \
	elif [ -f "./cloudbuild.yaml" ]; then \
		CONFIG_PATH="./cloudbuild.yaml"; \
	else \
		printf "  $(ICON_FAIL) $(RED)Error: Cloud Build config file not found at './build/cloudbuild/cloudbuild.yaml' or './cloudbuild.yaml'.$(NC)\n"; \
		exit 1; \
	fi; \
	printf "  $(ICON_INFO) Using Cloud Build config: '%s'\n" "$$CONFIG_PATH"; \
	printf "  $(ICON_INFO) DEBUG: Substitutions string:\n>>> $(GCLOUD_BUILD_SUBSTITUTIONS) <<<\n"; \
	printf "  $(ICON_INFO) $(YELLOW)Submitting to Google Cloud Build and awaiting completion...$(NC)\n"; \
	printf "  $(ICON_INFO) $(YELLOW)Begin gcloud output:----------------------------------------------$(NC)\n"; \
	if gcloud builds submit . --config="$$CONFIG_PATH" --project=$(PROJECT_ID) --substitutions="$(GCLOUD_BUILD_SUBSTITUTIONS)"; then \
		printf "  $(ICON_INFO) $(YELLOW)End gcloud output:------------------------------------------------$(NC)\n"; \
		printf "  $(ICON_OK) $(GREEN)Cloud Build completed successfully.$(NC)\n"; \
		printf "  $(ICON_INFO) Monitor detailed build logs at: https://console.cloud.google.com/cloud-build/builds?project=$(PROJECT_ID)\n"; \
		printf "  $(ICON_START) $(BLUE)Fetching deployed service URL...$(NC)\n"; \
		SERVICE_URL=$$(gcloud run services describe $(SERVICE_NAME) --platform=managed --region=$(GCP_REGION) --project=$(PROJECT_ID) --format="value(status.url)" 2>/dev/null); \
		if [ -n "$$SERVICE_URL" ]; then \
			printf "  $(ICON_OK) $(GREEN)Service URL: $$SERVICE_URL$(NC)\n"; \
			make health-check HEALTH_CHECK_URL="$$SERVICE_URL/health" EXPECTED_VERSION="$(LOCAL_VERSION)" EXPECTED_COMMIT="$(LOCAL_COMMIT_HASH)"; \
		else \
			printf "  $(ICON_WARN) $(YELLOW)Could not retrieve service URL. Skipping health check. Please check the Cloud Run console.$(NC)\n"; \
		fi; \
	else \
		printf "  $(ICON_INFO) $(YELLOW)End gcloud output:------------------------------------------------$(NC)\n"; \
		printf "  $(ICON_FAIL) $(RED)Cloud Build submission or execution failed. Review gcloud output above and build logs for details.$(NC)\n"; \
		exit 1; \
	fi
	@printf "\n"

# --- Health Check ---
health-check:
ifndef HEALTH_CHECK_URL
	$(error HEALTH_CHECK_URL is not set. Usage: make health-check HEALTH_CHECK_URL=... EXPECTED_VERSION=... EXPECTED_COMMIT=...)
endif
ifndef EXPECTED_VERSION
	$(error EXPECTED_VERSION is not set)
endif
ifndef EXPECTED_COMMIT
	$(error EXPECTED_COMMIT is not set)
endif
	@printf "$(ICON_START) $(BOLD)$(BLUE)Performing health check (Expected Version: $(EXPECTED_VERSION), Commit: $(EXPECTED_COMMIT)) on $(HEALTH_CHECK_URL)...$(NC)\n"
	@SUCCESS=false; \
	for i in $$(seq 1 5); do \
		printf "  $(ICON_INFO) Health check attempt #$$i...\n"; \
		RESPONSE=$$(curl -s -w "\n%{http_code}" "$(HEALTH_CHECK_URL)"); \
		HTTP_CODE=$$(echo "$$RESPONSE" | tail -n1); \
		BODY=$$(echo "$$RESPONSE" | sed '$$d'); \
		if [ "$$HTTP_CODE" -eq 200 ]; then \
			RESPONSE_VERSION=$$(echo "$$BODY" | jq -r .version); \
			RESPONSE_COMMIT=$$(echo "$$BODY" | jq -r .commit); \
			if [ "$$RESPONSE_VERSION" = "$(EXPECTED_VERSION)" ] && [ "$$RESPONSE_COMMIT" = "$(EXPECTED_COMMIT)" ]; then \
				printf "  $(ICON_OK) $(GREEN)Health check PASSED. Version and Commit match expected values.$(NC)\n"; \
				printf "    Response: %s\n" "$$BODY"; \
				SUCCESS=true; break; \
			else \
				printf "  $(ICON_WARN) $(YELLOW)Health check attempt #$$i: Service responded OK, but content mismatch.$(NC)\n"; \
				printf "    Expected Version: $(EXPECTED_VERSION), Got: $$RESPONSE_VERSION\n"; \
				printf "    Expected Commit: $(EXPECTED_COMMIT), Got: $$RESPONSE_COMMIT\n"; \
				printf "    Full Response: %s\n" "$$BODY"; \
			fi \
		else \
			printf "  $(ICON_WARN) $(YELLOW)Health check attempt #$$i: Failed with HTTP status $$HTTP_CODE.$(NC)\n"; \
			printf "    Response Body: %s\n" "$$BODY"; \
		fi; \
		if [ $$i -lt 5 ]; then printf "  $(ICON_INFO) Retrying in 10 seconds... $(NC)\n"; sleep 10; fi; \
	done; \
	if [ "$$SUCCESS" = "false" ]; then \
		printf "  $(ICON_FAIL) $(RED)$(BOLD)Post-deployment health check FAILED after 5 attempts.$(NC)\n"; \
		exit 1; \
	fi
	@printf "\n"

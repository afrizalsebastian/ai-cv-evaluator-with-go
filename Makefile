generate:
	@echo "Creating directory...."
	@mkdir -p ./internal/generated

	@echo "Generating Gorilla Mux Go Code"
	@oapi-codegen -generate types,gorilla -package generated ./api/openapi.yaml > ./internal/generated/api.gen.go || (echo "Failed to generate Go code"; exit 1)

	@echo "Generation complete!"
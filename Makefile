

.PHONY: fmt fmt-frontend fmt-backend
fmt: fmt-frontend fmt-backend
	
fmt-frontend:
	cd frontend && npm run fmt

fmt-backend:
	cd backend && $(MAKE) fmt

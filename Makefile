#Config
main_pkg := cmd/music_store
bin_name := music_store.bin
proj_dir := $(shell pwd)

# Extract project info
mod_name := $(shell head -1 go.mod | cut -d' ' -f2)
proj_dir := $(shell pwd)
proj_name := $(shell basename "$(PWD)")
git_hash := $(shell git rev-parse --short=16 HEAD)
git_branch := $(shell git rev-parse --abbrev-ref HEAD)
git_stat := $(shell if [ "$$(git diff --stat)" != '' ]; then echo "dirty"; else echo "clean"; fi)
git_tag := $(shell tag=$$(git describe --tags --abbrev=0 2>/dev/null); if [ $$? -ne 0 ]; then echo "latest"; else echo "$$tag"; fi)
build_date:= $(shell date +%s)

_print:
	@echo "--------------------------------------------"
	@echo "              $(proj_name) MAKE FILE              "
	@echo "--------------------------------------------"
	@printf "Bin:			%s\n" "$(bin_name)"
	@printf "Module:			%s\n" "$(mod_name)"
	@printf "Proj Dir:		%s\n" "$(proj_dir)"
	@printf "Proj Name:		%s\n" "$(proj_name)"
	@printf "Main Pkg:		%s\n" "$(main_pkg)"
	@printf "Git Hash:		%s\n" "$(git_hash)"
	@printf "Git Branch:		%s\n" "$(git_branch)"
	@printf "Git Stat:		%s\n" "$(git_stat)"
	@printf "Git Tag:		%s\n" "$(git_tag)"
	@printf "Build Date:		%s\n" "$(build_date)"
	@printf "Mode:			%s\n" "$(mode)"
	@echo "--------------------------------------------"
	@printf "VERSION : %s\n" "$(version)"
	@echo "--------------------------------------------"

# -- Go Build Functions
_go_clean_bin:
	@echo " > Delete: rm -f $(proj_dir)/bin/$(bin_name) 2>/dev/null"
	@rm -f $(proj_dir)/bin/$(bin_name) 2>/dev/null

_go_ensure:
	@echo " > Ensure: go mod tidy"
	@go mod tidy -compat=1.19

_go_test:
	@echo " > Run Test: go test ./..."
	@echo " > --------- "
	@go test ./...
	@echo " > --------- "

_go_build:
	@echo " > Build Binary: env GOOS=linux GOARCH=amd64 go build -o $(bin_name)  $(mod_name)/$(main_pkg)"
	@env  GOSUMDB=off GOOS=linux GOARCH=amd64 go build -o $(bin_name) $(mod_name)/$(main_pkg)

_go_build_mac:
	@echo " > Build Binary: env GOOS=darwin GOARCH=amd64 go build -o $(bin_name)  $(mod_name)/$(main_pkg)"
	@env  GOSUMDB=off GOOS=darwin GOARCH=amd64 go build -o $(bin_name) $(mod_name)/$(main_pkg)

_move_bin:
	@echo " > Make Dir:  mkdir -p $(proj_dir)/bin"
	@mkdir -p ./bin
	@echo " > Move: mv $(bin_name) $(proj_dir)/bin"
	@mv $(bin_name) $(proj_dir)/bin

_build_post:
	@echo " > --------------------------------------------"
	@echo " > Go Binary Generated "
	@echo " > --------------------------------------------"

_build:
	@-$(MAKE) --no-print-directory _print _go_clean_bin _go_ensure _go_test _go_build_mac _move_bin _build_post

_build_linux:
	@-$(MAKE) --no-print-directory _print _go_clean_bin _go_ensure _go_test _go_build _move_bin _build_post

run:
	./bin/music_store.bin start

.PHONY: _go_clean_bin _go_ensure _go_test _go_build _go_build_mac _move_bin _build_post _build

# docker build
ifndef docker_bin
	$(error "missing docker bin, can't docker build")
endif
	@echo " > Found: [docker] binary"

_docker_build:
	@echo " > Build: docker build -t $(proj_name) $(proj_dir)"
	@docker build -t $(proj_name):latest $(proj_dir)

## docker: builds docker image
docker:
	@-$(MAKE) --no-print-directory _print _docker_build

run-dev:
	@go run ./...

up:
	@docker compose up -d

down:
	@docker compose down

up-redis:
	@docker compose up -d redis

test-api: up-redis
	@bash scripts/test_api.sh

seed: up-redis
	@bash scripts/seed_data.sh

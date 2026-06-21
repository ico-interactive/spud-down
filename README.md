# spud down bot
monorepo in the future

## developer setup

### prerequisites:
- mise-en-place
  - for auto environment variables

### setup environment
```sh
git clone https://github.com/ico-interactive/spud-down.git
cd spud-down
cp mise.toml mise.local.toml
# enter ur secrets
vim mise.local.toml
go mod download
```

### build and run
```sh
go run ./bot
```

## deploy with compose
```sh
git clone https://github.com/ico-interactive/spud-down.git
cd spud-down
# enter ur secrets
vim compose.yaml
docker compose up -d
```

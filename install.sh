#!/usr/bin/env sh
set -eu

REPO="github.com/inezhinskiy/pre-commit-gitleaks"
BINARY_PATH="./cmd/gitleaks-hook"

echo "[install] Перевірка наявності Go..."
if ! command -v go >/dev/null 2>&1; then
  echo "[install] Помилка: Go не знайдено. Встанови Go (https://go.dev/dl/) і повтори." >&2
  exit 1
fi

TARGET_REPO="$(pwd)"
if [ ! -d "$TARGET_REPO/.git" ]; then
  echo "[install] Помилка: $TARGET_REPO не є git-репозиторієм." >&2
  exit 1
fi

echo "[install] Завантаження та збірка pre-commit hook..."
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

git clone --depth 1 "https://${REPO}.git" "$TMP_DIR/src"

( cd "$TMP_DIR/src" && go build -o "$TARGET_REPO/.git/hooks/pre-commit" "$BINARY_PATH" )
chmod +x "$TARGET_REPO/.git/hooks/pre-commit"

echo "[install] Увімкнення автоматичного встановлення gitleaks..."
git config hooks.gitleaks.autoinstall true

echo "[install] Готово! pre-commit hook встановлено в $TARGET_REPO/.git/hooks/pre-commit"
echo "[install] Перевірка секретів (gitleaks) працюватиме автоматично перед кожним комітом."

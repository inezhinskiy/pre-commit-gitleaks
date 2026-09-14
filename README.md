# pre-commit-gitleaks

Pre-commit hook на Go, що автоматично перевіряє staged-зміни на наявність секретів
(API-токенів, ключів, паролів тощо) перед кожним комітом за допомогою
[gitleaks](https://github.com/gitleaks/gitleaks). Якщо секрет знайдено — коміт відхиляється.

## Швидке встановлення (рекомендовано)

Виконайте в корені вашого git-репозиторію:

```bash
curl -sSL https://raw.githubusercontent.com/inezhinskiy/pre-commit-gitleaks/main/install.sh | sh
```

Скрипт:
1. перевіряє наявність Go;
2. клонує та білдить hook;
3. встановлює його як `.git/hooks/pre-commit` у поточному репозиторії;
4. вмикає автоматичне встановлення `gitleaks` (якщо він ще не встановлений) через
   `git config hooks.gitleaks.autoinstall true`.

## Ручне встановлення

Якщо не хочете використовувати `curl | sh`:

```bash
git clone https://github.com/inezhinskiy/pre-commit-gitleaks.git /tmp/pcg
go build -o .git/hooks/pre-commit /tmp/pcg/cmd/gitleaks-hook
chmod +x .git/hooks/pre-commit
```

## Режими роботи

**Без автовстановлення (за замовчуванням).**
Hook очікує, що `gitleaks` вже встановлений і доступний у `PATH`. Якщо його немає —
коміт блокується з підказкою, як встановити вручну:
https://github.com/gitleaks/gitleaks#installing

**З автовстановленням.**
```bash
git config hooks.gitleaks.autoinstall true
```
Якщо `gitleaks` відсутній, hook сам встановить його залежно від ОС:
- **macOS** — через `brew`
- **Linux** — через `apt-get`, або через `go install` як fallback
- **Windows** — через `choco` або `scoop`
- **будь-яка ОС** — універсальний fallback через `go install github.com/zricethezav/gitleaks/v8@latest`,
  якщо жоден пакетний менеджер не знайдено

Щоб вимкнути автовстановлення:
```bash
git config hooks.gitleaks.autoinstall false
```

## Як це працює

Hook викликає:
```bash
gitleaks protect --staged --redact --no-banner
```

Команда сканує саме **staged**-зміни (те, що піде в коміт), а не всю історію репозиторію.
`--redact` приховує саме значення секрету в логах, залишаючи лише тип правила, яке спрацювало.

> **Важливо:** `gitleaks` використовує keyword-префільтр — правило застосовується лише
> якщо в рядку присутнє відповідне ключове слово (наприклад, `telegram` для правила
> `telegram-bot-api-token`). Змінна на кшталт `TELE_TOKEN` це правило не активує;
> потрібно `TELEGRAM_BOT_TOKEN` чи подібне.

## Перевірка роботи (приклад із Telegram bot token)

```bash
echo 'TELEGRAM_BOT_TOKEN="1234567890:AAExampleFakeTokenForTestingOnly350"  # gitleaks:allow' > test.env
git add test.env
git commit -m "test"
```

Очікуваний результат — коміт відхилено:
```
WRN leaks found: 1
[pre-commit:gitleaks] gitleaks виявив секрети у staged-змінах. Коміт відхилено.
```

Приберіть тестовий файл після перевірки:
```bash
git reset test.env
rm test.env
```

## Структура репозиторію

```
pre-commit-gitleaks/
├── cmd/gitleaks-hook/main.go   # логіка hook'а: перевірка секретів + автовстановлення gitleaks
├── install.sh                  # bootstrap-інсталятор через curl | sh
├── go.mod
└── README.md
```

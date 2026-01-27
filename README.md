# Light Defender Launcher

Лаунчер для Light Defender.

С помощью него через удобный интерфейс вы можете:
- Устанавливать Light Defender на удалённый сервер.
- Менять конфигурации на удалённом сервере.
- Обновляться до новых версий за пару кликов.


- [Текстовая инструкция по использованию лаунчера](./INSTRUCTIONS/LAUNCHER/LAUNCHER.md);
- [Текстовая инструкция по установке Light Defender вручную](./INSTRUCTIONS/MANUAL/MANUAL.md).

## Compiling

1) Клонируем репозиторий

```bash
git clone https://github.com/fridalif/light-defender-launcher
```

2) Устанавливаем зависимости фронтэнда

```bash
cd frontend
npm i --force
```

3) Устанавливаем wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

и экспортируем чтобы получить доступ к исполняемому файлу (пример для Linux)

```bash
export PATH=$PATH:/home/user/go/bin
```

4) Собираем проект

```bash
wails build
```

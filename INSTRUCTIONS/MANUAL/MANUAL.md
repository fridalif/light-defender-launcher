# Инструкция по запуску Light Defender вручную на сервере

***При возникновении вопросов пишите автору в telegram: @frdlf, или на почту: contact@light-defender.ru***

Для того чтобы установить Light Defender на свой сервер нужно выполнить 6 простых шагов:
1) Получить конфигурацию у бота @light_defender_bot или создать её в dashboard.light-defender.ru;
2) Подключиться к серверу по SSH;
3) Скачать подходящую версию клиента

```bash
wget -q -O - https://dashboard.light-defender.ru/load_ld.sh | sudo bash
```

4) Перейти в скачанную директорию

```bash
cd Название_выводится_в_конце_предыдущего_скрипта
```

5) Загрузить конфигурацию с сервера (или загрузить её вручную без использования программы)

```bash
sudo ./ldclient.bin --load-config
```

6) Установить Light Defender

```bash
sudo ./ldclient.bin --install
```

Для того чтобы обновить версию Light Defender до новой можно воспользоваться лаунчером или ввести следующие команды в директории с Light Defender:

```bash
sudo systemctl stop ldclient
sudo mv ldclient.bin ldclient.bin.save
sudo ./ldclient.bin.save -uni && sudo rm -rf ldclient.bin.save
sudo systemctl start ldclient
```

Если хотите обновиться не до самой новой версии или возникли проблемы с определением архитектуры можно обновить в интерактивном режиме

```bash
sudo systemctl stop ldclient
sudo mv ldclient.bin ldclient.bin.save
sudo ./ldclient.bin.save -update && sudo rm -rf ldclient.bin.save
sudo systemctl start ldclient
```
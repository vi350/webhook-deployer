# Webhook deployer

Этот проект создавался для запуска: будучи субмодулем, в докере,
используя докер компоуз, из корневой директории проекта. Для добавления
как субмодуля:
    
```bash
git submodule add https://github.com/vi350/webhook-deployer
```

Для обновления:

```bash
git submodule init
git submodule update --remote
git submodule foreach git pull origin master
```

Для запуска нужно создать контейнер и примаунтить к нему:

```text
/var/run/docker.sock:/var/run/docker.sock
./:/repo/
~/.ssh/:/root/.ssh/
```
.env будет скопирован в образ из докерфайла
контекст билда должен соответствовать папке субмодуля

### архитектура

### middleware авторизации запроса

## можно добавить:

### healthcheck каждого сервиса/всех сразу

не очень понятно зачем, требует рассмотрения

### вместо своей реализации подтянуть ее из библиотеки гитхаба

так как там реализован парсинг всех типов ивентов, но хотелось разобараться
самому, поэтому пока так

### критичное: перейти с down на stop + remove
для того чтобы не пересоздавать контейнер в котором находится деплоер


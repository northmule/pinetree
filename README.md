# Pinetree

Обновляет описание для всех видео в указанной группе Вконтакте.

![Alt Text](https://github.com/northmule/pinetree/blob/master/doc/pin.gif)

## Настройка client.yaml
```yaml
VK:
   ApiVersion: "5.199" # версия API которую использовать
   AccessToken: "" # токен доступа полученный от ВК - подробнее https://dev.vk.com/ru/api/access-token/getting-started
   GroupID: "" # Идентификатор пользователя или сообщества, которому принадлежат видеозаписи. Идентификатор сообщества должен начинаться со знака -
   AlbumID: "1" # Идентификатор альбома

Log:
   FilePath: "app.log" # путь к файлу.
   Level: "info" # уроверь логирования
```

Файл client.yaml должен располагаться рядом с запускаемым приложением. Собранные приложения находятся в папке /build


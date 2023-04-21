### Модуль для работы с плейлистом
Модуль находится в каталоге `internal/model/playlist`. <br> Для запуска тестов модуля можно воспользоваться командой `make test_model`


Модуль реализован на основе двусвязного списка и обладает следующими методами: <br>
<ul>
<li>Play - начинает воспроизведение</li>
<li>Pause - приостанавливает воспроизведение</li>
<li>AddSong - добавляет в конец плейлиста песню</li>
<li>Next воспроизвести след песню</li>
<li>Prev воспроизвести предыдущую песню</li>
</ul>
 Воспроизведение песен эмулируется длительной операцией.

### Сервис для управления музыкальным плейлистом
Доступ к сервису осуществляется с помощью API, который имеет возможность выполнять CRUD операции с песнями в плейлисте, а также воспроизводить, приостанавливать, переходить к следующему и предыдущему трекам. Для хранения песен используется PostgreSQL. В качестве протокола взаимодействия используется gRPC. 

Запуск тестов:<br>
`make compose_database`<br>
`make migrate_up`<br>
`make test_server` <br>
`make remove_database`


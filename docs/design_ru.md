<!-- markdownlint-disable MD033 -->

# Проектирование системы EduFlow

---

## 1. Концепция и решаемые задачи

**EduFlow** представляет собой информационную систему управления образовательным процессом, спроектированную для автоматизации ведения классного журнала, учета посещаемости и мониторинга академической успеваемости.

Основная цель проекта — минимизация бюрократической нагрузки на преподавательский состав и обеспечение прозрачного, оперативного доступа к учебным данным для администрации и учащихся.

## 2. Архитектурное решение

Система базируется на **микросервисной архитектуре**, что обеспечивает независимое масштабирование модулей и высокую отказоустойчивость. Взаимодействие компонентов организовано через единую точку входа — **API Gateway**, выполняющий функции обратного проксирования и агрегации данных в рамках паттерна **BFF (Backend for Frontend)**.

### Ключевые технологические стеки

* **Транспортный уровень:** Гибридная модель взаимодействия. Внешний контур использует **REST (JSON)** для обеспечения совместимости с клиентскими приложениями (Flutter, Web). Внутренний контур построен на высокопроизводительном протоколе **gRPC (Protobuf)** с кодогенерации через **protoc** для go lang, обеспечивающем строгую типизацию и бинарную сериализацию данных.
* **Слой бизнес-логики:** Микросервисы реализованы на языке **Go (Golang)**. Выбор обусловлен эффективной моделью конкурентности, малым потреблением ресурсов и высокой скоростью обработки сетевых запросов.
* **Слой данных:** Реализован паттерн **Database-per-Service**. Каждая функциональная область обладает собственной изолированной БД (**PostgreSQL**) и выделенным слоем кэширования (**Redis**). Для управления уникальными идентификаторами в распределенной среде используется стандарт **UUID v7**.
* **Безопасность и контроль доступа:** Централизованный сервис авторизации интегрирован с **SpiceDB** (реализация модели Google Zanzibar). Это позволяет внедрить графовую систему прав доступа (**ReBAC**), обеспечивая гибкое управление полномочиями (учитель, староста, администратор) на уровне конкретных объектов системы.

---

```mermaid
graph TD

    subgraph Clients
        direction LR
        Flutter -- "Rest" <--> GW
        Web -- "Rest" <--> GW
    end

    GW[API Gateway / Aggregator]

    subgraph Cluster [System Core]
        direction TB
        
        subgraph S_Auth [Auth Service]
            direction LR
            Auth_L[Go Logic] --- Auth_D[("Postgres + SpiceDB<br/>Redis (Sessions)")]
        end

        subgraph S_Jour [Journal Service]
            direction LR
            Jour_L[Go Logic] --- Jour_D[("Postgres (Grades)<br/>Redis (Cache)")]
        end

    end

    %% Основной трафик
    GW == "gRPC: get" ==> S_Auth
    GW == "gRPC: get" ==> S_Jour

    %% Межсервисные связи
    S_Jour -. "gRPC: Check" .-> S_Auth

    style Cluster fill:#f9f9f9,stroke:#333,stroke-dasharray: 5 5
    style Auth_D fill:#e3f2fd,stroke:#2196f3
    style Jour_D fill:#fff3e0,stroke:#ff9800

```

## 3. API Gateway (KrakenD)

В качестве единой точки входа и агрегатора данных в системе используется высокопроизводительный L7-шлюз **KrakenD**. Его основная задача — реализация паттерна **BFF (Backend for Frontend)**, что позволяет абстрагировать сложность микросервисной архитектуры от клиентских приложений (Flutter, Web).

### Основные функции и механизмы

* **Трансформация протоколов (gRPC-JSON Transcoding):** KrakenD принимает внешние HTTP/REST запросы в формате JSON и преобразует их во внутренние gRPC-вызовы к сервисам `Auth` и `Journal`. Это позволяет использовать преимущества бинарного протокола внутри кластера, сохраняя простоту интеграции для фронтенд-разработки.

* **Агрегация запросов (Endpoints Aggregation):** Шлюз позволяет минимизировать количество сетевых задержек (RTT) для мобильных устройств. Например, при запросе «страницы журнала» KrakenD выполняет параллельные запросы к сервисам для получения оценок, данных учеников и проверки прав доступа, формируя единый консолидированный ответ.

* **Безопасность и Валидация:**
  * **JWT Validation:** Проверка подписи и срока действия токена на уровне шлюза до проброса запроса в бизнес-логику.
  * **Header Manipulation:** Извлечение `user_id` из токена и его инъекция в gRPC-метаданные для последующей идентификации пользователя внутри системы.
* **Оптимизация трафика:**
  * **Throttling & Rate Limiting:** Защита микросервисов от лавинообразных запросов.
  * **Response Manipulation:** Удаление избыточных полей из ответов микросервисов для экономии мобильного трафика пользователей.

### Логика работы на примере запроса «Журнал класса»

1. **Клиент:** Отправляет `GET /v1/journal/{class_id}` с Bearer-токеном.
2. **KrakenD:**
   * Проверяет валидность JWT.
   * Инициирует gRPC-вызов в **Auth Service** для подтверждения прав субъекта на доступ к ресурсу `class_id` (через SpiceDB).
   * Параллельно запрашивает данные об успеваемости и профилях учащихся из **Journal Service**.
3. **Агрегатор:** Сопоставляет полученные данные и возвращает фронтенду структурированный JSON, готовый к отрисовке без дополнительной обработки на клиенте.

## Journal service

Journal Service является центральным хранилищем и поставщиком всех данных, связанных с организацией и ведением образовательного процесса. Сервис управляет:

* учебными классами и группами;
* списком предметов и учебных планов;
* привязкой учителей к предметам и классам;
* составом учащихся по классам;
* успеваемостью (оценки, средние баллы);
* посещаемостью занятий;
* расписанием и домашними заданиями (при необходимости).

### Endpoints: Journal Microservice

#### 1. Students (Студенты)

| Method | Request | Response | Description |
| :--- | :--- | :--- | :--- |
| **CreateStudent** | CreateStudentRequest | Student | Добавление нового студента |
| **GetStudent** | GetStudentRequest | Student | Получение данных студента по ID |
| **ListStudents** | ListStudentsRequest | ListStudentsResponse | Список студентов с фильтрацией |
| **UpdateStudent** | UpdateStudentRequest | Student | Обновление информации о студенте |
| **DeleteStudent** | DeleteStudentRequest | Empty | Удаление студента из системы |

#### 2. Teachers (Учителя)

| Method | Request | Response | Description |
| :--- | :--- | :--- | :--- |
| **CreateTeacher** | CreateTeacherRequest | Teacher | Добавление нового преподавателя |
| **ListTeachers** | ListTeachersRequest | ListTeachersResponse | Получение списка всех учителей |
| **UpdateTeacher** | UpdateTeacherRequest | Teacher | Обновление данных преподавателя |
| **DeleteTeacher** | DeleteTeacherRequest | Empty | Удаление преподавателя |

#### 3. Classes & Subjects (Классы и Предметы)

| Method | Request | Response | Description |
| :--- | :--- | :--- | :--- |
| **CreateClass** | CreateClassRequest | Class | Создание нового учебного класса |
| **GetClass** | GetClassRequest | Class | Получение данных класса по ID |
| **ListTeacherClasses** | ListTeacherClassesRequest | ListClassesResponse | Список классов конкретного учителя |
| **UpdateClass** | UpdateClassRequest | Class | Изменение данных класса |
| **DeleteClass** | DeleteClassRequest | Empty | Удаление класса |
| **CreateSubject** | CreateSubjectRequest | Subject | Добавление учебного предмета |
| **UpdateSubject** | UpdateSubjectRequest | Subject | Изменение названия или данных предмета |
| **DeleteSubject** | DeleteSubjectRequest | Empty | Удаление предмета |

#### 4. Grades (Оценки)

| Method | Request | Response | Description |
| :--- | :--- | :--- | :--- |
| **RecordGrade** | RecordGradeRequest | Grade | Выставить оценку студенту |
| **ListGrades** | ListGradesRequest | ListGradesResponse | Ведомость оценок (фильтр по классу/предмету) |
| **UpdateGrade** | UpdateGradeRequest | Grade | Изменение выставленной оценки |
| **DeleteGrade** | DeleteGradeRequest | Empty | Удаление записи об оценке |

#### 5. Homework (Домашняя работа)

| Method | Request | Response | Description |
| :--- | :--- | :--- | :--- |
| **RecordHomework** | RecordHomeworkRequest | Homework | Написать домашнее задание |
| **UpdateHomework** | UpdateHomeworkRequest | Homework | Обновить данные домашнего задания |
| **ListHomework** | ListHomeworkRequest | ListHomeworkResponse | Получение списка домашних заданий (Метод сильно зависит от того кто делает запрос) |
| **DeleteHomework** | DeleteHomeworkRequest | Empty | Удаление домашнего задания |

#### 6. Status code (Статус коды)

| Method | Request | Response | Description |
| :--- | :--- | :--- | :--- |
| **CreateStatusCode** | CreateStatusCodeRequest | StatusCode | Создать статус-код |
| **UpdateStatusCode** | UpdateStatusCodeRequest | StatusCode | Обновить статус-код |
| **ListStatusCode** | ListStatusCodeRequest | ListStatusCodeResponse | Получение список статус-кодов |
| **DeleteStatusCode** | DeleteStatusCodeRequest | Empty | Удаление статус-кода |

#### 7. Teaching Load (Учебная нагрузка)

| Method | Request | Response | Description |
| :--- | :--- | :--- | :--- |
| **CreateTeachingLoad** | CreateTeachingLoadRequest | TeachingLoad | Назначить учителя на предмет в конкретном классе |
| **UpdateTeachingLoad** | UpdateTeachingLoadRequest | TeachingLoad | Изменить параметры нагрузки (смена учителя) |
| **ListTeachingLoad** | ListTeachingLoadRequest | ListTeachingLoadResponse | Получить список всех распределений нагрузки |
| **GetTeachingLoad** | GetTeachingLoadRequest | TeachingLoad | Получить детали конкретного назначения |
| **DeleteTeachingLoad** | DeleteTeachingLoadRequest | Empty | Удалить запись о нагрузке |

<details>
<summary>Структуры (proto)</summary>

[embedmd]:# (../api/proto/journal/v1/journal.proto protobuf)

</details>

### Модель данных

Все сущности хранятся в изолированной PostgreSQL, а для повышения производительности используется Redis

```mermaid
erDiagram
    %% Core Entities
    "Class" ||--o{ TeachingLoad  : assigned
    "Class" ||--o{ Student : studies
    "Class" ||--o{ Homework : studies

    Teacher ||--o{ TeachingLoad  : teaches
    Teacher ||--o{ Homework : teaches
    Subject ||--o{ TeachingLoad  : defines
    
    %% Student Interactions (Dotted for clarity)
    TeachingLoad  ||--o{ Grade : results_in
    StatusCode ||--o{ Grade: gets
    Student ||--o{ Grade : gets

    "Class" {
        uuid id PK
        string class_name
        int year_of_study
        int graduation_year
    }

    Student {
        uuid id PK
        uuid class_id FK
        string full_name
    }

    Teacher {
        uuid id PK
        string full_name
    }

    Subject {
        uuid id PK
        string name
    }

    TeachingLoad  {
        uuid id PK
        uuid teacher_id FK
        uuid subject_id FK
        uuid class_id FK
    }

    Grade {
        uuid id PK
        uuid student_id FK
        uuid tl_id FK
        uuid status_code_id FK
        int score
        int lesson_number
        DATE lesson_date
        string note
    }

    StatusCode {
        uuid id
        VARCHAR(6) abbreviation
    }

    Homework {
        uuid id PK
        uuid teacher_id FK
        uuid class_id FK
        string description_task
        DATE deadline_at
        DATE assigned_at
    }
```

#### Описание статус кодов

| Code | description |
| :--- | :--- |
| DEBT | Задолженность. Работа не сдана в срок («точка» в журнале) |
| ABSENT | Отсутствие. Ученика не было на уроке без объяснения причины |
| SICK | Болезнь. Отсутствие по медицинской справке |
| EXCUSE | Уважительная. Семейные обстоятельства, заявление от родителей |
| RETAKE | Пересдача. Оценка, полученная при исправлении предыдущей |
| EXEMPT | Освобожден. Например, освобождение от физкультуры или зачета |

<details>
<summary>SQL код</summary>

[embedmd]:# (../migrations/000001_init.up.sql)

</details>

## Auth service

...

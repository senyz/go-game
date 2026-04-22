package handler

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/senyz/go-game/interfaces"
)

type WebhookHandler struct {
	gameService interfaces.GameService
}

func NewWebhookHandler(gameService interfaces.GameService) *WebhookHandler {
	return &WebhookHandler{gameService: gameService}
}

func (h *WebhookHandler) HandleMessage(c *gin.Context) {
	var request struct {
		UserID  string `json:"user_id"`
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// Парсим ID пользователя
	userID, err := parseUserID(request.UserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Получаем текущий прогресс пользователя, чтобы определить текущую сцену
	currentSceneID, err := h.gameService.GetCurrentSceneID(userID)
	if err != nil {
		log.Printf("Failed to get current scene for user %d: %v", userID, err)
		c.JSON(500, gin.H{"error": "Failed to retrieve game progress"})
		return
	}

	// Обрабатываем ответ пользователя
	nextScene, err := h.gameService.ProcessAnswer(userID, currentSceneID, request.Message)
	if err != nil {
		log.Printf("Game processing error for user %d: %v", userID, err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Формируем ответ в зависимости от типа сцены
	response := gin.H{
		"scene_id": nextScene.ID,
		"title":    nextScene.Title,
		"text":     nextScene.Description,
	}

	if nextScene.Question != "" {
		response["question"] = nextScene.Question
	} else {
		// Это финальная сцена — добавляем сообщение о завершении
		response["message"] = "Поздравляем! Вы завершили приключение!"
	}

	c.JSON(200, response)
}

// parseUserID преобразует строку в uint, возвращает ошибку при неудаче
func parseUserID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

/*
Передача токена через query-параметры больше не поддерживается — используйте заголовок Authorization: <token>

API (Application Programming Interface) — это посредник между разработчиком приложений и средой, с которой это приложение должно взаимодействовать. API упрощает написание кода за счёт набора готовых классов, функций или структур для работы с данными

API MAX — это интерфейс, который позволяет ботам взаимодействовать с платформой и получать необходимые данные с помощью HTTPS-запросов к серверу. В этом разделе расскажем, как подготовиться к использованию API приложения

Методы
HTTPS-запросы на домен platform-api.max.ru вызывают методы — условные команды, которые соответствуют той или иной операции с базой данных. Например, получение, запись или удаление какой-либо информации
Параметры запроса должны содержать HTTP-метод, соответствующий необходимой операции:

GET — получить ресурсы
POST — создать ресурсы (например, отправить новые сообщения)
PUT — редактировать ресурсы
DELETE — удалить ресурсы
PATCH — исправить ресурсы
В зависимости от конкретного метода, параметры запроса будут отображаться в пути, URL-параметрах или теле запроса

Примеры запросов:

GET https://platform-api.max.ru/messages/{messageId} — получить сообщения
POST https://platform-api.max.ru/messages — отправить сообщения
PATCH https://platform-api.max.ru/chats/{chatId} — изменить информацию о чате
В ответ сервер вернёт JSON-объект с запрошенными данными или сообщение об ошибке, если что-то пойдёт не так
JSON — это формат записи данных в виде пар <ИМЯ_СВОЙСТВА>: <ЗНАЧЕНИЕ>. Прочитайте об особенностях формата JSON, если вы ещё не работали с ним
Пример ответа на запрос к методу GET /me:

JSON
{
	"user_id": 1,
	"name": "My Bot",
	"username": "my_bot",
	"is_bot": true,
	"last_activity_time": 1737500130100
}
Также, помимо JSON, сервер вернет трёхзначный HTTP-код, информирующий об успешном выполнении запроса или ошибке.

Коды ответов HTTP
200 — успешная операция
400 — недействительный запрос
401 — ошибка аутентификации
404 — ресурс не найден
405 — метод не допускается
429 — превышено количество запросов
503 — сервис недоступен
Рекомендации по работе с API
Когда вы настраиваете получение обновлений о действиях в чат-боте, используйте:

Long Polling — для разработки и тестирования
только Webhook — для production-окружения
Для стабильной работы сервисов MAX убедитесь, что максимальное количество запросов в секунду на platform-api.max.ru — 30 rps

Клавиатура для чат-бота
Клавиатура позволяет отправлять боту запросы кнопками, а не сообщениями. Чтобы клавиатура была удобной для пользователей, рекомендуем заранее продумать её наполнение и учитывать обязательные параметры:

Текст на кнопке выравнивается по центру и обрезается, если выходит за её границы
Кнопки в одной строке всегда одинаковой ширины
Ширина каждого ряда кнопок равна ширине клавиатуры
Высота у всех кнопок по умолчанию одинаковая
Вы можете подключить к чат-боту в MAX inline-клавиатуру. Она позволяет разместить под сообщением бота до 210 кнопок, сгруппированных в 30 рядов — до 7 кнопок в каждом (до 3, если это кнопки типа link, open_app, request_geo_location или request_contact)

Для кнопки с видом link максимальный размер ссылки составляет 2048 символов

Типы кнопок
callback — сервер MAX отправляет событие с типом message_callback (через Webhook или Long polling)
Обратите внимание: для отправки вебхуков поддерживается только протокол HTTPS, включая самоподписанные сертификаты. HTTP не поддерживается

link — позволяет открыть ссылку в новой вкладке
request_contact — запрашивает у пользователя его контакт и номер телефона
request_geo_location — запрашивает у пользователя его местоположение
open_app — открывает мини-приложение
message — отправляет боту текстовое сообщение
clipboard — копирует текст, указанный в свойстве payload, в буфер обмена
Кнопка clipboard
При нажатии на кнопку с типом clipboard текст, указанный в свойстве payload, копируется в буфер обмена

В свойстве payload можно передать любой текст, например промокод, трек-номер, платёжные реквизиты

JSON
{
  "type": "clipboard", // Тип кнопки
  "text": "Скопировать", // Текст кнопки
  "payload": "123456" // Текст, который будет скопирован
}
Как добавить кнопки
Чтобы добавить кнопки, отправьте сообщение с InlineKeyboardAttachment

JSON
{
  "text": "It is message with inline keyboard",
  "attachments": [
    {
      "type": "inline_keyboard",
      "payload": {
        "buttons": [
          [
            {
              "type": "callback", // Тип кнопки
              "text": "Press me!", // Текст кнопки
              "payload": "button1 pressed" // Описание действия
            }
          ]
        ]
      }
    }
  ]
}
Форматирование текста
Текст сообщения в чат-боте можно улучшить с помощью базового форматирования. Для этого вы можете использовать либо Markdown, либо HTML

Markdown
Чтобы включить разбор Markdown, установите свойство format в NewMessageBody на значение markdown

Markdown	Отображение
курсив	*empasized* или _empasized_
жирный	**strong** или __strong__
зачёркнутый	~~strikethough~~
подчёркнутый	++underline++
моноширинный	`code` (переводы строк внутри этого блока обрабатываются как пробелы)
ссылка	[Inline URL](https://dev.max.ru/)
@упоминание пользователя	"text": "[Имя Фамилия](max://user/user_id)", "format": "markdown"
Вместо User mention указывайте полное имя пользователя из профиля в MAX, в том числе фамилию. Если фамилия отсутствует — только имя
HTML
Чтобы включить разбор HTML, установите свойство format в NewMessageBody на значение html

Markdown	Отображение
курсив	<i> или <em>
жирный	<b> или <strong>
зачёркнутый	<del> или <s>
подчёркнутый	<ins> или <u>
моноширинный	<pre> или <code>
ссылка	<a href="https://dev.max.ru">Docs</a>
@упоминание пользователя	"text": "<a href=\"max://user/user_id\">Имя Фамилия</a>", "format": "html"
Вместо User mention указывайте полное имя пользователя из профиля в MAX, в том числе фамилию. Если фамилия отсутствует — только имя
*/

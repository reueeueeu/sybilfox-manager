<img src="https://raw.githubusercontent.com/reueeueeu/sybilfox-manager/refs/heads/main/banner.png" align="center">

<h1 align="center">Sybilfox</h1>

<h4 align="center">Антидект браузер и GUI менеджер для управления множеством профилей</h4>

<p align="center">
  <img src="https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20MacOS-green" alt="platforms" style="width: 250px; height: 25px;" width="250" height="55"/></a>
</p>



---

## Описание

**Все анонсы в телеграм канале [Sybilfox](https://t.me/sybilfox)**

Описание сгенерировано с помощью ИИ.

**Sybilfox Manager** — это графический менеджер для работы с множественными профилями антидетект браузеров с расширенными возможностями по инъекции и ротации отпечатков браузера. 

**Sybilfox** - Антидетект браузер. Разработанный так, чтобы оставаться незаметным для всех анти-бот систем, Sybilfox превосходит большинство коммерческих анти-бот браузеров, предоставляя надёжный контроль над ротацией отпечатков — и всё это без использования JS-инъекций.

---

<img width="1657" alt="screenshot" align="center" src="https://github.com/user-attachments/assets/e70523b8-ca90-4ad9-928f-3e495b8522e5" />

---

## Скачать
**[Windows](https://github.com/reueeueeu/sybilfox-manager/releases/download/0.1/sybilfox-manager_win_x86_64.exe)**
**[Linux](https://github.com/reueeueeu/sybilfox-manager/releases)**
**[Mac](https://github.com/reueeueeu/sybilfox-manager/releases/download/0.1/SybilfoxManager.dmg)**

---




## Особенности

- **Управление профилями браузеров:** Создание и управление множеством профилей антидетект браузеров.
- **Инъекция и ротация отпечатков:**  
  - Инъекция и ротация всех свойств `navigator` (устройство, ОС, аппаратное обеспечение, браузер и т.д.).
  - Подмена геолокации, часового пояса и локали.
  - Подмена шрифтов и защита от фингерпринтинга.
  - Подмена параметров WebGL, поддерживаемых расширений, атрибутов контекста и форматов точности шейдеров.
  - Подмена IP-адресов WebRTC на уровне протокола.
  - Подмена медиа-устройств, голосов и скорости воспроизведения речи.
- **Поддержка прокси:** Возможность назначения различных прокси для каждого профиля.
- **Автоматическое переключение профилей:** Лёгкое переключение между сохранёнными конфигурациями профилей.
- **Сохранение сессий:** Сохранение и загрузка сессий браузер.
- **Кроссплатформенность:** Работает на Windows, Linux и macOS.
- **Удобный графический интерфейс:** Интуитивно понятное управление профилями.

---

## Расширенные возможности

### Незаметность для анти-бот систем
- **Стелс-производительность:** Sybilfox разработан так, чтобы быть незаметным для всех анти-бот систем, превосходя большинство коммерческих решений.

### Инъекция и ротация отпечатков (без JS-инъекций!)
- **Полная Подмена свойств navigator:**  
  - Все свойства `navigator` (устройство, ОС, аппаратное обеспечение, браузер и т.д.) ✅
- **Параметры отображения:**  
  - Подмена размеров экрана, разрешения, окна и viewport ✅
- **Геолокация и локаль:**  
  - Подмена геолокации, часового пояса и локали ✅
- **Шрифты и анти-фингерпринтинг:**  
  - Подмена шрифтов и защита от фингерпринтинга ✅
- **Возможности WebGL:**  
  - Подмена параметров WebGL, поддерживаемых расширений, атрибутов контекста и форматов точности шейдеров ✅
- **WebRTC:**  
  - Подмена IP-адресов WebRTC на уровне протокола ✅
- **Медиа-устройства:**  
  - Подмена медиа-устройств, голосов и скорости воспроизведения речи. Battery API ✅

### Стелс-патчи
- **Предотвращение утечек при выполнении кода:** Изоляция всего JavaScript-кода страницы в безопасном окружении.
- **Изоляция фреймов:** Предотвращение утечек данных в контексте выполнения фреймов.
- **Удаление потенциально утечивающих патчей:** Исключение патчей, связанных с анти-zoom и обработкой meta viewport.
- **Изоляция контента:** Повторное включение механизмов изоляции контента.
- **Обработка PDF:** Повторное включение PDF.js для улучшённого отображения документов.

### Анти-фингерпринтинг шрифтов
- **Оптимизация системных шрифтов:** Автоматический выбор корректных системных шрифтов в зависимости от User Agent.
- **Встроенные шрифты:** Поставляется с системными шрифтами для Windows, Mac и Linux.
- **Случайное смещение метрик:** Защита от фингерпринтинга шрифтов путём случайного смещения межбуквенных интервалов.

---

## Тестирование

Sybilfox успешно проходит тестирование по всем основным системам защиты (исходные тесты с Botright):
| Test                                                                                               | Status                                                    |
| -------------------------------------------------------------------------------------------------- | --------------------------------------------------------- |
| [**CreepJS**](https://abrahamjuliot.github.io/creepjs/)                                            | ✔️ 71.5%.         |
| [**Rebrowser Bot Detector**](https://bot-detector.rebrowser.net/)                                  | ✔️                                        |
| [**BrowserScan**](https://browserscan.net/)                                                        | ✔️ 100%. |
| **reCaptcha Score**                                                                                | ✔️                                                        |
| ‣ [nopecha.com](https://nopecha.com/demo/recaptcha)                                                | ✔️                                                        |
| ‣ [recaptcha-demo.appspot.com](https://recaptcha-demo.appspot.com/recaptcha-v3-request-scores.php) | ✔️ 0.9                                                    |
| ‣ [berstend.github.io](https://berstend.github.io/static/recaptcha/v3-programmatic.html)           | ✔️ 0.9                                                    |
| **DataDome**                                                                                       | ✔️                                                        |
| ‣ [DataDome bot bounty](https://yeswehack.com/programs/datadome-bot-bounty#program-description)    | ✔️ Все тесты пройдены.                                   |
| ‣ [hermes.com](https://www.hermes.com/us/en/)                                                      | ✔️                                                        |
| **Imperva**                                                                                        | ✔️                                                        |
| ‣ [ticketmaster.es](https://www.ticketmaster.es/)                                                  | ✔️                                                        |
| **Cloudflare**                                                                                     | ✔️                                                        |
| ‣ [Turnstile](https://nopecha.com/demo/turnstile)                                                  | ✔️                                                        |
| ‣ [Interstitial](https://nopecha.com/demo/cloudflare)                                              | ✔️                                                        |
| **WebRTC IP Spoofing**                                                                             | ✔️                                                        |
| ‣ [Browserleaks WebRTC](https://browserleaks.net/webrtc)                                           | ✔️                         |
| ‣ [CreepJS WebRTC](https://abrahamjuliot.github.io/creepjs/)                                       | ✔️ Host & STUN IP correctly.                       |
| ‣ [BrowserScan WebRTC](https://www.browserscan.net/webrtc)                                         | ✔️  Host & STUN IP correctly.                       |
| **Font Fingerprinting**                                                                            | ✔️                                                        |
| ‣ [Browserleaks Fonts](https://browserleaks.net/fonts)                                             | ✔️                                    |
| ‣ [CreepJS TextMetrics](https://abrahamjuliot.github.io/creepjs/tests/fonts.html)                  | ✔️                                    |
| [**Incolumitas**](https://bot.incolumitas.com/)                                                    | ✔️ 0.8-1.0                                                |
| [**SannySoft**](https://bot.sannysoft.com/)                                                        | ✔️                                                        |
| [**Fingerprint.com**](https://fingerprint.com/products/bot-detection/)                             | ✔️                                                        |
| [**IpHey**](https://iphey.com/)                                                                    | ✔️                                                        |
| [**Bet365**](https://www.bet365.com/#/AC/B1/C1/D1002/E79147586/G40/)                               | ✔️                                                        |

---

## Дорожная карта

- [ ] Создать антидетекст браузер на хроме.
- [ ] Реализовать облачное хранение профилей.
- [ ] Интегрировать ИИ для рандомизации отпечатков.

---

## Установка

### Требования
- **[Wails](https://wails.io/docs/gettingstarted/installation)**


### Шаги установки

```sh
# Клонируем репозиторий
git clone https://github.com/reueeueeu/sybilfox-manager.git
cd sybilfox-manager

# Билд
wails build

```

---

## Вклад в проект

Мы приветствуем вклады в проект! Чтобы начать:

1. Форкните репозиторий.
2. Создайте новую ветку.
3. Сделайте коммиты и отправьте изменения.
4. Откройте pull request.

---

## Контакты

Для вопросов или поддержки, пожалуйста, свяжитесь:

- **Telegram:** [@reueeueeu](https://t.me/reueeueeu)

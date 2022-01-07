# Go Blog Updates - Scanner
This service scans go blog ([go.dev](https://go.dev)) and publishes new posts to message broker ([rabbitmq](https://www.rabbitmq.com/)).
It uses [mongodb](https://www.mongodb.com/) as a storage for already published posts.

### Consumers:
 - [gbu-telegram-bot](https://github.com/don2quixote/gbu-telegram-bot) service
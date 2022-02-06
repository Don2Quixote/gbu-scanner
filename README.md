# Go Blog Updates - Scanner
This service scans go blog ([go.dev](https://go.dev)) and publishes new posts to message broker ([rabbitmq](https://www.rabbitmq.com/)).
It uses [mongodb](https://www.mongodb.com/) as a storage for already published posts.

## ENV Configuration:
| name                      | type   | description                                                                        |
| ------------------------- | ------ | ---------------------------------------------------------------------------------- |
| BLOG_HOST                 | string | Host where blog is. You definitely want to set it to "go.dev"                      |
| BLOG_PATH                 | string | Path to find all posts. You definitely want to set it to "/blog/all"               |
| BLOG_HTTPS                | string | Flag to use https protocol instead of http. You definitely want to set it to true  |
| BLOG_SCAN_INTERVAL        | int    | Duration between scan's interations (seconds)                                      |
| BLOG_SCAN_NETWORK_TIMEOUT | int    | Duration after which timeout error will happen during getting posts (seconds)      |
| MONGO_HOST                | string | Database host                                                                      |
| MONGO_USER                | string | Database user                                                                      |
| MONGO_PASS                | string | Database password                                                                  |
| MONGO_SRV                 | bool   | Flag to use mongodb+srv protocol instead of mongodb                                |
| MONGO_DATABASE            | string | Database name                                                                      |
| RABBIT_HOST               | string | Rabbit host                                                                        |
| RABBIT_USER               | string | Rabbit user                                                                        |
| RABBIT_PASS               | string | Rabbit password                                                                    |
| RABBIT_VHOST              | string | Rabbit vhost                                                                       |
| RABBIT_AMQPS              | bool   | Flag to use amqps protocol instead of amqp                                         |
| RABBIT_RECONNECT_DELAY    | int    | Delay (seconds) before attempting to reconnect to rabbit after loosing connection  |

Env template for sourcing is [deployments/local.env](deployments/local.env)
```
$ source deployments/local.env
```

## MongoDB schema
**publishedPosts collection**
```
{
    posts: [{
        title: string,
        date: ISODate,
        author string,
        summary: string,
        url: string
    }]
}
```
Почему one document with posts array instead of one document per one post?

Потому что mongodb's atomic transactions are only available with single document (mongodb is nice choice XDDDDDD) <!-- или я просто ничего не понял -->

> [this single-document atomicity obviates the need for multi-document transactions for many practical use cases](https://docs.mongodb.com/manual/core/transactions/#transactions)

Also because of mongodb's transactions usage it's [impossible to use standalone instance](https://docs.mongodb.com/manual/core/transactions/#feature-compatibility-version--fcv-) XD

## Makefile commands:
| name | description                                                                            |
| ---- | -------------------------------------------------------------------------------------- |
| lint | Runs linters                                                                           |
| test | Runs tests, but there are no tests                                                     |
| run  | Sources env variables from [deployments/local.env](deployments/local.env) and runs app |
| stat | Prints stats information about project (packages, files, lines, chars count)           |

Direcotry [scripts](/scripts) contains scripts which invoked from [Makefile](Makefile)

## Consumers:
 - [gbu-telegram-bot](https://github.com/don2quixote/gbu-telegram-bot) service
 - [gbu-queue-api](https://github.com/don2quixote/gbu-queue-api) service
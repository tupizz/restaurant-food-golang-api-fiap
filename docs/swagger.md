### Add dependecy in local
O Swag é uma ferramenta que analisa seus comentários no código para gerar arquivos de documentação Swagger.

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Observação: Certifique-se de que o diretório $GOPATH/bin está no seu PATH para que você possa usar o comando swag.

### Instalar as bibliotecas para o Gin

```bash
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

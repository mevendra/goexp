## Lab CLoud APP

Aplicação para consulta de temperatura

###  Como rodar/testar a aplicação?
* Serviço disponível através da porta 8080:
  * Através do Dockerfile;
  * Através do arquivo `cmd/temperature/main.go`
* Serivço disponível através do Google Run no endereço: https://temperature-yffengdbxq-vp.a.run.app
* Cep pode ser consultado atraves de query parameter cep: `?cep=99500000`
* Exemplo `.http` de consulta disponível em: `api/get_temperature.http` 
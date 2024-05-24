## Lab Open Telemetry

### Como rodar a aplicação?
* Executar comando `docker-compose up` para iniciar a aplicação
  * Server A
  * Server B
  * Zipkin
* Exemplo de requisições estão na pasta `/api`
  * Requisição GET `/temperature?cep=99500000`
  * Requisição POST `/validate`
* Tracing aparecerá no zipkin, na porta `9411 (http://localhost:9411/zipkin/)`
  * Requisição `temperature` gera um tracing standalone;
  * Requisição `validate` gera um tracing distribuido, entre os dois sistemas;

### Variáveis de ambiente para server
* PORT: Porta de acesso ao sistema web;
* TEMPERATURE_URI: Url de acesso ao serviço de temperatura;
* NAME: Nome do serviço;
* ZIPKIN_URL: Url de envio ao serviço zipkin;

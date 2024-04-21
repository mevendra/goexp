## Rate Limiter APP

Aplicação para desafio de sistema para Rate Limit;

Aplicação disponibiliza um RateLimiter, um middleware HTTP utilizando o RateLimiter e exemplo de uso em um servidor utilizando Redis para controle de rota padrão `/`.

A aplicação permite realizar controle por Token ou IP, dando prioridade ao controle por Token;

O exemplo possui arquivo `docker-compose.yaml` para subir redis, e as seguintes variáveis podem ser definidas no arquivo `cmd/limiter/.env`:
* TOKEN_LIMIT: Número de requisições para tokens disponível em determinado time-frame. Para desativar controle de rate limit, remover ou deixar com valor 0;
  * Ex: `10`;
* TOKEN_BLOCK_TIME: Tempo em que um token irá ficar bloqueado após atingir seu rate-limit;
  * Ex: `10m`;
* TOKEN_FRAME_TIME: Time-frame em que as requisições serão acumuladas;
  * Ex: `1s`;
* IP_LIMIT: Número de requisições para IP disponível em determinado time-frame. Para desativar controle de rate limit, remover ou deixar com valor 0;
  * Ex: `10`;
* IP_BLOCK_TIME: Tempo em que um IP irá ficar bloqueado após atingir seu rate-limit;
    * Ex: `10m`;
* IP_FRAME_TIME: Time-frame em que as requisições serão acumuladas;
    * Ex: `1s`;
* REDIS_ADDR: Endereço de acesso ao Redis;
* REDIS_USERNAME: Nome do usuário para acesso ao Redis, caso necessário;
* REDIS_PASSWORD: Senha para acesso ao Redis, caso necessário;
* PORT: Porta em o serviço de exemplo será executado;
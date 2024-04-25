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

### Aplicação

A aplicação implementa um `RateLimiter`, que recebe os parâmetros de configuração na função `NewRateLimiter`. O `RateLimiter` utiliza uma interface de `Memory` para realizar suas operações;

A aplicação também implementa um `Middleware`, que recebe um objeto `RateLimiter` e o utiliza para realizar controle de Rate Limit para uma função de `http.HandlerFunc`;

### Testes

A aplicação possui um exemplo de implementação:
* No exemplo da aplicação está uma implementação de `Memory` utilizando Redis;
* Requests para a rota `/`. As configurações do rate limiter de exemplo podem ser alteradas no arquivo `cmd/limiter/.env`;
* O arquivo `docs/api.http` possui exemplo de request para testar a API, utilizando a aplicação rodando localmente na porta 8080;
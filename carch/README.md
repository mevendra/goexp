## Clean Architecture APP

Aplicação para desafio de listagem de `orders`

### Como rodar/testar a aplicação?
* Executar comando `docker-compose up` para subir os serviços;
* Chamadas http para porta 8000 (.env);
  * Exemplo de istagem em `/api/get_order.http`
* Chamadas gql para porta 8080 (.env);
  * Playground gql em `locahost:8080`
* Chamadas grpc para porta 50051 (.env);
  * Evans conectando em `pb.OrderService`
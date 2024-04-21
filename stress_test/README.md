## Stress Test APP

Aplicação para desafio de sistema para StressTest;

Realiza teste de carga para uma determinada URL, utilizando os seguintes parâmetros via CLI:
* `url`: Url para a qual as requests serão enviadas;
* `requests`: Número total de requests realizadas para URL;
* `concurrency`: Número total de chamadas simultâneas realizadas para URL;

A aplicação irá informar um erro caso:
* Parâmetros CLI estejam faltando;
* Url mal formada;
* Falha na obtenção de respostas para url informada, de acordo com `requests``;

### Como rodar a aplicação?
* Realizar build/execução do arquivo `main.go`; ou
* Realizar build/execução da imagem docker, gerada atráves do arquivo `Dockerfile`;
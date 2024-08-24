
## Motivação

Este é um desafio da Pós-Graduação da Fullcycle: Go Expert, Desenvolvimento Avançado em Golang

## Desafio

**Objetivo:** 
Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

**Descrição:**
O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

1. Endereço IP: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
2. Token de Acesso: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
API_KEY: <TOKEN>
3. As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

**Requisitos:**

- O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
- O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
- O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
- As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
- Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
- O sistema deve responder adequadamente quando o limite é excedido:
	- Código HTTP: 429
	- Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame
- Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
- Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
- A lógica do limiter deve estar separada do middleware.

**Exemplos:**

1. Limitação por IP: Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.
2. Limitação por Token: Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.
3. Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.

**Dicas:**

- Teste seu rate limiter sob diferentes condições de carga para garantir que ele funcione conforme esperado em situações de alto tráfego.
Entrega:

- O código-fonte completo da implementação.
- Documentação explicando como o rate limiter funciona e como ele pode ser configurado.
- Testes automatizados demonstrando a eficácia e a robustez do rate limiter.
- Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
- O servidor web deve responder na porta 8080.

## Documentação


Para configurar o rate limiter, configure o arquivo .env conforme a sessão **Variáveis de Ambiente**.

Para configurar o Token, coloque o valor que você quer no Header: API_KEY, conforme a sessão **Headers**.

Para subir o servidor siga as instruções da sessão **Execução**.

Para testar unitáriamente e em lote siga as intruções da sessão 

---
### Quickstart Bash

Passo 1:

```Bash
git clone git@gitlab.com:thalesfaggiano/ge-rate-limiter.git
```

Passo 2:

```Bash
cd ge-rate-limiter
```

Passo 3:

```Bash
docker compose up --build
```

Passo 4:

Em outro terminal acesse o diretório ge-rate-limiter/teste

```Bash
cd ge-rate-limiter/teste
```

Passo 5:

```Bash
./teste.sh 5 8rm2332prqqfr4rh3d8ghei30ks3
```

Você pode mudar o primeiro parâmetro para executar a quantidade de iterações e o segundo parametro para testar os tokens.

---
### Variáveis de Ambiente

```Bash 
MAX_REQUESTS_PER_SECOND=1
BLOCK_DURATION_SECONDS=10
#Número máximo e tempo de bloqueio de requisições por segundo

IP_MAX_REQUESTS_PER_SECOND=1
IP_BLOCK_DURATION_SECONDS=15
#Número máximo e tempo de bloqueio de requisições por segundo por IP

TOKEN_MAX_REQUESTS_PER_SECOND=10
TOKEN_BLOCK_DURATION_SECONDS=15
#Número máximo e tempo de bloqueio de requisições por segundo por token
```

---
### Headers

Para limitação baseada em token, o cliente deve enviar o seguinte header:

```API_KEY: seu_token```

---

### Execução

O servidor funcionará na porta 8080 e para subi-lo execute o seguinte comando:

```Bash
docker compose up -d
```
---
#### Automated Tests via Shell Script

1. In the api folder, there are also two shell scripts to test the rate limiter:
    1. `api/test.sh `
    2. `api/test_rate_limiter.sh`
2. To run first give execution permissions to the scripts:
   ```sh
   chmod +x api/test.sh api/test_rate_limiter.sh
   ```
3. To run the test scripts, use the commands:
    ```sh
    ./api/test.sh
    ./api/test_rate_limiter.sh
    ```

4. View the container logs:
    ```sh
    docker logs -f rate_limiter_api
    ```
---
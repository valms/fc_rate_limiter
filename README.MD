# Rate Limiter - FULL CYCLE

Este projeto implementa um rate limiter em Go, utilizando o framework Fiber e o banco de dados Redis para controlar o
tráfego de requisições para um serviço web. O rate limiter pode ser configurado para limitar o número de requisições por
segundo com base em um endereço IP específico ou em um token de acesso.

O objetivo do rate limiter é proteger o serviço web contra abusos, garantindo que um número excessivo de requisições não
sobrecarregue o sistema. Ele é implementado como um middleware que pode ser facilmente integrado a qualquer aplicação Go
que utilize o Fiber.

## Como Funciona

### Principais Funcionalidades

- **Limitação por IP**: Restringe o número de requisições de um único endereço IP dentro de um intervalo de tempo
  definido.
- **Limitação por Token de Acesso**: Limita as requisições baseadas em um token de acesso único, permitindo diferentes
  limites de tempo de expiração para diferentes tokens.
- **Prioridade de Token**: As configurações de limite do token de acesso se sobrepõem às do IP.
- **Persistência com Redis**: Utiliza Redis para armazenar informações sobre o número de requisições feitas por cada IP
  ou token.

### Fluxo de Trabalho

1. **Identificação do Cliente**: Verifica se a requisição contém um token de acesso no cabeçalho `API_KEY`. Se não,
   utiliza o endereço IP.
2. **Verificação de Limite**: Verifica se o cliente excedeu o limite de requisições permitido.
3. **Bloqueio de Requisições**: Se o limite for excedido, a requisição é bloqueada e uma resposta com o código HTTP 429
   é retornada.
4. **Atualização de Contadores**: Se a requisição for permitida, o contador de requisições é incrementado.

## Configuração

A configuração do rate limiter é feita através de variáveis de ambiente, que podem ser definidas no arquivo
`docker-compose.yml`.

### Variáveis de Ambiente

- **REDIS_HOST**: Endereço do servidor Redis.
- **REDIS_PORT**: Porta do servidor Redis.
- **IPLIMIT**: Número máximo de requisições permitidas por segundo por IP.
- **APIKEYLIMIT**: Número máximo de requisições permitidas por segundo por token de acesso.
- **BLOCKDURATION**: Duração do bloqueio em segundos após o limite ser excedido.

## Instalação

Para instalar e executar o projeto localmente, siga os passos abaixo:

1. Clone o repositório:
   ```bash
   git clone https://github.com/valms/fc_rate_limiter
   cd seu_repositorio
   ```

2. Configure as variáveis de ambiente no arquivo `docker-compose.yml`.

3. Inicie os serviços usando Docker Compose:
   ```bash
   docker-compose up --build
   ```

## Uso

Após a instalação, o serviço estará disponível na porta 8080. Você pode testar o rate limiter enviando requisições HTTP
para o servidor e verificando se os limites são aplicados corretamente.

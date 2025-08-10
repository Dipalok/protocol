# Protocol (Go + Rust)

## Descrição

Este projeto tem como objetivo criar um protocolo de comunicação multiplataforma para envio e recebimento de mensagens (e-mail, SMS), utilizando as linguagens Go e Rust.  
A arquitetura contempla servidores e workers responsáveis pelo processamento, roteamento e integração das mensagens, além de suporte a banco de dados e possibilidade de expansão para outros tipos de comunicação.

## Estrutura do Projeto

```
src/
├── main.rs                # Lógica principal em Rust
├── database/
│   └── database.sql       # Script de banco de dados
├── go/
│   ├── server.go          # Servidor principal em Go
│   ├── db/
│   │   └── db.go          # Conexão e lógica de banco de dados em Go
│   ├── tls/
│   │   └── tls.Config.go  # Configuração TLS para segurança
│   └── worker/
│       └── worker.go      # Worker para processamento de mensagens
```

## Funcionalidades

- Estrutura modular em Go e Rust
- Suporte inicial para e-mail e SMS
- Integração com banco de dados relacional
- Workers para processamento assíncrono de mensagens
- Configuração de segurança via TLS
- Possibilidade de expansão para novos protocolos

## Como executar

### Pré-requisitos

- Go >= 1.20
- Rust >= 1.70
- PostgreSQL ou outro banco compatível
- Redis (opcional para fila de mensagens)

### Passos

1. Clone o repositório:
   ```bash
   git clone https://github.com/kovarike/protocol.git
   cd protocol
   ```

2. Configure o banco de dados:
   - Edite o arquivo `src/go/db/db.go` com suas credenciais
   - Execute o script `src/database/database.sql` para criar as tabelas

3. Compile e execute o servidor Go:
   ```bash
   cd src/go
   go run server.go
   ```

4. Compile e execute o módulo Rust:
   ```bash
   cd src
   cargo run
   ```

## Contribuição

1. Fork este repositório
2. Crie uma branch (`git checkout -b feature/nova-feature`)
3. Commit suas alterações (`git commit -m 'Minha nova feature'`)
4. Push para o branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## Roadmap

- [x] Estrutura inicial do projeto
- [x] Correção de imports e organização dos workers
- [ ] Modelagem das mensagens (e-mail, SMS, etc.)
- [ ] Implementação do servidor TCP/HTTP
- [ ] Integração entre Go e Rust
- [ ] Testes automatizados
- [ ] Documentação detalhada

## Licença

Este projeto está sob a licença MIT.

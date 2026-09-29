# Synkgo — Kanban de Evolução

> Painel de controle do ciclo autônomo de engenharia. Cada execução move exatamente
> **uma** tarefa do Backlog para In Progress, implementa em branch dedicada, valida com
> `go test -v ./...` (100% verde obrigatório) e só então move para Done.

## Estimativa
- **Escopo analisado em:** 2026-09-29
- **Ciclos previstos:** 8 madrugadas (1 análise/setup + 7 execução)
- **Critério de conclusão:** Backlog vazio, todas as tarefas em Done, `go test -v ./...` OK.

---

## Backlog

- [ ] **T2 - Testabilidade: extrair interfaces e isolar lógica pura**
  Extrair `isValidTableName`, construção de queries (`buildInsertQuery`,
  `buildPlaceholders`) e comparação de tabelas (`CompareTables`) como funções puras
  testáveis sem necessidade de banco real. Preparar terreno para os testes unitários.

- [ ] **T3 - Testes unitários base (funções puras)**
  Criar `internal/engine/engine_test.go` e `internal/config/config_test.go` cobrindo:
  validação de nome de tabela (casos válidos/inválidos/SQL injection), construção de
  query de insert, `CompareTables` com mocks de slices.

- [ ] **T4 - Tratamento de erros idiomático (remover panic/log.Fatal de libs)**
  Substituir `panic(err)` e `log.Fatal` dentro de `internal/engine` e
  `internal/config` por retorno de `error` explícito, propagado até `cmd/`. Comandos
  Cobra decidem se abortam ou logam.

- [ ] **T5 - Transação por tabela + insert em batch**
  Envolver truncate+insert de cada tabela em uma `sql.Tx`, com rollback em caso de
  erro parcial. Trocar insert linha-a-linha por batch (multi-row insert ou
  `Tx.Prepare` reaproveitado com `Exec` em lote configurável).

- [ ] **T6 - Cópia paralela configurável via goroutines**
  Adicionar campo `max_workers`/`concurrency` em `synkgo.json` (dto.Config).
  Implementar worker pool (goroutines + `sync.WaitGroup` + channel de tabelas) no
  `engine.Sync`, respeitando o limite configurado. Cada worker usa sua própria
  conexão/transação para evitar contenção.

- [ ] **T7 - Testes de integração com sqlmock para o pool de workers**
  Testar `engine.Sync` com `go-sqlmock` simulando N tabelas e verificando que o
  número de goroutines simultâneas nunca excede `max_workers`, e que falha em uma
  tabela não aborta as demais.

- [ ] **T8 - Modo dry-run e safety-check antes do TRUNCATE**
  Adicionar flag `--dry-run` no comando `sync` que apenas reporta contagem de linhas
  origem/destino e quais tabelas seriam truncadas, sem executar. Adicionar
  confirmação (`-y`/`--yes`) obrigatória para rodar de fato sem prompt interativo.

---

## In Progress

_(vazio)_

---

## Done

- [x] **T1 - Eliminar conexões hardcoded e código duplicado** *(2026-09-29)*
  `engine.Sync()` agora carrega o `synkgo.json` antes de conectar e usa o novo
  `config.OpenDB()`/`buildDSN()`; funções duplicadas `getTables`/`GetSourceTables`/
  `GetDestinationTables` removidas do engine (o pacote `config` é o único dono);
  `cmd/tables.go` aponta para `config.*`; `Sync()` retorna `error` propagado até
  `cmd/sync.go`. Testes: `TestBuildDSN` e `TestBuildDSN_EmptyFields` (2/2 PASS),
  build/vet limpos, zero strings `root:root@tcp` restantes.


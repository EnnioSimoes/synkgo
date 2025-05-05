-- Script SQL para criar tabelas e inserir 1000 registros em cada uma
-- Tabelas: clientes, fornecedores, produtos, pedidos e vendas

-- Remover tabelas se já existirem (em ordem inversa devido às chaves estrangeiras)
DROP TABLE IF EXISTS vendas;
DROP TABLE IF EXISTS pedidos;
DROP TABLE IF EXISTS produtos;
DROP TABLE IF EXISTS fornecedores;
DROP TABLE IF EXISTS clientes;

-- -----------------------------------------------------
-- Criação da tabela clientes
-- -----------------------------------------------------
CREATE TABLE clientes (
    id_cliente INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    telefone VARCHAR(20),
    endereco VARCHAR(200),
    cidade VARCHAR(50),
    estado CHAR(2),
    cep VARCHAR(10),
    data_cadastro DATETIME DEFAULT CURRENT_TIMESTAMP,
    limite_credito DECIMAL(10,2) DEFAULT 0.00,
    ativo BOOLEAN DEFAULT TRUE
);

-- -----------------------------------------------------
-- Criação da tabela fornecedores
-- -----------------------------------------------------
CREATE TABLE fornecedores (
    id_fornecedor INT AUTO_INCREMENT PRIMARY KEY,
    razao_social VARCHAR(100) NOT NULL,
    nome_fantasia VARCHAR(100),
    cnpj VARCHAR(20) NOT NULL,
    email VARCHAR(100),
    telefone VARCHAR(20),
    endereco VARCHAR(200),
    cidade VARCHAR(50),
    estado CHAR(2),
    cep VARCHAR(10),
    contato_nome VARCHAR(100),
    data_cadastro DATETIME DEFAULT CURRENT_TIMESTAMP,
    ativo BOOLEAN DEFAULT TRUE
);

-- -----------------------------------------------------
-- Criação da tabela produtos
-- -----------------------------------------------------
CREATE TABLE produtos (
    id_produto INT AUTO_INCREMENT PRIMARY KEY,
    id_fornecedor INT,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    preco_custo DECIMAL(10,2) NOT NULL,
    preco_venda DECIMAL(10,2) NOT NULL,
    estoque_atual INT DEFAULT 0,
    estoque_minimo INT DEFAULT 5,
    categoria VARCHAR(50),
    data_cadastro DATETIME DEFAULT CURRENT_TIMESTAMP,
    ativo BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (id_fornecedor) REFERENCES fornecedores(id_fornecedor)
);

-- -----------------------------------------------------
-- Criação da tabela pedidos
-- -----------------------------------------------------
CREATE TABLE pedidos (
    id_pedido INT AUTO_INCREMENT PRIMARY KEY,
    id_fornecedor INT,
    data_pedido DATETIME DEFAULT CURRENT_TIMESTAMP,
    valor_total DECIMAL(12,2) NOT NULL,
    status_pedido ENUM('Pendente', 'Aprovado', 'Em separação', 'Enviado', 'Entregue', 'Cancelado') DEFAULT 'Pendente',
    forma_pagamento ENUM('Boleto', 'Cartão', 'Transferência', 'A prazo'),
    prazo_entrega INT DEFAULT 15,
    observacoes TEXT,
    FOREIGN KEY (id_fornecedor) REFERENCES fornecedores(id_fornecedor)
);

-- -----------------------------------------------------
-- Criação da tabela vendas
-- -----------------------------------------------------
CREATE TABLE vendas (
    id_venda INT AUTO_INCREMENT PRIMARY KEY,
    id_cliente INT,
    id_produto INT,
    data_venda DATETIME DEFAULT CURRENT_TIMESTAMP,
    quantidade INT NOT NULL,
    preco_unitario DECIMAL(10,2) NOT NULL,
    valor_total DECIMAL(12,2) NOT NULL,
    forma_pagamento ENUM('Dinheiro', 'Cartão de Crédito', 'Cartão de Débito', 'Pix', 'Boleto'),
    status_venda ENUM('Concluída', 'Cancelada', 'Em processamento') DEFAULT 'Concluída',
    FOREIGN KEY (id_cliente) REFERENCES clientes(id_cliente),
    FOREIGN KEY (id_produto) REFERENCES produtos(id_produto)
);

-- -----------------------------------------------------
-- Inserção de dados nas tabelas
-- -----------------------------------------------------

-- Inserir 1000 clientes
DELIMITER $$

-- Apaga a procedure se ela já existir, para facilitar testes repetidos
DROP PROCEDURE IF EXISTS inserir_clientes
$$

CREATE PROCEDURE inserir_clientes()
BEGIN
    -- Variáveis de controle e configuração
    DECLARE i INT DEFAULT 1;
    DECLARE max_clientes INT DEFAULT 1000; -- Tornar o limite configurável

    -- Variáveis para os dados gerados em cada iteração
    DECLARE nome_cliente VARCHAR(100);
    DECLARE email_cliente VARCHAR(100);
    DECLARE telefone_gerado VARCHAR(20); -- Aumentado para acomodar '(XX) 9XXXXXXXX'
    DECLARE endereco_gerado VARCHAR(150);
    DECLARE cidade_nome VARCHAR(50);
    DECLARE estado CHAR(2);
    DECLARE cep_gerado VARCHAR(9); -- Formato XXXXX-XXX
    DECLARE limite_gerado DECIMAL(10, 2);
    DECLARE ativo_gerado BOOLEAN;
    DECLARE digito9 CHAR(1);

    -- Lista CORRIGIDA e ÚNICA de 27 siglas de estados brasileiros
    DECLARE estados CHAR(54) DEFAULT 'SP';

    -- Variáveis para construção do INSERT em lote (batch)
    DECLARE sql_values LONGTEXT DEFAULT ''; -- Armazena a parte VALUES (...) , (...)
    DECLARE batch_separator CHAR(1) DEFAULT ''; -- Usado para adicionar a vírgula entre as tuplas

    -- Handler para Erros: Faz ROLLBACK em caso de qualquer exceção SQL
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK; -- Desfaz a transação em caso de erro
        RESIGNAL; -- Opcional: Re-emite o erro para quem chamou a procedure
    END;

    -- Inicia a Transação: Garante atomicidade e melhora performance
    START TRANSACTION;

    -- Loop para gerar os dados (mas não para inserir linha a linha)
    WHILE i <= max_clientes DO
        -- 1. Gerar dados para um cliente
        SET nome_cliente = CONCAT('Cliente ', i);
        SET email_cliente = CONCAT('cliente', i, '@email.com');

        -- Seleciona o estado usando a lista corrigida (índice 0-26)
        SET estado = SUBSTRING(estados, 1 + ((i - 1) % 27) * 2, 2);

        -- Gera cidade baseada no estado (lógica mantida, pode ser melhorada para mais realismo)
        CASE estado
            WHEN 'SP' THEN SET cidade_nome = ELT(1 + (i % 5), 'São Paulo', 'Campinas', 'Ribeirão Preto', 'Santos', 'São José dos Campos');
            WHEN 'RJ' THEN SET cidade_nome = ELT(1 + (i % 3), 'Rio de Janeiro', 'Niterói', 'Petrópolis');
            WHEN 'MG' THEN SET cidade_nome = ELT(1 + (i % 4), 'Belo Horizonte', 'Uberlândia', 'Juiz de Fora', 'Contagem');
            WHEN 'RS' THEN SET cidade_nome = ELT(1 + (i % 3), 'Porto Alegre', 'Caxias do Sul', 'Pelotas');
            ELSE SET cidade_nome = CONCAT('Cidade ', ((i - 1) % 20) + 1, ' ', estado); -- Adiciona UF na cidade genérica
        END CASE;

        -- Gera telefone (70% chance de ter 9 dígitos após DDD)
        IF RAND() < 0.7 THEN SET digito9 = '9'; ELSE SET digito9 = ''; END IF;
        -- Ajuste leve no DDD para evitar 10 ou 99. Garante 8 ou 9 dígitos no número.
        SET telefone_gerado = CONCAT('(', FLOOR(11 + RAND() * 88), ') ', digito9, LPAD(FLOOR(10000000 + RAND() * 89999999), IF(digito9 = '9', 9, 8), '0'));

        -- Gera endereço
        SET endereco_gerado = CONCAT('Rua Exemplo ', CHAR(65 + ((i - 1) % 26)), ', ', FLOOR(1 + RAND() * 1999));

        -- Gera CEP formatado XXXXX-XXX (usa LPAD para garantir zeros à esquerda se necessário)
        SET cep_gerado = CONCAT(LPAD(FLOOR(10000 + RAND() * 89999), 5, '0'), '-', LPAD(FLOOR(100 + RAND() * 899), 3, '0'));

        -- Gera limite de crédito
        SET limite_gerado = ROUND(1000 + RAND() * 9000, 2);

        -- Gera status ativo (TRUE para ~90% dos casos)
        SET ativo_gerado = (RAND() > 0.1);

        -- 2. Construir a string de valores para o INSERT em lote
        -- Usa QUOTE() para tratar strings corretamente (evita SQL injection e erros com aspas)
        SET sql_values = CONCAT(sql_values,
            batch_separator, -- Adiciona ',' a partir da segunda linha
            '(',
            QUOTE(nome_cliente), ',',
            QUOTE(email_cliente), ',',
            QUOTE(telefone_gerado), ',',
            QUOTE(endereco_gerado), ',',
            QUOTE(cidade_nome), ',',
            QUOTE(estado), ',',
            QUOTE(cep_gerado), ',',
            limite_gerado, ',',  -- Numérico, sem QUOTE
            ativo_gerado,         -- Boolean, sem QUOTE
            ')'
        );

        SET batch_separator = ','; -- Define o separador para as próximas iterações
        SET i = i + 1;

    END WHILE;

    -- 3. Executar o INSERT em lote APENAS SE houver valores gerados
    IF LENGTH(sql_values) > 0 THEN
        -- Monta a query SQL final
        SET @full_sql = CONCAT(
            'INSERT INTO clientes (nome, email, telefone, endereco, cidade, estado, cep, limite_credito, ativo) VALUES ',
            sql_values
        );

        -- Prepara e executa a query dinâmica
        PREPARE stmt FROM @full_sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt; -- Libera a memória do statement preparado
    END IF;

    -- Finaliza a Transação com sucesso
    COMMIT;

END $$

-- Restaura o delimitador padrão
DELIMITER ;

CALL inserir_clientes();
DROP PROCEDURE inserir_clientes;

-- Inserir 1000 fornecedores
DELIMITER $$

DROP PROCEDURE IF EXISTS inserir_fornecedores
$$

CREATE PROCEDURE inserir_fornecedores()
BEGIN
    DECLARE i INT DEFAULT 1;
    DECLARE razao_social_fornecedor VARCHAR(100);
    DECLARE nome_fantasia_fornecedor VARCHAR(100);
    DECLARE cnpj_fornecedor VARCHAR(20);
    DECLARE estados CHAR(54) DEFAULT 'SP';
    DECLARE estado CHAR(2);
    
    WHILE i <= 1000 DO
        SET razao_social_fornecedor = CONCAT('Fornecedor ', i, ' LTDA');
        SET nome_fantasia_fornecedor = CONCAT('Fornec ', i);
        -- CNPJ formato XX.XXX.XXX/0001-XX (não validado)
        SET cnpj_fornecedor = CONCAT(
            LPAD(FLOOR(RAND() * 99), 2, '0'), '.',
            LPAD(FLOOR(RAND() * 999), 3, '0'), '.',
            LPAD(FLOOR(RAND() * 999), 3, '0'), '/0001-',
            LPAD(FLOOR(RAND() * 99), 2, '0')
        );
        SET estado = SUBSTRING(estados, 1 + (i % 27) * 2, 2);
        
        INSERT INTO fornecedores (
            razao_social, 
            nome_fantasia, 
            cnpj, 
            email, 
            telefone, 
            endereco, 
            cidade, 
            estado, 
            cep, 
            contato_nome, 
            ativo
        )
        VALUES (
            razao_social_fornecedor,
            nome_fantasia_fornecedor,
            cnpj_fornecedor,
            CONCAT('contato@fornecedor', i, '.com.br'),
            CONCAT('(', FLOOR(10 + RAND() * 89), ') ', FLOOR(10000000 + RAND() * 89999999)),
            CONCAT('Av. Principal, ', FLOOR(100 + RAND() * 9900)),
            CASE estado
                WHEN 'SP' THEN ELT(1 + (i % 5), 'São Paulo', 'Campinas', 'Ribeirão Preto', 'Santos', 'São José dos Campos')
                WHEN 'RJ' THEN ELT(1 + (i % 3), 'Rio de Janeiro', 'Niterói', 'Petrópolis')
                WHEN 'MG' THEN ELT(1 + (i % 4), 'Belo Horizonte', 'Uberlândia', 'Juiz de Fora', 'Contagem')
                WHEN 'RS' THEN ELT(1 + (i % 3), 'Porto Alegre', 'Caxias do Sul', 'Pelotas')
                ELSE CONCAT('Cidade ', (i % 20) + 1)
            END,
            estado,
            CONCAT(FLOOR(10000 + RAND() * 89999), '-', FLOOR(100 + RAND() * 899)),
            CONCAT('Contato ', CHAR(65 + (i % 26)), '. ', CHAR(65 + ((i+7) % 26))),
            IF(RAND() > 0.05, TRUE, FALSE)  -- 95% ativos
        );
        
        SET i = i + 1;
    END WHILE;
END $$
DELIMITER ;

CALL inserir_fornecedores();
DROP PROCEDURE inserir_fornecedores;

-- Inserir 1000 produtos
DELIMITER $$

DROP PROCEDURE IF EXISTS inserir_produtos
$$

CREATE PROCEDURE inserir_produtos()
BEGIN
    DECLARE i INT DEFAULT 1;
    DECLARE produto_nome VARCHAR(100);
    DECLARE categorias VARCHAR(300) DEFAULT 'Eletrônicos,Informática,Móveis,Decoração,Cozinha,Jardinagem,Ferramentas,Brinquedos,Livros,Material Escolar,Esportes,Vestuário,Calçados,Acessórios,Beleza,Saúde';
    DECLARE categoria VARCHAR(50);
    DECLARE preco_custo_valor DECIMAL(10,2);
    DECLARE qtd_fornecedores INT;
    
    SELECT COUNT(*) INTO qtd_fornecedores FROM fornecedores;
    
    WHILE i <= 1000 DO
        SET produto_nome = CONCAT('Produto ', i);
        -- Selecionar uma categoria da lista
        SET categoria = ELT(1 + (i % 16), 'Eletrônicos', 'Informática', 'Móveis', 'Decoração', 'Cozinha', 
                                         'Jardinagem', 'Ferramentas', 'Brinquedos', 'Livros', 'Material Escolar', 
                                         'Esportes', 'Vestuário', 'Calçados', 'Acessórios', 'Beleza', 'Saúde');
        -- Gerar preço de custo baseado na categoria
        CASE categoria
            WHEN 'Eletrônicos' THEN SET preco_custo_valor = 500 + RAND() * 2000;
            WHEN 'Informática' THEN SET preco_custo_valor = 300 + RAND() * 1500;
            WHEN 'Móveis' THEN SET preco_custo_valor = 200 + RAND() * 800;
            WHEN 'Decoração' THEN SET preco_custo_valor = 50 + RAND() * 200;
            WHEN 'Cozinha' THEN SET preco_custo_valor = 30 + RAND() * 300;
            WHEN 'Esportes' THEN SET preco_custo_valor = 80 + RAND() * 400;
            ELSE SET preco_custo_valor = 15 + RAND() * 150;
        END CASE;
        
        INSERT INTO produtos (
            id_fornecedor,
            nome,
            descricao,
            preco_custo,
            preco_venda,
            estoque_atual,
            estoque_minimo,
            categoria,
            ativo
        )
        VALUES (
            1 + FLOOR(RAND() * qtd_fornecedores),  -- Fornecedor aleatório
            CONCAT(categoria, ' - ', produto_nome),
            CONCAT('Descrição detalhada do produto ', i, ' da categoria ', categoria),
            ROUND(preco_custo_valor, 2),  -- Preço de custo
            ROUND(preco_custo_valor * (1.3 + RAND() * 0.7), 2),  -- Margem de 30% a 100%
            FLOOR(10 + RAND() * 200),  -- Estoque atual
            FLOOR(5 + RAND() * 25),  -- Estoque mínimo
            categoria,
            IF(RAND() > 0.05, TRUE, FALSE)  -- 95% ativos
        );
        
        SET i = i + 1;
    END WHILE;
END $$
DELIMITER ;

CALL inserir_produtos();
DROP PROCEDURE inserir_produtos;

-- Inserir 1000 pedidos
DELIMITER $$

DROP PROCEDURE IF EXISTS inserir_produtos
$$

CREATE PROCEDURE inserir_pedidos()
BEGIN
    DECLARE i INT DEFAULT 1;
    DECLARE qtd_fornecedores INT;
    DECLARE valor_total_pedido DECIMAL(12,2);
    DECLARE data_pedido_valor DATETIME;
    DECLARE status_pedido_valor VARCHAR(20);
    DECLARE forma_pagamento_valor VARCHAR(20);
    
    SELECT COUNT(*) INTO qtd_fornecedores FROM fornecedores;
    
    WHILE i <= 1000 DO
        -- Gerar um valor total aleatório para o pedido
        SET valor_total_pedido = ROUND(500 + RAND() * 15000, 2);
        
        -- Gerar uma data de pedido nos últimos 2 anos
        SET data_pedido_valor = DATE_SUB(NOW(), INTERVAL FLOOR(RAND() * 730) DAY);
        
        -- Determinar status do pedido
        SET status_pedido_valor = ELT(1 + FLOOR(RAND() * 6), 'Pendente', 'Aprovado', 'Em separação', 'Enviado', 'Entregue', 'Cancelado');
        
        -- Determinar forma de pagamento
        SET forma_pagamento_valor = ELT(1 + FLOOR(RAND() * 4), 'Boleto', 'Cartão', 'Transferência', 'A prazo');
        
        INSERT INTO pedidos (
            id_fornecedor,
            data_pedido,
            valor_total,
            status_pedido,
            forma_pagamento,
            prazo_entrega,
            observacoes
        )
        VALUES (
            1 + FLOOR(RAND() * qtd_fornecedores),  -- Fornecedor aleatório
            data_pedido_valor,
            valor_total_pedido,
            status_pedido_valor,
            forma_pagamento_valor,
            10 + FLOOR(RAND() * 20),  -- Prazo de entrega entre 10 e 30 dias
            CASE WHEN RAND() > 0.7 THEN CONCAT('Observações do pedido ', i) ELSE NULL END  -- 30% tem observações
        );
        
        SET i = i + 1;
    END WHILE;
END $$
DELIMITER ;

CALL inserir_pedidos();
DROP PROCEDURE inserir_pedidos;

-- Inserir 1000 vendas
DELIMITER $$

DROP PROCEDURE IF EXISTS inserir_vendas
$$

CREATE PROCEDURE inserir_vendas()
BEGIN
    DECLARE i INT DEFAULT 1;
    DECLARE qtd_clientes INT;
    DECLARE qtd_produtos INT;
    DECLARE id_cliente_valor INT;
    DECLARE id_produto_valor INT;
    DECLARE preco_unitario_valor DECIMAL(10,2);
    DECLARE quantidade_valor INT;
    DECLARE data_venda_valor DATETIME;
    DECLARE forma_pagamento_valor VARCHAR(20);
    DECLARE status_venda_valor VARCHAR(20);
    
    SELECT COUNT(*) INTO qtd_clientes FROM clientes;
    SELECT COUNT(*) INTO qtd_produtos FROM produtos;
    
    WHILE i <= 1000 DO
        -- Selecionar cliente e produto aleatórios
        SET id_cliente_valor = 1 + FLOOR(RAND() * qtd_clientes);
        SET id_produto_valor = 1 + FLOOR(RAND() * qtd_produtos);
        
        -- Obter preço do produto
        SELECT preco_venda INTO preco_unitario_valor FROM produtos WHERE id_produto = id_produto_valor;
        
        -- Gerar quantidade aleatória
        SET quantidade_valor = 1 + FLOOR(RAND() * 10);
        
        -- Gerar data de venda nos últimos 12 meses
        SET data_venda_valor = DATE_SUB(NOW(), INTERVAL FLOOR(RAND() * 365) DAY);
        
        -- Determinar forma de pagamento
        SET forma_pagamento_valor = ELT(1 + FLOOR(RAND() * 5), 'Dinheiro', 'Cartão de Crédito', 'Cartão de Débito', 'Pix', 'Boleto');
        
        -- Determinar status da venda
        SET status_venda_valor = ELT(1 + FLOOR(RAND() * 3), 'Concluída', 'Cancelada', 'Em processamento');
        
        INSERT INTO vendas (
            id_cliente,
            id_produto,
            data_venda,
            quantidade,
            preco_unitario,
            valor_total,
            forma_pagamento,
            status_venda
        )
        VALUES (
            id_cliente_valor,
            id_produto_valor,
            data_venda_valor,
            quantidade_valor,
            preco_unitario_valor,
            ROUND(quantidade_valor * preco_unitario_valor, 2),
            forma_pagamento_valor,
            status_venda_valor
        );
        
        SET i = i + 1;
    END WHILE;
END $$
DELIMITER ;

CALL inserir_vendas();
DROP PROCEDURE inserir_vendas;

-- Confirmar número de registros em cada tabela
SELECT 'clientes' AS tabela, COUNT(*) AS registros FROM clientes
UNION ALL
SELECT 'fornecedores', COUNT(*) FROM fornecedores
UNION ALL
SELECT 'produtos', COUNT(*) FROM produtos
UNION ALL
SELECT 'pedidos', COUNT(*) FROM pedidos
UNION ALL
SELECT 'vendas', COUNT(*) FROM vendas;
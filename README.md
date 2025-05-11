# Desafio 1 - Pós Go Full Cycle

## Executando o Sistema

Para executar o sistema é necessário executar inicialmente o servidor e, então, realizar chamadas com o cliente.

### **Executando o servidor**

1. Abra um terminal 
2. A partir da pasta raiz, mude para a pasta do servidor:

```bash
cd src/server
```

3. rode o código do servidor

```bash
go run .
```

O servidor deverá rodar indefinidamente até que seja interrompido.

### **Executando o cliente**

Executa-se o cliente de forma semelhante ao servidor.

1. Abra um terminal
2. A partir da pasta raiz, mude para a pasta do cliente
```bash
cd src/client
```
3. rode o código do cliente

```bash
go run .
```

cada execução do cliente fará apenas uma chamada.
A saída do sistema está, como posto nos requisitos, no arquivo cotação.txt
# KpowBIT Node Operator Guide

This guide covers the recommended ways to launch and operate a KpowBIT node using either the native client or Docker.

## Requirements

You will need:

* Linux, macOS, or Windows
* A stable network connection
* Storage for blockchain data
* Either:

  * [KpowBIT Geth](https://github.com/KpowBIT-Project/kpow-geth/releases)
  * Docker

---

## Native Client

### Start a Node

```bash
geth --kpow \
  --bootnodes enode://76b5b9e1cdcf1c188e82268cd9a85a52f9d622e6bb6fad412fd6edeca73574eb58fe4e531187841d5a4de7b1a71009d6666aecff3dcecc2140862e334b0d5b41@175.110.114.159:30303,enode://365d0bc3bc6df393a69a332d96bbf9e820db181fe56cf665e411ce9e1d4545f66775c44d9af4438b8ca951254555140795f211874b88998e866a031cfcf0fda6@159.223.187.113:30303 \
  --syncmode full \
  --http \
  --http.addr 127.0.0.1 \
  --http.port 8545 \
  --http.api eth,net,web3
```

This starts a full-sync KpowBIT node and exposes JSON-RPC locally at:

```text
http://127.0.0.1:8545
```

To allow RPC access from other hosts, change:

```bash
--http.addr 127.0.0.1
```

to:

```bash
--http.addr 0.0.0.0
```

Only do this when the RPC endpoint is protected by firewall rules or other access controls.

---

## Mining Mode

Mining can be enabled by adding the miner configuration:

```bash
geth --kpow \
  --bootnodes enode://76b5b9e1cdcf1c188e82268cd9a85a52f9d622e6bb6fad412fd6edeca73574eb58fe4e531187841d5a4de7b1a71009d6666aecff3dcecc2140862e334b0d5b41@175.110.114.159:30303,enode://365d0bc3bc6df393a69a332d96bbf9e820db181fe56cf665e411ce9e1d4545f66775c44d9af4438b8ca951254555140795f211874b88998e866a031cfcf0fda6@159.223.187.113:30303 \
  --syncmode full \
  --http \
  --http.addr 127.0.0.1 \
  --http.port 8545 \
  --http.api eth,net,web3,miner \
  --mine \
  --miner.threads 1 \
  --miner.etherbase 0xYourAddress
```

Replace `0xYourAddress` with the address that should receive mining rewards.

Key options:

```text
--mine              Enable mining
--miner.threads     Number of mining threads
--miner.etherbase   Reward address
```

---

## Docker Deployment

Create a persistent data directory:

```bash
mkdir -p data
```

Start the node:

```bash
docker run -d \
  --name kpowbit-node \
  --restart unless-stopped \
  -v "$(pwd)/data:/node" \
  -p 30303:30303/tcp \
  -p 30303:30303/udp \
  -p 127.0.0.1:8545:8545 \
  ghcr.io/kpowbit-project/kpow-geth:stable --kpow \
  --datadir /node \
  --bootnodes enode://76b5b9e1cdcf1c188e82268cd9a85a52f9d622e6bb6fad412fd6edeca73574eb58fe4e531187841d5a4de7b1a71009d6666aecff3dcecc2140862e334b0d5b41@175.110.114.159:30303,enode://365d0bc3bc6df393a69a332d96bbf9e820db181fe56cf665e411ce9e1d4545f66775c44d9af4438b8ca951254555140795f211874b88998e866a031cfcf0fda6@159.223.187.113:30303 \
  --syncmode full \
  --http \
  --http.addr 0.0.0.0 \
  --http.port 8545 \
  --http.api eth,net,web3
```

Blockchain data is stored in:

```text
./data
```

The data remains available after the container is restarted or recreated.

### Container Operations

View logs:

```bash
docker logs -f kpowbit-node
```

Restart:

```bash
docker restart kpowbit-node
```

Stop:

```bash
docker stop kpowbit-node
```

Check status:

```bash
docker ps --filter name=kpowbit-node
```

---

## Network Endpoints

Typical ports used by the node:

| Endpoint  |            Port | Purpose                          |
| --------- | --------------: | -------------------------------- |
| P2P       | `30303` TCP/UDP | Peer communication and discovery |
| HTTP RPC  |          `8545` | JSON-RPC requests                |
| WebSocket |          `8546` | Persistent RPC connections       |

Custom ports can be configured with:

```bash
--port 30304 \
--http.port 9545 \
--ws.port 9546
```

When using Docker, update the published ports accordingly.

---

## WebSocket Access

Enable WebSocket RPC when required by applications:

```bash
geth --kpow \
  --bootnodes enode://76b5b9e1cdcf1c188e82268cd9a85a52f9d622e6bb6fad412fd6edeca73574eb58fe4e531187841d5a4de7b1a71009d6666aecff3dcecc2140862e334b0d5b41@175.110.114.159:30303,enode://365d0bc3bc6df393a69a332d96bbf9e820db181fe56cf665e411ce9e1d4545f66775c44d9af4438b8ca951254555140795f211874b88998e866a031cfcf0fda6@159.223.187.113:30303 \
  --syncmode full \
  --ws \
  --ws.addr 127.0.0.1 \
  --ws.port 8546 \
  --ws.api eth,net,web3
```

Endpoint:

```text
ws://127.0.0.1:8546
```

---

## Node Health Checks

### Peer Connectivity

```bash
curl -s \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' \
  http://127.0.0.1:8545
```

A value greater than `0x0` indicates connected peers.

### Latest Block

```bash
curl -s \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://127.0.0.1:8545
```

The returned block number should continue increasing as the chain progresses.

### Sync State

```bash
curl -s \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' \
  http://127.0.0.1:8545
```

After synchronization completes:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": false
}
```

---

## Recommended RPC Policy

For most installations, keep RPC bound to localhost:

```bash
--http.addr 127.0.0.1
```

A minimal API set is:

```bash
--http.api eth,net,web3
```

Avoid exposing unrestricted RPC directly to the public internet.

For remote access, place the node behind infrastructure such as:

```text
Application
    │
    ▼
Reverse Proxy / Firewall
    │
    ▼
KpowBIT RPC
```

This makes it easier to add authentication, TLS, IP filtering, and rate limits.

---

## Common Issues

### No Peer Connections

Check:

* Internet connectivity
* TCP/UDP access to the P2P port
* Firewall configuration
* Node logs

Docker logs:

```bash
docker logs --tail 200 kpowbit-node
```

### RPC Is Unreachable

Confirm that RPC is enabled:

```bash
--http
```

Also make sure the configured address and port match the endpoint you are using.

### Docker Stops After Launch

Inspect the startup output:

```bash
docker logs kpowbit-node
```

Typical causes include invalid flags, port conflicts, or data-directory permissions.

### Port Conflict

Linux/macOS:

```bash
lsof -i :8545
lsof -i :30303
```

Use another port or stop the conflicting service.

---

## Recommended Production Setup

A typical long-running node should use:

* Persistent SSD storage
* Automatic restart
* Public P2P connectivity
* Private RPC access
* Minimal RPC APIs
* Firewall rules
* Regular log and disk monitoring

Example Docker configuration:

```bash
docker run -d \
  --name kpowbit-node \
  --restart unless-stopped \
  -v "$(pwd)/data:/node" \
  -p 30303:30303/tcp \
  -p 30303:30303/udp \
  -p 127.0.0.1:8545:8545 \
  ghcr.io/kpowbit-project/kpow-geth:stable --kpow \
  --bootnodes enode://76b5b9e1cdcf1c188e82268cd9a85a52f9d622e6bb6fad412fd6edeca73574eb58fe4e531187841d5a4de7b1a71009d6666aecff3dcecc2140862e334b0d5b41@175.110.114.159:30303,enode://365d0bc3bc6df393a69a332d96bbf9e820db181fe56cf665e411ce9e1d4545f66775c44d9af4438b8ca951254555140795f211874b88998e866a031cfcf0fda6@159.223.187.113:30303 \
  --datadir /node \
  --syncmode full \
  --http \
  --http.addr 0.0.0.0 \
  --http.port 8545 \
  --http.api eth,net,web3
```

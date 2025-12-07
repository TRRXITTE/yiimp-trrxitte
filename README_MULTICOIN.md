# Multicoin Merged Mining Pool

A high-performance merged mining pool supporting Litecoin as the primary chain and multiple Scrypt-based auxiliary chains simultaneously. Built in Go with a modern React dashboard.

## Features

### Core Features
- **Multicoin Merged Mining**: Mine Litecoin + up to 35 auxiliary chains simultaneously
- **Merkle Tree AuxPow**: Proper auxiliary chain Merkle tree implementation for 3+ chains
- **Variable Difficulty (varDiff)**: Dynamic difficulty adjustment per miner
- **PPLNS Payout Scheme**: Fair pay-per-last-N-shares distribution
- **Solo Mining Support**: Optional solo mining mode
- **Real-time Blockchain Sync Monitoring**: Track daemon synchronization status
- **Percentage-based Pool Fees**: Configurable fee distribution per chain
- **RESTful API**: Comprehensive API for pool statistics and management
- **Modern Web Dashboard**: React-based dashboard with real-time updates

### Supported Coins (35+ chains)
Primary: Litecoin

Auxiliary chains: Dogecoin, Bellscoin, Pepecoin, Luckycoin, Junkcoin, Dingocoin, Dogmcoin, Craftcoin, Newyorkcoin, Earthcoin, Worldcoin, Shibacoin, Beerscoin, Dogecoinev, Bonkcoin, Flincoin, Marscoin, BBQcoin, Goldcoin, Catcoin, Cyberyen, Infinitecoin, IBithub, Newenglandcoin, Bitbar, Ferrite, Flopcoin, Stohncoin, Sorachancoin, Mooncoin, Fairbrix, Lebowskiscoin, Bit, Trumpow, Mydogecoin

## Architecture

### Merged Mining Process

1. **Template Fetching**: Pool fetches block templates from all configured chains
2. **Merkle Tree Building**: Constructs Merkle tree of all auxiliary block hashes
3. **Coinbase Embedding**: Embeds Merkle root in primary chain's coinbase transaction
4. **Work Distribution**: Sends work to miners with embedded Merkle data
5. **Share Validation**: Validates submitted shares against all chain targets
6. **Block Submission**: Submits blocks to all qualifying chains with proper AuxPow proofs

### Key Components

- **bitcoin/auxmerkle.go**: Merkle tree builder for auxiliary chains
- **bitcoin/auxpow.go**: AuxPow structure and Merkle branch handling
- **bitcoin/coins.go**: Individual coin implementations (35+ chains)
- **pool/work.go**: Work generation and block submission
- **pool/share.go**: Share validation and difficulty calculation
- **pool/server.go**: Stratum server and RPC template fetching
- **api/enhanced_api.go**: REST API endpoints
- **dashboard/**: React TypeScript web dashboard

## Installation

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Node.js 18+ (for dashboard)
- Running daemons for each chain you want to mine

### Building the Pool

```bash
# Build the pool server
go mod tidy
go build -o poolserver

# The binary will be created as 'poolserver'
```

### Setting Up the Database

```bash
# Create PostgreSQL database
createdb miningpool

# Run migrations (if migrations exist)
# psql miningpool < migrations/001_initial.sql
```

### Configuration

Copy one of the example configurations:

```bash
# For testing with 5 chains
cp config_test_multicoin.json config.json

# For production with 8 chains
cp config_all_coins.json config.json
```

Edit `config.json` and update:
- RPC URLs and credentials for your daemons
- Reward addresses for each chain
- Pool fee addresses and percentages
- PostgreSQL connection details
- Port configurations and difficulty settings

### Running the Pool

```bash
# Start the pool server
./poolserver -config config.json

# The pool will listen on the configured ports (default: 3333-3342)
# The API will be available on port 8001
```

### Setting Up the Dashboard

```bash
cd dashboard
npm install
npm run build

# For development
npm run dev

# For production, serve the dist/ folder with nginx or another web server
```

## Configuration Guide

### Basic Structure

```json
{
  "pool_name": "Your Pool Name",
  "merged_blockchain_order": ["litecoin", "dogecoin", "bellscoin"],
  "blockchains": {
    "litecoin": [{ "rpc_url": "...", "reward_to": "..." }],
    "dogecoin": [{ "rpc_url": "...", "reward_to": "..." }]
  },
  "ports": {
    "3333": {
      "difficulty": 1024,
      "varDiff": { "enabled": true, "minDiff": 256, "maxDiff": 16384 },
      "payoutScheme": "PPLNS"
    }
  },
  "payouts": {
    "chains": {
      "litecoin": {
        "miner_min_payment": 0.1,
        "pool_rewards": [{"address": "...", "percentage": 1.0}]
      }
    }
  }
}
```

### Port Configuration

Configure multiple ports for different miner types:

- **Low difficulty ports** (3334): For small miners, GPU miners
- **Standard ports** (3333): For medium hashrate miners
- **High difficulty ports** (3336): For ASIC farms
- **Solo ports** (3340-3342): For solo mining

### Variable Difficulty Settings

```json
"varDiff": {
  "enabled": true,
  "minDiff": 256,        // Minimum difficulty
  "maxDiff": 16384,      // Maximum difficulty
  "targetTime": 15,      // Target seconds per share
  "retargetTime": 90,    // Retarget interval in seconds
  "variancePercent": 30  // Allowed variance percentage
}
```

### Pool Fees

Configure percentage-based fees per chain:

```json
"pool_rewards": [
  {"address": "LTC_FEE_ADDRESS", "percentage": 1.0},
  {"address": "LTC_DEV_ADDRESS", "percentage": 0.5}
]
```

## API Documentation

The pool provides a comprehensive REST API on port 8001.

### Pool Endpoints

- `GET /api/pool/stats` - Pool statistics (hashrate, miners, blocks, chains)
- `GET /api/pool/blocks?limit=50` - Recent blocks from all chains
- `GET /api/pool/payments?limit=50` - Recent payments
- `GET /api/pool/miners` - Active miners list
- `GET /api/pool/sync` - Blockchain synchronization status
- `GET /api/pool/hashrate` - Pool hashrate history

### Miner Endpoints

- `GET /api/miner/{address}/stats` - Miner statistics
- `GET /api/miner/{address}/payments` - Miner payment history
- `GET /api/miner/{address}/blocks` - Blocks found by miner
- `GET /api/miner/{address}/workers` - Worker statistics
- `GET /api/miner/{address}/balance` - Balance per chain

### Chain Endpoints

- `GET /api/chain/{chain}/stats` - Chain-specific statistics
- `GET /api/chain/{chain}/blocks` - Blocks for specific chain
- `GET /api/chain/{chain}/difficulty` - Chain difficulty history

### Configuration Endpoints

- `GET /api/config/chains` - List of configured chains
- `GET /api/config/fees` - Pool fee percentages
- `GET /api/config/ports` - Pool ports configuration

## Mining

### Connect to the Pool

```bash
# Standard Stratum connection
stratum+tcp://your-pool-address:3333

# Username format
YOUR_WALLET_ADDRESS.WORKER_NAME

# Password (optional)
x
```

### Example with cpuminer

```bash
cpuminer -a scrypt \
  -o stratum+tcp://pool.example.com:3333 \
  -u LYourLitecoinAddress.worker1 \
  -p x
```

### Solo Mining

Connect to a solo port (3340-3342) to mine solo. If you find a block, you receive the full block reward (minus pool fee) for all chains.

## Dashboard

The web dashboard provides:

- Real-time pool statistics
- Multi-chain monitoring
- Blockchain sync status
- Block explorer
- Payment history
- Miner lookup and statistics
- Worker performance monitoring

Access at: `http://your-pool:3000`

## Daemon Configuration

Each daemon must be configured for merged mining:

### Litecoin (Primary Chain)

```conf
server=1
rpcuser=litecoinrpc
rpcpassword=changeme
rpcport=9332
rpcallowip=127.0.0.1

# Block notifications
blocknotify=curl -X POST http://127.0.0.1:8001/notify/litecoin
```

### Auxiliary Chains (Dogecoin, etc.)

```conf
server=1
rpcuser=dogecoinrpc
rpcpassword=changeme
rpcport=22555
rpcallowip=127.0.0.1

# Enable auxiliary proof-of-work
auxpow=1

# Block notifications
blocknotify=curl -X POST http://127.0.0.1:8001/notify/dogecoin
```

## Performance Tuning

### Database Optimization

```sql
-- Increase connection pool
ALTER SYSTEM SET max_connections = 200;

-- Optimize for writes
ALTER SYSTEM SET synchronous_commit = off;
```

### Pool Settings

```json
{
  "max_connections": 8192,
  "share_flush_interval": "30s",
  "hashrate_window": "5m",
  "pool_stats_interval": "60s"
}
```

## Monitoring

### Health Checks

```bash
# Check API
curl http://localhost:8001/api/pool/stats

# Check sync status
curl http://localhost:8001/api/pool/sync

# Check specific chain
curl http://localhost:8001/api/chain/dogecoin/stats
```

### Logs

The pool logs important events:
- Block found notifications
- Share validation failures
- RPC connection issues
- Payout processing

## Troubleshooting

### Daemon Not Syncing

Check sync status via API:
```bash
curl http://localhost:8001/api/pool/sync
```

### Shares Rejected

Common causes:
- Incorrect difficulty settings
- Daemon not fully synced
- Network latency
- Stale shares

### Blocks Not Submitting

Check:
- Daemon RPC credentials
- Reward addresses are valid
- Daemon is configured for auxpow (auxiliary chains)
- Coinbase size limits

## Security

- Use strong RPC passwords
- Firewall RPC ports (only allow pool server)
- Use SSL/TLS for public-facing API
- Regular database backups
- Monitor for suspicious mining activity

## License

[Your License Here]

## Support

- GitHub Issues: [Your Repo URL]
- Discord: [Your Discord URL]
- Email: [Your Email]

## Credits

Based on merged mining concepts from:
- Bitcoin merged mining (Namecoin)
- Litecoin/Dogecoin merged mining pools
- Modern pool implementations

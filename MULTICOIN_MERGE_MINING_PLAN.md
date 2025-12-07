# Multicoin Merge Mining Implementation Plan

## Current State Analysis

The pool currently supports **2-chain merged mining** (1 primary + 1 auxiliary):
- **Primary chain**: Litecoin (generates the actual proof of work)
- **Auxiliary chain (Aux1)**: Dogecoin (embedded in primary's coinbase)

### Existing Infrastructure
- ✅ Config already supports `merged_blockchain_order` array with N chains
- ✅ `Pair.AuxBlocks` is already a slice: `[]bitcoin.AuxBlock`
- ✅ Pool server allocates space: `amountOfChains = len(BlockChainOrder) - 1`
- ✅ Helper methods exist: `GetAuxN(n int)` for accessing any aux chain
- ❌ **But**: Only Aux1 is actually fetched, validated, and submitted
- ❌ **But**: AuxPow Merkle branch is hardcoded empty (single chain only)

### Current Limitations
1. `fetchAllBlockTemplatesFromRPC()` only fetches Aux1
2. `validateAndWeighShare()` only checks primary + aux1 (returns: invalid, valid, primary, aux1, dual)
3. `recieveWorkFromClient()` only submits to primary and aux1
4. `AuxMerkleBranch` is hardcoded empty (no Merkle tree)
5. Coinbase embeds single aux block hash directly (no Merkle root)

## Merge Mining Protocol Overview

### Single Auxiliary Chain (Current)
```
Primary Coinbase Script:
  [height] [pool signature] [fabe6d6d] [aux_block_hash] [01000000] [00000000] [00002632]
                             ↑ MM header  ↑ Dogecoin hash  ↑ size=1    ↑ nonce    ↑ chain ID
```

### Multiple Auxiliary Chains (Target)
```
Primary Coinbase Script:
  [height] [pool signature] [fabe6d6d] [merkle_root] [merkle_size] [merkle_nonce] [00002632]
                             ↑ MM header  ↑ Root of     ↑ Tree size   ↑ Position
                                          aux Merkle tree

Auxiliary Merkle Tree (example with 3 aux chains):
                    merkle_root
                   /            \
              H(A,B)              H(C,C)
             /      \            /      \
         Aux1Hash  Aux2Hash  Aux3Hash  Aux3Hash  (duplicated for power-of-2)

Each aux chain gets a Merkle branch proof showing its path to the root.
```

## Implementation Plan

### Phase 1: Core Merkle Tree Infrastructure

#### 1.1 Create Auxiliary Chain Merkle Tree Builder
**File**: `bitcoin/auxmerkle.go` (new file)

```go
type AuxChainMerkleTree struct {
    Root     string              // Merkle root hash
    Size     uint                // Number of auxiliary chains
    Branches map[int]MerklePath  // chainIndex => proof path
}

type MerklePath struct {
    Index    uint      // Position in tree
    Siblings []string  // Sibling hashes for proof
}

func BuildAuxChainMerkleTree(auxBlocks []AuxBlock) AuxChainMerkleTree
func (tree *AuxChainMerkleTree) GetBranchForChain(chainIndex int) MerklePath
```

**Implementation Details**:
- Build balanced binary Merkle tree from aux block hashes
- Duplicate last element if count is not power of 2
- Calculate Merkle root using double SHA256
- Generate proof paths for each auxiliary chain
- Chain index determines position in tree (0, 1, 2, ...)

#### 1.2 Update AuxBlock to Include Merkle Proof
**File**: `bitcoin/auxpow.go`

```go
type AuxBlock struct {
    Hash              string
    ChainID           int
    PreviousBlockHash string
    CoinbaseHash      string
    CoinbaseValue     uint
    Bits              string
    Height            uint64
    Target            string
    MerkleIndex       uint      // NEW: position in aux Merkle tree
    MerkleBranch      []string  // NEW: proof path to Merkle root
}

// Update GetWork() to use Merkle root when multiple chains
func (b *AuxBlock) GetWork(auxMerkleRoot string, auxMerkleSize uint, merkleNonce uint) string
```

#### 1.3 Update AuxMerkleBranch for Multiple Chains
**File**: `bitcoin/auxpow.go`

```go
// Current: hardcoded empty branch
type AuxMerkleBranch struct {
    numberOfBranches string  // "00" = empty
    mask             string  // "00000000"
}

// NEW: support actual Merkle branches
type AuxMerkleBranch struct {
    BranchCount uint      // Number of sibling hashes
    Branches    []string  // Sibling hashes
    Mask        uint      // Bitmask for branch path (left/right)
}

func makeAuxChainMerkleBranch(auxBlock AuxBlock) AuxMerkleBranch {
    // Build from auxBlock.MerkleBranch and MerkleIndex
    return AuxMerkleBranch{
        BranchCount: uint(len(auxBlock.MerkleBranch)),
        Branches:    auxBlock.MerkleBranch,
        Mask:        auxBlock.MerkleIndex,
    }
}

func (am *AuxMerkleBranch) Serialize() string {
    // Format: [count][branch1][branch2]...[mask]
}
```

### Phase 2: Update Work Generation

#### 2.1 Fetch All Auxiliary Block Templates
**File**: `pool/server.go`

```go
// Current: fetchAllBlockTemplatesFromRPC() only gets primary + aux1
// NEW: fetch all auxiliary chains

func (p *PoolServer) fetchAllBlockTemplatesFromRPC() (bitcoin.Template, []bitcoin.AuxBlock, error) {
    // Get primary template
    template := fetchPrimaryTemplate()

    // Get ALL auxiliary blocks
    auxBlocks := make([]bitcoin.AuxBlock, 0)
    for i := 1; i < len(p.config.BlockChainOrder); i++ {
        chainName := p.config.BlockChainOrder[i]
        node := p.activeNodes[chainName]
        auxBlock := fetchAuxBlock(node)
        auxBlocks = append(auxBlocks, auxBlock)
    }

    return template, auxBlocks, nil
}
```

#### 2.2 Build Auxiliary Merkle Tree During Work Generation
**File**: `pool/work.go`

```go
func (p *PoolServer) fetchRpcBlockTemplatesAndCacheWork() error {
    template, auxBlocks, err := p.fetchAllBlockTemplatesFromRPC()

    // Build auxiliary chain Merkle tree
    var auxMerkleTree bitcoin.AuxChainMerkleTree
    if len(auxBlocks) > 0 {
        auxMerkleTree = bitcoin.BuildAuxChainMerkleTree(auxBlocks)

        // Annotate each aux block with its Merkle proof
        for i := range auxBlocks {
            branch := auxMerkleTree.GetBranchForChain(i)
            auxBlocks[i].MerkleIndex = branch.Index
            auxBlocks[i].MerkleBranch = branch.Siblings
        }
    }

    // Build coinbase auxiliary data
    auxillary := p.config.BlockSignature
    if len(auxBlocks) > 0 {
        // Embed Merkle root instead of single hash
        mergedPOW := buildMergedMiningData(auxMerkleTree)
        auxillary = auxillary + hexStringToByteString(mergedPOW)

        p.templates.AuxBlocks = auxBlocks
    }

    // Generate work...
}

func buildMergedMiningData(tree bitcoin.AuxChainMerkleTree) string {
    // Format: [fabe6d6d][merkle_root][size][nonce][chainID]
    return "fabe6d6d" + tree.Root +
           varUint(tree.Size) +
           "00000000" + // nonce (can be used for positioning)
           "00002632"   // primary chain ID (Litecoin)
}
```

### Phase 3: Update Share Validation

#### 3.1 Validate Against All Chains
**File**: `pool/share.go`

```go
// Current: only supports primary, aux1, dual
// NEW: support N auxiliary chains

const (
    shareInvalid = iota
    shareValid
    blockCandidate  // Generic block candidate
)

type BlockCandidateResult struct {
    IsCandidate     bool
    PrimaryMeetsTarget bool
    AuxChainsMetTargets []int  // Indices of aux chains meeting target
    ShareDifficulty float64
}

func validateAndWeighShare(
    primary *bitcoin.BitcoinBlock,
    auxBlocks []bitcoin.AuxBlock,
    poolDifficulty float64,
) BlockCandidateResult {

    primarySum, _ := primary.Sum()
    poolTarget, _ := bitcoin.TargetFromDifficulty(poolDifficulty / primary.ShareMultiplier())
    shareDifficulty, _ := poolTarget.ToDifficulty()

    result := BlockCandidateResult{
        ShareDifficulty: shareDifficulty,
        AuxChainsMetTargets: make([]int, 0),
    }

    // Check if valid share
    poolTargetBig, _ := poolTarget.ToBig()
    if primarySum.Cmp(poolTargetBig) > 0 {
        result.IsCandidate = false
        return result
    }

    result.IsCandidate = true

    // Check primary chain target
    primaryTarget := bitcoin.Target(primary.Template.Target)
    primaryTargetBig, _ := primaryTarget.ToBig()
    if primarySum.Cmp(primaryTargetBig) <= 0 {
        result.PrimaryMeetsTarget = true
    }

    // Check ALL auxiliary chain targets
    for i, auxBlock := range auxBlocks {
        auxTarget := bitcoin.Target(reverseHexBytes(auxBlock.Target))
        auxTargetBig, _ := auxTarget.ToBig()

        if primarySum.Cmp(auxTargetBig) <= 0 {
            result.AuxChainsMetTargets = append(result.AuxChainsMetTargets, i)
        }
    }

    return result
}
```

### Phase 4: Update Block Submission

#### 4.1 Submit to All Qualifying Chains
**File**: `pool/work.go`

```go
func (p *PoolServer) recieveWorkFromClient(share bitcoin.Work, client *stratumClient) error {
    // ... existing validation code ...

    result := validateAndWeighShare(&primaryBlockTemplate, p.templates.AuxBlocks, p.config.PoolDifficulty)

    if !result.IsCandidate {
        return errors.New("Invalid share")
    }

    // Submit to all qualifying auxiliary chains
    for _, auxIndex := range result.AuxChainsMetTargets {
        auxBlock := p.templates.AuxBlocks[auxIndex]
        chainName := p.config.BlockChainOrder[auxIndex + 1] // +1 because primary is index 0

        err = p.submitAuxBlockForChain(primaryBlockTemplate, auxBlock, chainName)
        if err != nil {
            log.Printf("Failed to submit aux block for %s: %v", chainName, err)
            // Try failover node
            p.rpcManagers[chainName].CheckAndRecoverRPCs()
            err = p.submitAuxBlockForChain(primaryBlockTemplate, auxBlock, chainName)
        }

        if err == nil {
            // Record successful block
            p.recordFoundBlock(auxBlock, chainName, minerAddress)
        }
    }

    // Submit primary if it meets target
    if result.PrimaryMeetsTarget {
        err = p.submitBlockToChain(primaryBlockTemplate)
        if err == nil {
            p.recordFoundBlock(primaryBlockTemplate, p.config.GetPrimary(), minerAddress)
        }
    }

    return nil
}
```

#### 4.2 Create Chain-Specific Aux Submission
**File**: `pool/nodes.go`

```go
func (p *PoolServer) submitAuxBlockForChain(
    primaryBlock bitcoin.BitcoinBlock,
    auxBlock bitcoin.AuxBlock,
    chainName string,
) error {
    // Build AuxPow with proper Merkle branch for this specific chain
    auxpow := bitcoin.MakeAuxPow(primaryBlock, auxBlock)

    node := p.activeNodes[chainName]
    success, err := node.RPC.SubmitAuxBlock(auxBlock.Hash, auxpow.Serialize())

    if !success {
        m := "⚠️  %v node failed to submit aux block: %v"
        m = fmt.Sprintf(m, chainName, err.Error())
        return errors.New(m)
    }

    return err
}
```

#### 4.3 Update AuxPow Generation
**File**: `bitcoin/auxpow.go`

```go
// Current: MakeAuxPow(parentBlock)
// NEW: include specific aux block for proper Merkle branch

func MakeAuxPow(parentBlock BitcoinBlock, auxBlock AuxBlock) AuxPow {
    if parentBlock.hash == "" {
        panic("Set parent block hash first")
    }

    return AuxPow{
        ParentCoinbase:       parentBlock.coinbase,
        ParentHeaderHash:     parentBlock.hash,
        ParentMerkleBranch:   makeParentMerkleBranch(parentBlock.merkleSteps),
        auxMerkleBranch:      makeAuxChainMerkleBranch(auxBlock),  // Use aux block's proof
        ParentHeaderUnhashed: parentBlock.header,
    }
}
```

### Phase 5: Configuration & Helper Updates

#### 5.1 Update Config Helpers
**File**: `config/config.go`

```go
// Add helper methods for accessing any auxiliary chain
func (b BlockChainOrder) GetAuxN(n int) string {
    index := n + 1  // +1 because primary is at index 0
    if index >= len(b) {
        return ""
    }
    return b[index]
}

func (b BlockChainOrder) GetAuxChains() []string {
    if len(b) <= 1 {
        return []string{}
    }
    return b[1:]  // All chains except primary
}

func (b BlockChainOrder) AuxChainCount() int {
    return len(b) - 1
}
```

#### 5.2 Update Pool Server Helpers
**File**: `pool/nodes.go`

```go
func (p *PoolServer) GetAuxNode(chainName string) blockChainNode {
    return p.activeNodes[chainName]
}

func (p *PoolServer) GetAuxNodeByIndex(n int) blockChainNode {
    chainName := p.config.GetAuxN(n)
    return p.activeNodes[chainName]
}
```

### Phase 6: Database & Persistence Updates

#### 6.1 Update Found Blocks Recording
**File**: `pool/work.go`

```go
func (p *PoolServer) recordFoundBlock(
    auxBlock bitcoin.AuxBlock,
    chainName string,
    minerAddress string,
) error {
    aux1Target := bitcoin.Target(reverseHexBytes(auxBlock.Target))
    auxDifficulty, _ := aux1Target.ToDifficulty()
    auxDifficulty = auxDifficulty * bitcoin.GetChain(chainName).ShareMultiplier()

    found := persistence.Found{
        PoolID:               p.config.PoolName,
        Chain:                chainName,
        Created:              time.Now(),
        Hash:                 auxBlock.Hash,
        NetworkDifficulty:    auxDifficulty,
        BlockHeight:          uint(auxBlock.Height),
        TransactionConfirmationData: reverseHexBytes(auxBlock.CoinbaseHash),
        Status:               persistence.StatusPending,
        Type:                 "Auxiliary",
        ConfirmationProgress: 0,
        Miner:                minerAddress,
    }

    return persistence.Blocks.Insert(found)
}
```

### Phase 7: Testing & Validation

#### 7.1 Unit Tests
Create tests for:
- Merkle tree building with 1, 2, 3, 4, 5+ auxiliary chains
- Merkle branch proof generation and verification
- AuxPow serialization with proper Merkle branches
- Share validation against multiple targets

#### 7.2 Integration Testing
1. Test with 3-chain setup (LTC + DOGE + another)
2. Test with 4+ chains
3. Test failover when auxiliary chain is unavailable
4. Test submission when multiple aux chains meet target simultaneously

## Migration Strategy

### Backward Compatibility
- ✅ Single auxiliary chain (current 2-chain setup) continues to work
- ✅ Config format unchanged (just add more chains to `merged_blockchain_order`)
- ✅ Existing database schema compatible

### Deployment Steps
1. Deploy code changes
2. Update config to add new auxiliary chains
3. Restart pool server
4. Monitor logs for successful aux block fetching
5. Verify blocks submitted to all chains

## Example Configuration

```json
{
    "merged_blockchain_order": [
        "litecoin",      // Primary (index 0)
        "dogecoin",      // Aux1 (index 0 in AuxBlocks array)
        "viacoin",       // Aux2 (index 1 in AuxBlocks array)
        "myriadcoin",    // Aux3 (index 2 in AuxBlocks array)
        "namecoin"       // Aux4 (index 3 in AuxBlocks array)
    ],
    "blockchains": {
        "litecoin": [...],
        "dogecoin": [...],
        "viacoin": [...],
        "myriadcoin": [...],
        "namecoin": [...]
    },
    "payouts": {
        "chains": {
            "litecoin": {...},
            "dogecoin": {...},
            "viacoin": {...},
            "myriadcoin": {...},
            "namecoin": {...}
        }
    }
}
```

## Key Architectural Decisions

1. **Merkle Tree Position**: Chain position determined by index in `merged_blockchain_order`
2. **Merkle Tree Type**: Balanced binary tree with duplicate last element for non-power-of-2
3. **Hash Function**: Double SHA256 (Bitcoin standard)
4. **Backward Compatibility**: Single aux chain uses same format (tree of size 1)
5. **Submission Strategy**: Submit to ALL chains meeting target (parallel submissions)
6. **Failover**: Per-chain RPC failover (existing infrastructure)

## Performance Considerations

- **Merkle tree building**: O(n) where n = number of aux chains (minimal overhead)
- **RPC calls**: Parallel fetching of aux blocks (non-blocking)
- **Validation**: Check all targets in single pass (O(n) additional checks)
- **Submission**: Parallel submissions to qualifying chains

## Success Criteria

✅ Pool successfully fetches block templates from 3+ chains
✅ Primary chain coinbase contains auxiliary Merkle root
✅ Each auxiliary chain receives proper Merkle branch proof in AuxPow
✅ Blocks accepted by all configured auxiliary chains
✅ Database tracks blocks found for all chains
✅ Payouts distributed for all chains according to config
✅ Existing 2-chain setup continues to work without changes

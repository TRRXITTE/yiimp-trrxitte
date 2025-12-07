package bitcoin

import "fmt"

const (
	mergedMiningHeader  = "fabe6d6d"
	mergedMiningTrailer = "010000000000000000002632"
)

type AuxBlock struct {
	Hash              string `json:"hash"`
	ChainID           int    `json:"chainid"`
	PreviousBlockHash string `json:"previousblockhash"`
	CoinbaseHash      string `json:"coinbasehash"`
	CoinbaseValue     uint   `json:"coinbasevalue"`
	Bits              string `json:"bits"`
	Height            uint64 `json:"height"`
	Target            string `json:"target"`
	MerkleIndex       uint
	MerkleBranch      []string
}

func (b *AuxBlock) GetWork() string {
	return mergedMiningHeader + b.Hash + mergedMiningTrailer
}

func (b *AuxBlock) GetWorkWithMerkleRoot(merkleRoot string, merkleSize uint) string {
	sizeHex := fmt.Sprintf("%08x", merkleSize)
	return mergedMiningHeader + merkleRoot + sizeHex + "00000000" + "00002632"
}

type AuxPow struct {
	ParentCoinbase   string
	ParentHeaderHash string
	ParentMerkleBranch
	auxMerkleBranch      AuxMerkleBranch
	auxMerkleBranches    []string
	ParentHeaderUnhashed string
}

func MakeAuxPow(parentBlock BitcoinBlock) AuxPow {
	if parentBlock.hash == "" {
		panic("Set parent block hash first")
	}
	// debugAuxPow(parentBlock, makeParentMerkleBranch(parentBlock.merkleSteps), makeAuxChainMerkleBranch())

	return AuxPow{
		ParentCoinbase:       parentBlock.coinbase,
		ParentHeaderHash:     parentBlock.hash,
		ParentMerkleBranch:   makeParentMerkleBranch(parentBlock.merkleSteps),
		auxMerkleBranch:      makeAuxChainMerkleBranch(),
		ParentHeaderUnhashed: parentBlock.header,
	}
}

func MakeAuxPowWithBranch(parentBlock BitcoinBlock, auxBlock AuxBlock) AuxPow {
	if parentBlock.hash == "" {
		panic("Set parent block hash first")
	}

	return AuxPow{
		ParentCoinbase:       parentBlock.coinbase,
		ParentHeaderHash:     parentBlock.hash,
		ParentMerkleBranch:   makeParentMerkleBranch(parentBlock.merkleSteps),
		auxMerkleBranch:      makeAuxChainMerkleBranchFromBlock(auxBlock),
		auxMerkleBranches:    auxBlock.MerkleBranch,
		ParentHeaderUnhashed: parentBlock.header,
	}
}

func (p *AuxPow) Serialize() string {
	auxMerkleSerialized := p.auxMerkleBranch.Serialize()
	if len(p.auxMerkleBranches) > 0 {
		auxMerkleSerialized = p.auxMerkleBranch.SerializeWithBranches(p.auxMerkleBranches)
	}

	return p.ParentCoinbase +
		p.ParentHeaderHash +
		p.ParentMerkleBranch.Serialize() +
		auxMerkleSerialized +
		p.ParentHeaderUnhashed
}

type ParentMerkleBranch struct {
	Length uint
	Items  []string
	mask   string
}

func makeParentMerkleBranch(items []string) ParentMerkleBranch {
	length := uint(len(items))
	return ParentMerkleBranch{
		Length: length,
		Items:  items,
		mask:   "00000000",
	}
}

func (pm *ParentMerkleBranch) Serialize() string {
	items := ""
	for _, item := range pm.Items {
		items = items + item
	}
	return varUint(pm.Length) + items + pm.mask
}

type AuxMerkleBranch struct {
	numberOfBranches string
	mask             string
}

func makeAuxChainMerkleBranch() AuxMerkleBranch {
	return AuxMerkleBranch{
		numberOfBranches: "00",
		mask:             "00000000",
	}
}

func makeAuxChainMerkleBranchFromBlock(auxBlock AuxBlock) AuxMerkleBranch {
	if len(auxBlock.MerkleBranch) == 0 {
		return makeAuxChainMerkleBranch()
	}

	branchCount := uint(len(auxBlock.MerkleBranch))
	mask := fmt.Sprintf("%08x", auxBlock.MerkleIndex)

	return AuxMerkleBranch{
		numberOfBranches: varUint(branchCount),
		mask:             mask,
	}
}

func (am *AuxMerkleBranch) Serialize() string {
	return am.numberOfBranches + am.mask
}

func (am *AuxMerkleBranch) SerializeWithBranches(branches []string) string {
	branchesHex := ""
	for _, branch := range branches {
		branchesHex += branch
	}
	return am.numberOfBranches + branchesHex + am.mask
}

func debugAuxPow(parentBlock BitcoinBlock, parentMerkle ParentMerkleBranch, auxchainMerkle AuxMerkleBranch) {
	fmt.Println()
	fmt.Println("coinbase", parentBlock.coinbase)
	fmt.Println("hash", parentBlock.hash)
	fmt.Println("merkleSteps", parentBlock.merkleSteps)
	fmt.Println("merkleDigested", parentMerkle.Serialize())
	fmt.Println("chainmerklebranch", auxchainMerkle.Serialize())
	fmt.Println("header", parentBlock.header)
	fmt.Println()
}

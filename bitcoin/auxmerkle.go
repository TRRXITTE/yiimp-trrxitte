package bitcoin

import (
	"encoding/hex"
)

type MerklePath struct {
	Index    uint
	Siblings []string
}

type AuxChainMerkleTree struct {
	Root     string
	Size     uint
	Branches map[int]MerklePath
}

func BuildAuxChainMerkleTree(auxBlocks []AuxBlock) AuxChainMerkleTree {
	if len(auxBlocks) == 0 {
		return AuxChainMerkleTree{
			Root:     "",
			Size:     0,
			Branches: make(map[int]MerklePath),
		}
	}

	if len(auxBlocks) == 1 {
		tree := AuxChainMerkleTree{
			Root:     auxBlocks[0].Hash,
			Size:     1,
			Branches: make(map[int]MerklePath),
		}
		tree.Branches[0] = MerklePath{
			Index:    0,
			Siblings: []string{},
		}
		return tree
	}

	hashes := make([]string, len(auxBlocks))
	for i, block := range auxBlocks {
		hashes[i] = block.Hash
	}

	tree := AuxChainMerkleTree{
		Size:     uint(len(auxBlocks)),
		Branches: make(map[int]MerklePath),
	}

	levelHashes := padToNextPowerOfTwo(hashes)
	originalCount := len(hashes)

	tree.Root = buildMerkleRoot(levelHashes)

	for i := 0; i < originalCount; i++ {
		tree.Branches[i] = generateMerklePath(i, levelHashes)
	}

	return tree
}

func padToNextPowerOfTwo(hashes []string) []string {
	length := len(hashes)
	if isPowerOfTwo(length) {
		return hashes
	}

	nextPowerOfTwo := 1
	for nextPowerOfTwo < length {
		nextPowerOfTwo <<= 1
	}

	paddedHashes := make([]string, nextPowerOfTwo)
	copy(paddedHashes, hashes)

	for i := length; i < nextPowerOfTwo; i++ {
		paddedHashes[i] = hashes[length-1]
	}

	return paddedHashes
}

func isPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

func buildMerkleRoot(hashes []string) string {
	if len(hashes) == 0 {
		return ""
	}

	if len(hashes) == 1 {
		return hashes[0]
	}

	currentLevel := hashes

	for len(currentLevel) > 1 {
		nextLevel := make([]string, len(currentLevel)/2)

		for i := 0; i < len(currentLevel); i += 2 {
			combinedHash := combineHashes(currentLevel[i], currentLevel[i+1])
			nextLevel[i/2] = combinedHash
		}

		currentLevel = nextLevel
	}

	return currentLevel[0]
}

func generateMerklePath(index int, hashes []string) MerklePath {
	path := MerklePath{
		Index:    uint(index),
		Siblings: make([]string, 0),
	}

	currentIndex := index
	currentLevel := hashes

	for len(currentLevel) > 1 {
		siblingIndex := currentIndex ^ 1

		if siblingIndex < len(currentLevel) {
			path.Siblings = append(path.Siblings, currentLevel[siblingIndex])
		}

		nextLevel := make([]string, len(currentLevel)/2)
		for i := 0; i < len(currentLevel); i += 2 {
			nextLevel[i/2] = combineHashes(currentLevel[i], currentLevel[i+1])
		}

		currentIndex /= 2
		currentLevel = nextLevel
	}

	return path
}

func combineHashes(left, right string) string {
	leftBytes, err := hex.DecodeString(left)
	if err != nil {
		panic("Invalid hex in left hash: " + err.Error())
	}

	rightBytes, err := hex.DecodeString(right)
	if err != nil {
		panic("Invalid hex in right hash: " + err.Error())
	}

	combined := append(leftBytes, rightBytes...)
	hashed := doubleSha256Bytes(combined)

	return hex.EncodeToString(hashed[:])
}

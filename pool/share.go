package pool

import (
	"designs.capital/dogepool/bitcoin"
)

const (
	shareInvalid = iota
	shareValid
	blockCandidate
)

type BlockCandidateResult struct {
	Status              int
	PrimaryMeetsTarget  bool
	AuxChainsMetTargets []int
	ShareDifficulty     float64
}

func validateAndWeighShare(primary *bitcoin.BitcoinBlock, auxBlocks []bitcoin.AuxBlock, poolDifficulty float64) BlockCandidateResult {
	result := BlockCandidateResult{
		Status:              shareInvalid,
		PrimaryMeetsTarget:  false,
		AuxChainsMetTargets: make([]int, 0),
		ShareDifficulty:     0,
	}

	primarySum, err := primary.Sum()
	logOnError(err)

	poolTarget, _ := bitcoin.TargetFromDifficulty(poolDifficulty / primary.ShareMultiplier())
	shareDifficulty, _ := poolTarget.ToDifficulty()
	result.ShareDifficulty = shareDifficulty

	poolTargetBig, _ := poolTarget.ToBig()
	if primarySum.Cmp(poolTargetBig) > 0 {
		return result
	}

	result.Status = shareValid

	primaryTarget := bitcoin.Target(primary.Template.Target)
	primaryTargetBig, _ := primaryTarget.ToBig()

	if primarySum.Cmp(primaryTargetBig) <= 0 {
		result.PrimaryMeetsTarget = true
		result.Status = blockCandidate
	}

	for i, auxBlock := range auxBlocks {
		if auxBlock.Hash == "" {
			continue
		}

		auxTarget := bitcoin.Target(reverseHexBytes(auxBlock.Target))
		auxTargetBig, _ := auxTarget.ToBig()

		if primarySum.Cmp(auxTargetBig) <= 0 {
			result.AuxChainsMetTargets = append(result.AuxChainsMetTargets, i)
			result.Status = blockCandidate
		}
	}

	return result
}

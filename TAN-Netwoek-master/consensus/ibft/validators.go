package ibft

import (
	"math"

	"github.com/tarality/tan-network/types"
	"github.com/tarality/tan-network/validators"
)

func CalcMaxFaultyNodes(s validators.Validators) int {
	// N -> number of nodes in the validator set
	// F -> number of faulty nodes
	//
	// The protocol tolerates 40% of the nodes being faulty.
	// F = floor(0.4 * N)
	//
	// For example:
	// If N = 10, F = floor(0.4 * 10) = floor(4) = 4
	// If N = 15, F = floor(0.4 * 15) = floor(6) = 6
	return int(0.4 * float64(s.Len()))
}

type QuorumImplementation func(validators.Validators) int

// LegacyQuorumSize returns the legacy quorum size for the given validator set
func LegacyQuorumSize(set validators.Validators) int {
	// According to the IBFT spec, the number of valid messages
	// 60% agreement required
	return int(math.Ceil(0.6 * float64(set.Len())))
}

// TODO
func OptimalQuorumSize(set validators.Validators) int {
	quorumPercentage := 0.6 // 60% agreement required
	quorumSize := int(math.Ceil(quorumPercentage * float64(set.Len())))
	return quorumSize
}

// OptimalQuorumSize returns the optimal quorum size for the given validator set
// func OptimalQuorumSize(set validators.Validators) int {
// 	//	if the number of validators is less than 4,
// 	//	then the entire set is required
// 	// if CalcMaxFaultyNodes(set) == 0 {
// 	// 	/*
// 	// 		N: 1 -> Q: 1
// 	// 		N: 2 -> Q: 2
// 	// 		N: 3 -> Q: 3
// 	// 	*/
// 	// 	return set.Len()
// 	// }

// 	// 0.6 is not working
// 	// below 0.5 is working fine.

// 	// (quorum optimal)	Q = ceil(2/3 * N)
// 	// return int(math.Ceil(2 * float64(set.Len()) / 3))

// 	// Multiply totalVotingPower by 6 and then divide by 10
// 	quorum := (set.Len() * 6) / 10

// 	// Check if there is a remainder when dividing by 1
// 	remainder := (set.Len() * 6) % 10

// 	// If remainder is greater than 0, add 1 to the quorum
// 	if remainder > 0 {
// 		return quorum + 1
// 	}

// 	// If no remainder, return the quorum as it is
// 	fmt.Println("---------line no quorum", quorum, "-----remainder--", remainder)
// 	return quorum
// }

// }

func CalcProposer(
	validators validators.Validators,
	round uint64,
	lastProposer types.Address,
) validators.Validator {
	var seed uint64

	if lastProposer == types.ZeroAddress {
		seed = round
	} else {
		offset := int64(0)

		if index := validators.Index(lastProposer); index != -1 {
			offset = index
		}

		seed = uint64(offset) + round + 1
	}

	pick := seed % uint64(validators.Len())

	return validators.At(pick)
}

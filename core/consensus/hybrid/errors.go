package hybrid

import "errors"

var (
    ErrInvalidDifficulty   = errors.New("invalid proof of work difficulty")
    ErrInsufficientStake   = errors.New("validator has insufficient stake")
    ErrInsufficientVotes   = errors.New("insufficient votes for block approval")
    ErrNoValidators        = errors.New("no eligible validators available")
    ErrBlockTimeTooEarly   = errors.New("block time is too early")
    ErrInvalidBlockNumber  = errors.New("invalid block number")
)

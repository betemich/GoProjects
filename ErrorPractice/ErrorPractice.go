package main

import "fmt"

type RankedPlayer struct {
	Nickname string
	Rank     string
	Hours    int32
}

type RankedError struct {
	Err            error
	ErrCode        int16
	ErrDescruption string
}

func (re *RankedError) Error() string {
	return fmt.Sprintf("%s\nCode: %d\n%s", re.Err.Error(), re.ErrCode, re.ErrDescruption)
}

func (re *RankedError) PrintError() {
	fmt.Println(re.Error())
}

func CheckRanks() func(val string) bool {
	ranks := []string{"Bronze", "Silver", "Gold", "Platinum", "Diamond", "Master", "Predator"}
	return func(val string) bool {
		for _, v := range ranks {
			if v == val {
				return true
			}
		}
		return false
	}
}

func RankedProcessing(player *RankedPlayer) *RankedError {
	ranks_checker := CheckRanks()
	if !ranks_checker(player.Rank) {
		return &RankedError{
			Err:            fmt.Errorf("invalid rank error"),
			ErrCode:        52,
			ErrDescruption: fmt.Sprintf("Player %s has invalid rank %s", player.Nickname, player.Rank),
		}
	}
	return nil
}

func main() {
	player := RankedPlayer{"average 483 enjoyer", "Malachite", 4000}
	err := RankedProcessing(&player)
	if err != nil {
		err.PrintError()
		return
	}
	fmt.Println("OK")
}

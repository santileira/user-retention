package calculator

import (
	"fmt"
	"strconv"
	"time"
)

type UserRetentionServiceImpl struct{}

func (u *UserRetentionServiceImpl) Calculate(rawUsersActivity [][]string) (map[int][]int, error) {

	userRetention := make(map[int][]int, 14)
	for i := 0; i < 14; i++ {
		userRetention[i] = make([]int, 14)
	}

	usersActivityByUserID := make(map[string][]int, 0)
	for _, rawUserActivity := range rawUsersActivity {
		dayInUTC, err := u.getDayFromRawUserActivity(rawUserActivity)
		if err != nil {
			return nil, err
		}
		userID := rawUserActivity[1]

		userActivity := usersActivityByUserID[userID]
		lenUserActivity := len(userActivity)

		// if the user doesn't have activity or user has consecutive usage,
		//add the day as activity and follow with the next row.
		if lenUserActivity == 0 || userActivity[lenUserActivity-1] == dayInUTC-1 {
			usersActivityByUserID[userID] = append(userActivity, dayInUTC)
			continue
		}

		// if user has more than one event from the same day, follow with the next row
		if userActivity[lenUserActivity-1] == dayInUTC {
			continue
		}

		// if the user has not consecutive usage,
		// update the user retention for the first day and
		// clean the old user activity.
		u.updateUserRetention(userRetention, userActivity)
		usersActivityByUserID[userID] = []int{dayInUTC}
	}

	// update the user retention for the users that didn't register yet
	for _, userActivity := range usersActivityByUserID {
		u.updateUserRetention(userRetention, userActivity)
	}

	fmt.Println(userRetention)
	return userRetention, nil
}

func (u *UserRetentionServiceImpl) getDayFromRawUserActivity(rawUserActivity []string) (int, error) {
	timestampStr := rawUserActivity[0]
	timestampInUTC, err := strconv.Atoi(timestampStr)
	if err != nil {
		return 0, fmt.Errorf("error converting timestamp %s, err: %s", timestampStr, err.Error())
	}
	dayInUTC := time.Unix(int64(timestampInUTC), 0).UTC().Day()
	return dayInUTC, nil
}

func (u *UserRetentionServiceImpl) updateUserRetention(userRetention map[int][]int, userActivity []int) {
	// it uses userActivity[0] - 1 (ex: 1 - 1 = 0) to have the data of the day 1 in the position 0.
	// it uses len(userActivity) -1 (ex: 14 - 1 = 13) to have the data of the day 14 in the position 13.
	userRetention[userActivity[0]-1][len(userActivity)-1] += 1
}

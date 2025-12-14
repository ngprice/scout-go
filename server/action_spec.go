package server

const (
	ActionScout ActionType = iota
	ActionScoutReverse
	ActionShow
	ActionScoutAndShow
	ActionScoutAndShowReverse
	ActionReverseHand
)

type ActionType int

type ActionSpec struct {
	ID                            int
	Type                          ActionType
	ScoutTakeIndex, ScoutPutIndex int // scout action params
	ShowFirstIndex, ShowLength    int // show action params
}

func getAllActions() []ActionSpec {
	actions := make([]ActionSpec, 0)
	scoutActions := make([]ActionSpec, 0)
	scoutActionsReverse := make([]ActionSpec, 0)
	showActions := make([]ActionSpec, 0)

	id := 0

	// 1. scout actions
	for take := 0; take < MAX_ACTIVE_SET_SIZE; take++ {
		for put := 0; put <= MAX_HAND_SIZE; put++ {
			scoutAction := ActionSpec{
				ID: id, Type: ActionScout, ScoutTakeIndex: take, ScoutPutIndex: put,
			}
			actions = append(actions, scoutAction)
			scoutActions = append(scoutActions, scoutAction)
			id++
		}
	}

	for take := 0; take < MAX_ACTIVE_SET_SIZE; take++ {
		for put := 0; put <= MAX_HAND_SIZE; put++ {
			scoutAction := ActionSpec{
				ID: id, Type: ActionScoutReverse, ScoutTakeIndex: take, ScoutPutIndex: put,
			}
			actions = append(actions, scoutAction)
			scoutActionsReverse = append(scoutActionsReverse, scoutAction)
			id++
		}
	}

	// 2. show actions
	for start := 0; start < MAX_HAND_SIZE; start++ {
		for length := 1; length <= MAX_HAND_SIZE-start; length++ {
			showAction := ActionSpec{
				ID: id, Type: ActionShow, ShowFirstIndex: start, ShowLength: length,
			}
			actions = append(actions, showAction)
			showActions = append(showActions, showAction)
			id++
		}
	}

	// 3. scoutandshow = Cartesian product
	for _, scout := range scoutActions {
		for _, show := range showActions {
			actions = append(actions, ActionSpec{
				ID:             id,
				Type:           ActionScoutAndShow,
				ScoutTakeIndex: scout.ScoutTakeIndex, ScoutPutIndex: scout.ScoutPutIndex,
				ShowFirstIndex: show.ShowFirstIndex, ShowLength: show.ShowLength,
			})
			id++
		}
	}

	// 4. scoutandshowreverse = Cartesian product
	for _, scout := range scoutActionsReverse {
		for _, show := range showActions {
			actions = append(actions, ActionSpec{
				ID:             id,
				Type:           ActionScoutAndShowReverse,
				ScoutTakeIndex: scout.ScoutTakeIndex, ScoutPutIndex: scout.ScoutPutIndex,
				ShowFirstIndex: show.ShowFirstIndex, ShowLength: show.ShowLength,
			})
			id++
		}
	}

	// 5. reversehand
	actions = append(actions, ActionSpec{
		ID: id, Type: ActionReverseHand,
	})

	return actions
}

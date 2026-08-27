// Activity selection
// General layout split: Activity selection area & Activity interaction area

package main

import "net/http"

// (Get+Post Query Response) Interaction activity
type Activity interface {
	GetID() string
	GetDesc() string
	Respond(w http.ResponseWriter, r *http.Request)
}

// All available activities
var activities []Activity = []Activity{
	sylWordActivity{},
	sylParaActivity{},
	decompActivity{},
	verbsActivity{},
	miscActivity{},
	numActivity{},
	kuralActivity{},
}

func getDefaultActivity() Activity { return activities[0] }

func selectActivityById(id string) Activity {
	for _, a := range activities {
		if a.GetID() == id {
			return a
		}
	}
	return getDefaultActivity()
}

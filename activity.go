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
}

// Currently active activity
var currentActivity = activities[0]

func getActivity() Activity { return currentActivity }

func selectActivityById(id string) {
	for _, a := range activities {
		if a.GetID() == id {
			currentActivity = a
			return
		}
	}
}

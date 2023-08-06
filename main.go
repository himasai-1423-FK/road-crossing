package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/mux"
)

type RequestData struct {
	Position   int   `json:"position"`
	Speeds     []int `json:"speeds"`
	EntryTimes []int `json:"entryTimes"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/input", requestProcessor).Methods(http.MethodPost)
	log.Println("Listening on port 8080!")

	http.ListenAndServe(":8080", router)
}

func requestProcessor(w http.ResponseWriter, r *http.Request) {
	var body RequestData
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Unable to parse JSON request", http.StatusBadRequest)
		return
	}

	if len(body.Speeds) != len(body.EntryTimes) {
		http.Error(w, "Required VechileSpeeds or VechileTimes are missing", http.StatusBadRequest)
		return
	}

	prevCrossTimeVehicle := 0
	var maxTime int = 0
	curr := false
	res := 0
	for i := 0; i < len(body.Speeds); i++ {
		var currCrossTimeVehicle int = int(math.Ceil(float64(body.Position)/float64(body.Speeds[i]))) + body.EntryTimes[i]
		fmt.Printf("%v %v %v \n", currCrossTimeVehicle, prevCrossTimeVehicle, i)

		if res == 0 && currCrossTimeVehicle-prevCrossTimeVehicle >= 2 {
			fmt.Printf("Entered at %v\n", i)
			res = prevCrossTimeVehicle
			curr = true
		}
		fmt.Printf("vals res %v %v %v\n", res, currCrossTimeVehicle, i)
		if res != 0 && currCrossTimeVehicle-res <= 2 && !curr {
			fmt.Printf("Entered at res %v %v %v\n", res, currCrossTimeVehicle, i)
			res = 0
		}
		prevCrossTimeVehicle = currCrossTimeVehicle
		curr = false

		if prevCrossTimeVehicle > maxTime {
			maxTime = prevCrossTimeVehicle
		}
	}

	if res == 0 {
		res = maxTime
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(res)
}

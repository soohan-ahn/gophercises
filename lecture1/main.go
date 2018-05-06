package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

func trimQuestion(raw string) string {
	return strings.TrimSpace(strings.TrimFunc(raw, func(r rune) bool {
		return !unicode.IsNumber(r)
	}))
}

func main() {
	file, err := os.Open("problems.csv")
	if err != nil {
		log.Print("Err: ", err)
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Print("Err: ", err)
	}

	reader := bufio.NewReader(os.Stdin)
	if err != nil {
		log.Print("Err: ", err)
	}
	log.Printf("Press the Enter key to start quiz.")
	_, err = reader.ReadString('\n')
	if err != nil {
		log.Print("Err: ", err)
	}

	timer := time.NewTimer(30 * time.Second)

	correctCnt := 0
	for i, _ := range records {
		responseCh := make(chan string)
		go func() {

			question := records[i][0]
			log.Printf("Q: %s\n", question)
			if err != nil {
				log.Print("Err: ", err)
			}

			response, err := reader.ReadString('\n')
			if err != nil {
				log.Print("Err: ", err)
			}
			responseCh <- response
		}()
		answer := trimQuestion(records[i][1])
		quitQuiz := false
		select {
		case <-timer.C:
			log.Print("Time expired!\n")
			quitQuiz = true
			break
		case response := <-responseCh:
			response = strings.TrimSpace(response)
			if strings.Compare(answer, response) == 0 {
				correctCnt++
			}
		}
		if quitQuiz {
			break
		}
	}
	stop := timer.Stop()
	if stop {
		log.Println("Timer stopped")
	}

	log.Print("Correct: ", correctCnt)
	log.Print("Total: ", len(records))
}

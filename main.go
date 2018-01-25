package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

// RandomInt returns a random int between high and low
func RandomInt(low, high int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	timeout := r.Int31n(int32(high - low))
	return int(timeout) + low
}

type Server struct {
	Self            string
	State           State
	ElectionTimeout int // in ms
	Peers           []string
	ProposedLeader  string // url of leader
	ConfirmedLeader string
	mu              sync.RWMutex // guards server
}

func New(self string, peers []string) *Server {
	response := &Server{
		Self:            self,
		State:           Follower,
		ElectionTimeout: RandomInt(150, 300),
		Peers:           peers,
	}
	return response
}

func (s *Server) ElectMe() {
	s.mu.Lock()
	defer s.mu.Unlock()
	votes := []bool{}
	for _, peer := range s.Peers {
		url := fmt.Sprintf("%s/vote?leader=%s", peer, s.Self)
		response, err := http.Get(url)
		if err != nil {
			log.Print(err)
			return
		}
		if response.StatusCode == 200 {
			votes = append(votes, true)
			log.Printf("%s voted for me to be leader", peer)
		} else {
			log.Printf("%s didn't vote for me to be leader", peer)
		}
	}
	if len(votes) >= len(s.Peers)/2 {
		// I am the candidate
		s.State = Candidate
		log.Printf("I am the candidate!")

		confirmVotes := []bool{}

		for _, peer := range s.Peers {
			url := fmt.Sprintf("%s/confirm?leader=%s", peer, s.Self)
			response, err := http.Get(url)
			if err != nil {
				log.Print(err)
				return
			}
			if response.StatusCode == 200 {
				confirmVotes = append(confirmVotes, true)
				log.Printf("%s confirmed me as leader", peer)
			} else {
				log.Printf("%s didn't confirm me as leader", peer)
			}
		}
		if len(confirmVotes) >= len(s.Peers)/2 {
			log.Printf("I am the Leader!")
			s.State = Leader
		}

	}
}

func (s *Server) EventLoop() {
	tick := 0
	for range time.Tick(time.Millisecond * time.Duration(s.ElectionTimeout)) {

		log.Print("Heartbeat")

		// No leader elected yet
		if s.ProposedLeader == "" && s.ConfirmedLeader == "" && s.State == Follower && tick >= 1 {
			log.Print("No leader, proposing self")
			s.ElectMe()
		}
		tick++
	}

}

func (s *Server) HandleVote(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	leader := r.URL.Query().Get("leader")
	log.Printf("Got vote for %s", leader)
	if s.ProposedLeader == "" && s.ConfirmedLeader == "" {
		log.Printf("No leader, accepting request from %s", leader)
		s.ProposedLeader = leader
	} else {
		log.Printf("already have a proposed leader, denying request from %s", leader)
		http.Error(w, "No", 405)
	}
}

func (s *Server) HandleConfirm(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	leader := r.URL.Query().Get("leader")
	log.Printf("Got confirm for %s", leader)
	if s.ConfirmedLeader == "" && s.ProposedLeader == leader {
		log.Printf("Confirming %s as leader", leader)
		s.State = Follower
		s.ConfirmedLeader = leader
	} else {
		log.Printf("already have a leader, denying request from %s", leader)
		http.Error(w, "No", 405)
	}
}

func main() {
	self := flag.String("self", "", "My hostname")
	nodes := flag.String("nodes", "", "the other nodes, separated by commas")
	flag.Parse()

	hosts := strings.Split(*nodes, ",")
	for i, host := range hosts {
		hosts[i] = fmt.Sprintf("http://%s:3000", host)
	}

	server := New(
		fmt.Sprintf("http://%s:3000", *self),
		hosts,
	)
	log.Printf("Starting event loop %s", server.Self)

	http.HandleFunc("/vote", server.HandleVote)
	http.HandleFunc("/confirm", server.HandleConfirm)
	go func() {
		log.Print("starting server")
		log.Fatal(http.ListenAndServe(":3000", nil))
	}()
	server.EventLoop()
}

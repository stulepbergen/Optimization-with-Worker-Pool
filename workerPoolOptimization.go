package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

const (
	tasksCount, usersCount = 100, 20
)

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	rand.Seed(time.Now().Unix())

	tasks := make(chan int, tasksCount)
	users := make(chan User, tasksCount)
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < usersCount; i++ {
		go worker(i+1, tasks, users, &wg1)
		go saver(i+1, users, &wg2)

		wg1.Add(1)
		wg2.Add(1)
	}

	for t := 0; t < 100; t++ {
		tasks <- t
	}
	close(tasks)

	wg1.Wait()
	close(users)
	wg2.Wait()

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func saveUserInfo(user User) {
	fmt.Printf("saving user %d\n", user.id)
	fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

	filename := fmt.Sprintf("users/uid%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(user.getActivityInfo())
	time.Sleep(time.Second)
}

func worker(id int, tasks <-chan int, users chan<- User, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tasks {
		users <- User{
			id:    t + 1,
			email: fmt.Sprintf("user%d@company.com", t+1),
			logs:  generateLogs(rand.Intn(1000)),
		}
		fmt.Printf("worker #%d finished\n", id)
	}
}

func saver(id int, users chan User, wg *sync.WaitGroup) {
	defer wg.Done()
	for user := range users {
		saveUserInfo(user)
	}
}

func generateUsers(count int) []User {
	users := make([]User, count)

	for i := 0; i < count; i++ {
		users[i] = User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@company.com", i+1),
			logs:  generateLogs(rand.Intn(1000)),
		}
		fmt.Printf("generated user %d\n", i+1)
		time.Sleep(time.Millisecond * 100)
	}

	return users
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}

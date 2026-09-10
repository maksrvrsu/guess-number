package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

type GameResult struct {
	Date     string `json:"date"`
	Outcome  string `json:"outcome"`
	Attempts int    `json:"attempts"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		playGame(reader)

		fmt.Print("\nХотите сыграть ещё раз?")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "д" && answer != "y" && answer != "yes" && answer != "да" {
			fmt.Println("Спасибо за игру! До встречи.")
			break
		}
	}
}

func playGame(reader *bufio.Reader) {
	maxNumber, maxAttempts := selectDifficulty(reader)
	secretNumber := generateSecretNumber(maxNumber)

	fmt.Printf("\n%sИгра 'Угадай число' - от 1 до %d началась!%s\n", colorCyan, maxNumber, colorReset)
	fmt.Printf("%sУгадайте число за %d попыток!%s\n\n", colorCyan, maxAttempts, colorReset)

	var history []int
	attemptsUsed := 0
	won := false

	for attemptsUsed < maxAttempts {
		currentAttempt := attemptsUsed + 1

		if len(history) > 0 {
			fmt.Printf("Предыдущие попытки: %v\n", history)
		}

		prompt := fmt.Sprintf("%sПопытка #%d%s - Введите число: ", colorYellow, currentAttempt, colorReset)
		guess, ok := readValidInt(reader, prompt)
		if !ok {
			continue
		}

		history = append(history, guess)
		attemptsUsed++

		if guess == secretNumber {
			fmt.Printf("%sВы угадали!🙌%s\n", colorGreen, colorReset)
			won = true
			break
		}

		if guess < secretNumber {
			fmt.Println("Секретное число больше 👆")
		} else {
			fmt.Println("Секретное число меньше 👇")
		}

		printDistanceHint(guess, secretNumber)
		fmt.Println()
	}
	outcome := "победа"
	if !won {
		outcome = "проигрыш"
		fmt.Printf("%sВы проиграли! 🌚%s\n", colorRed, colorReset)
		fmt.Printf("Секретное число было: %d\n", secretNumber)
	}

	saveResult(outcome, attemptsUsed)
	fmt.Println("Игра закончена!")
}

func selectDifficulty(reader *bufio.Reader) (maxNumber int, attempts int) {
	fmt.Println("\nВыберите уровень сложности:")
	fmt.Println("1. Easy (1-50, 15 попыток)")
	fmt.Println("2. Medium (1-100, 10 попыток)")
	fmt.Println("3. Hard (1-200, 5 попыток)")

	for {
		fmt.Print("Введите 1, 2 или 3: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			return 50, 15
		case "2":
			return 100, 10
		case "3":
			return 200, 5
		default:
			fmt.Println("Некорректный выбор, попробуйте снова.")
		}
	}
}

func generateSecretNumber(max int) int {
	return rand.Intn(max) + 1
}

func readValidInt(reader *bufio.Reader, prompt string) (int, bool) {
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения ввода.")
		return 0, false
	}

	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		fmt.Println("Ввод не может быть пустым. Попытка не списана.")
		return 0, false
	}

	val, err := strconv.Atoi(trimmed)
	if err != nil {
		fmt.Println("Ошибка: введено не число! Попытка не списана.")
		return 0, false
	}

	return val, true
}

func printDistanceHint(guess, secret int) {
	diff := int(math.Abs(float64(guess - secret)))
	if diff <= 5 {
		fmt.Println("Подсказка: 🔥 'Горячо'")
	} else if diff <= 15 {
		fmt.Println("Подсказка: 😃 'Тепло'")
	} else {
		fmt.Println("Подсказка: ❄️ 'Холодно'")
	}
}

func saveResult(outcome string, attempts int) {
	filename := "results.json"
	var results []GameResult

	fileData, err := os.ReadFile(filename)
	if err == nil {
		_ = json.Unmarshal(fileData, &results)
	}

	newEntry := GameResult{
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Outcome:  outcome,
		Attempts: attempts,
	}

	results = append(results, newEntry)

	bytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("Не удалось сохранить статистику.")
		return
	}

	_ = os.WriteFile(filename, bytes, 0644)
}

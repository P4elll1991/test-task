package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func sendChallenge(destination string) {
	url := "http://localhost:8000/sendChallenge"
	// Создание буфера для формы.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err := writer.WriteField("destination", destination)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Завершение формы.
	writer.Close()

	// Создание POST-запроса с формой и отправка его на сервер.
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()

	// Чтение ответа от сервера.
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Ошибка чтения ответа: %v\n", err)
		return
	}

	fmt.Printf("Ответ от сервера:\n%s\n", responseBody)
}

func sendMessage(destination string) {
	url := "http://localhost:8000/sendMessage"

	data := make([]byte, 320)
	rand.Read(data)

	// Создание буфера для формы.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err := writer.WriteField("destination", destination)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Открытие файла для чтения.
	file, err := os.Open("./test-client/test")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Создание части файла формы и его запись в запрос.
	part, err := writer.CreateFormFile("data", "./test-client/test")
	if err != nil {

	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writer.WriteField("destination", string(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Завершение формы.
	writer.Close()

	// Создание POST-запроса с формой и отправка его на сервер.
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()

	// Чтение ответа от сервера.
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Ошибка чтения ответа: %v\n", err)
		return
	}

	fmt.Printf("Ответ от сервера:\n%s\n", responseBody)
}

func main() {
	fmt.Println("УСПЕШНЫЙ вызов sendChallenge с валидным destination")
	sendChallenge("localhost:9000")
	fmt.Println("УСПЕШНЫЙ вызов sendMessage с валидным destination")
	sendMessage("localhost:9000")
	fmt.Println("УСПЕШНЫЙ вызов sendChallenge с не валидным destination")
	sendChallenge("localhost:9001")
	fmt.Println("УСПЕШНЫЙ вызов sendMessage с не валидным destination")
	sendMessage("localhost:9001")
	fmt.Println("НЕ УСПЕШНЫЙ вызов sendChallenge с валидным destination")
	sendChallenge("localhost:9002")
	fmt.Println("НЕ УСПЕШНЫЙ вызов sendMessage с валидным destination")
	sendMessage("localhost:9002")
	fmt.Println("НЕ УСПЕШНЫЙ вызов sendChallenge с не валидным destination")
	sendChallenge("localhost:9003")
	fmt.Println("НЕ УСПЕШНЫЙ вызов sendMessage с не валидным destination")
	sendMessage("localhost:9003")
}

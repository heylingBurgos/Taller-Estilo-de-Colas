package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	broker = "localhost:9092"
	topic  = "DigiturnoE"
)

func main() {

	router := mux.NewRouter()

	// Configurar el middleware CORS
	handler := cors.Default().Handler(router)

	// Agregar una ruta para recibir los datos del usuario desde Vue
	router.HandleFunc("http://localhost:8081/enviar-turno", RecibirTurno)

	// Define your routes here using router.HandleFunc()
	http.ListenAndServe(":8080", handler)
}

func RecibirTurno(w http.ResponseWriter, r *http.Request) {
	// Leer los datos enviados desde Vue
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	// Responder al cliente de Vue
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))

	// Configuración del cliente de Kafka
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

	if r.Method == "POST" {

		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		name := r.FormValue("name")
		cellphone, _ := strconv.ParseInt(r.FormValue("cellphone"), 10, 64)

		// Creación del cliente de Kafka
		producer, err := sarama.NewSyncProducer([]string{broker}, config)
		if err != nil {
			log.Fatalf("Error al crear el cliente de Kafka: %v", err)
		}
		defer func() {
			if err := producer.Close(); err != nil {
				log.Fatalf("Error al cerrar el cliente de Kafka: %v", err)
			}
		}()

		// Creación del mensaje de Kafka
		message := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("%d;%s;%d", id, name, cellphone)),
		}

		// Envío del mensaje a Kafka
		partition, offset, err := producer.SendMessage(message)
		if err != nil {
			log.Printf("Error al enviar el mensaje a Kafka: %v", err)
		} else {
			log.Printf("Mensaje enviado a la partición %d con offset %d", partition, offset)

			// Agregar offset al archivo
			f, err := os.OpenFile("offsets.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				log.Printf("Error al abrir el archivo de offsets: %v", err)
			} else {
				defer f.Close()
				if _, err := f.WriteString(fmt.Sprintf("%d;%d;%d;%s;%d\n", partition, offset, id, name, cellphone)); err != nil {
					log.Printf("Error al escribir el offset en el archivo: %v", err)
				}
			}
		}

		// Creación del cliente de Kafka
		consumer, err := sarama.NewConsumer([]string{broker}, config)
		if err != nil {
			log.Fatalf("Error al crear el consumidor de Kafka: %v", err)
		}
		defer func() {
			if err := consumer.Close(); err != nil {
				log.Fatalf("Error al cerrar el consumidor de Kafka: %v", err)
			}
		}()

		// Asignación de la partición al consumidor
		partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error al asignar la partición al consumidor de Kafka: %v", err)
		}
		defer func() {
			if err := partitionConsumer.Close(); err != nil {
				log.Fatalf("Error al cerrar el consumidor de la partición: %v", err)
			}
		}()

		// Configuración de señales para el cierre del programa
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		// Ciclo principal del programa
		for {

			// Comprobación de señales
			select {
			case msg := <-partitionConsumer.Messages():
				// Procesar el mensaje recibido
				log.Printf("Mensaje recibido: %s", string(msg.Value))
				parts := strings.Split(string(msg.Value), ";")
				id, _ := strconv.ParseInt(parts[0], 10, 64)
				name := parts[1]
				cellphone, _ := strconv.ParseInt(parts[2], 10, 64)

				// Enviar los datos de vuelta a Vue en un mensaje
				w.Header().Set("Content-Type", "application/json")
				resp := struct {
					ID        int64  `json:"id"`
					Name      string `json:"name"`
					Cellphone int64  `json:"cellphone"`
				}{
					ID:        id,
					Name:      name,
					Cellphone: cellphone,
				}
				json.NewEncoder(w).Encode(resp)

				f, err := os.OpenFile("recived.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					log.Printf("Error al abrir el archivo de offsets: %v", err)
				} else {
					defer f.Close()
					if _, err := f.WriteString(fmt.Sprintf("Mensaje recibido: %s\n", string(msg.Value))); err != nil {
						log.Printf("Error al escribir el offset en el archivo: %v", err)
					}
				}
			case err := <-partitionConsumer.Errors():
				log.Printf("Error al recibir el mensaje: %v", err)
			case <-signals:
				log.Println("Cerrando el programa...")
				partitionConsumer.AsyncClose()
				return

			}
		}
		// ...
	} else if r.Method == "OPTIONS" {
		// Handle OPTIONS request (preflight)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

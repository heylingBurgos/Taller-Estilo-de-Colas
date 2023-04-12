package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	broker  = "localhost:9092"
	GroupID = "GrupoDigiturno"
	topics  = []string{"eltopico"}
)

type Datos struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Cellphone int64  `json:"cellphone"`
}

var m = 0

type Usuario struct {
	ID      int64  `json:"id"`
	Nombre  string `json:"name"`
	Celular int64  `json:"cellphone"`
	Turno   int    `json:"turn"`
}

var usuarios []*Usuario

func main() {
	router := mux.NewRouter()

	// Configurar el middleware CORS
	handler := cors.Default().Handler(router)

	//Recibir la lista de turnos
	router.HandleFunc("/turnos", getTurnos)

	// Agregar una ruta para recibir los datos del usuario desde Vue
	router.HandleFunc("/", RecibirTurno)

	// Define your routes here using router.HandleFunc()
	err := http.ListenAndServe(":5000", handler)
	if err != nil {
		log.Fatal(err)
	}

}

func getTurnos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	fi, err := os.OpenFile("../logs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error al abrir el archivo de offsets: %v", err)
	} else {
		defer fi.Close()
		if _, err := fi.WriteString("Desde Golang recibiendo petición GET lista de usuarios. \n"); err != nil {
			log.Printf("Error al escribir el offset en el archivo: %v", err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usuarios)
}

func RecibirTurno(w http.ResponseWriter, r *http.Request) {
	// Leer los datos enviados desde Vue
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	if r.Method == "POST" {

		// Configuración del cliente de Kafka
		config := sarama.NewConfig()
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Retry.Max = 5
		config.Producer.Return.Successes = true
		config.Consumer.Group.Session.Timeout = 10 * time.Second
		config.Consumer.Group.Heartbeat.Interval = 3 * time.Second

		decoder := json.NewDecoder(r.Body)
		var datos Datos
		err := decoder.Decode(&datos)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Body.Close()
		// Utilizar los datos recibidos
		log.Print("Datos  ", datos.Id, datos.Name, datos.Cellphone)

		// Creación del cliente de Kafka
		producer, err := sarama.NewSyncProducer([]string{broker}, config)
		if err != nil {
			log.Fatalf("Error al crear el productor de Kafka: %v", err)
		}
		defer func() {
			if err := producer.Close(); err != nil {
				log.Fatalf("Error al cerrar el productor de Kafka: %v", err)
			}
		}()

		// Creación del mensaje de Kafka
		message := &sarama.ProducerMessage{
			Topic: topics[0],
			Value: sarama.StringEncoder(fmt.Sprintf("%d;%s;%d", datos.Id, datos.Name, datos.Cellphone)),
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
				if _, err := f.WriteString(fmt.Sprintf("%d;%d;%d;%s;%d\n", partition, offset, datos.Id, datos.Name, datos.Cellphone)); err != nil {
					log.Printf("Error al escribir el offset en el archivo: %v", err)
				}
			}
		}

		// Creación del cliente de Kafka ---------------------------------------------------------
		consumer, err := sarama.NewConsumerGroup([]string{broker}, "digiturno", config)
		if err != nil {
			log.Fatalf("Error al crear el consumidor de Kafka: %v", err)
		}
		defer func() {
			if err := consumer.Close(); err != nil {
				log.Fatalf("Error al cerrar el consumidor de Kafka: %v", err)
			}
		}()

		// Configuración de señales para el cierre del programa
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		//Canal para recibir los mensajes consumidos
		messages := make(chan *sarama.ConsumerMessage, 256)

		// Manejador del grupo de consumidores
		handler := &consumerHandler{
			messages: messages,
		}

		// Inicio del grupo de consumidores
		go func() {
			for {
				if err := consumer.Consume(context.Background(), topics, handler); err != nil {
					log.Fatalf("Error al consumir mensajes: %v", err)
				}
				if handler.ctx.Err() != nil {
					// Si el contexto está cancelado, salir del bucle
					return
				}

			}
		}()

		// Consumo de los mensajes
	consumingLoop:
		for {
			select {
			case msg := <-messages:
				// Procesar el mensaje recibido
				log.Printf("Mensaje recibido: %s", string(msg.Value))
				parts := strings.Split(string(msg.Value), ";")
				id, _ := strconv.ParseInt(parts[0], 10, 64)
				name := parts[1]
				cellphone, _ := strconv.ParseInt(parts[2], 10, 64)
				m = m + 1

				usuario := &Usuario{
					ID:      id,
					Nombre:  name,
					Celular: cellphone,
					Turno:   len(usuarios) + 1, // Asignar el siguiente número de turno disponible
				}
				usuarios = append(usuarios, usuario)

				resp := struct {
					ID      int64  `json:"id"`
					Nombre  string `json:"name"`
					Celular int64  `json:"cellphone"`
					Turno   int    `json:"turn"`
				}{
					ID:      usuario.ID,
					Nombre:  usuario.Nombre,
					Celular: usuario.Celular,
					Turno:   usuario.Turno,
				}

				respJSON, err := json.Marshal(resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// Enviar los datos de vuelta a Vue en un mensaje
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(respJSON)
				handler.session.MarkMessage(msg, "")
				log.Printf("Mensaje recibido: %s", string(respJSON))

				f, err := os.OpenFile("recived.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					log.Printf("Error al abrir el archivo de offsets: %v", err)
				} else {
					defer f.Close()
					if _, err := f.WriteString(fmt.Sprintf("Mensaje recibido: %s%s%d\n", string(msg.Value), ";", usuario.Turno)); err != nil {
						log.Printf("Error al escribir el offset en el archivo: %v", err)
					}
				}
				return

			case <-signals:
				log.Print("Señal recibida, deteniendo la aplicación...")
				handler.cancel()
				break consumingLoop
			}
		}

	} else if r.Method == "OPTIONS" {
		// Handle OPTIONS request (preflight)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

// Manejador del grupo de consumidores
type consumerHandler struct {
	messages chan *sarama.ConsumerMessage
	session  sarama.ConsumerGroupSession
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func (h *consumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	h.ctx, h.cancel = context.WithCancel(context.Background())
	h.session = session
	return nil
}

func (h *consumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	h.cancel()
	h.wg.Wait()
	return nil
}

func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	h.wg.Add(1)
	defer h.wg.Done()
	for msg := range claim.Messages() {
		h.messages <- msg
	}
	return nil
}

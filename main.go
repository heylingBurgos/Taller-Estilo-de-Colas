package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

const (
	broker = "localhost:9092"
	topic  = "DigiturnoE"
)

func main() {
	// Configuración del cliente de Kafka
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

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
		// Lectura del ID y nombre desde la entrada estándar
		var id int64
		var nombre string
		fmt.Print("Ingrese el ID: ")
		fmt.Scan(&id)
		fmt.Print("Ingrese el nombre: ")
		fmt.Scan(&nombre)

		// Creación del mensaje de Kafka
		message := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("%d;%s", id, nombre)),
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
				if _, err := f.WriteString(fmt.Sprintf("%d;%d;%d;%s\n", partition, offset, id, nombre)); err != nil {
					log.Printf("Error al escribir el offset en el archivo: %v", err)
				}
			}
		}

		// Espera de 1 segundo antes de continuar
		time.Sleep(1 * time.Second)

		// Comprobación de señales para el cierre del programa
		select {
		case msg := <-partitionConsumer.Messages():
			// Procesar el mensaje recibido
			log.Printf("Mensaje recibido: %s", string(msg.Value))
			// // Actualizar el archivo de offsets
			// f, err := os.OpenFile("offsets.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
			// if err != nil {
			// 	log.Printf("Error al abrir el archivo de offsets: %v", err)
			// } else {
			// 	defer f.Close()
			// 	partitions, err := consumer.Partitions(topic)
			// 	if err != nil {
			// 		log.Fatalf("Error al obtener las particiones del tópico %s: %v", topic, err)
			// 	}

			// 	// Obtener los offsets más recientes de cada partición
			// 	for _, p := range partitions {

			// 		if err != nil {
			// 			log.Fatalf("Error al obtener el offset más reciente de la partición %d: %v", p, err)
			// 		}
			// 		if _, err := f.WriteString(fmt.Sprintf("%d;%d;%s\n", p, offset, time.Now().Format(time.RFC3339))); err != nil {
			// 			log.Printf("Error al escribir el offset en el archivo: %v", err)
			// 		}
			// 	}
			// 	if err := f.Sync(); err != nil {
			// 		log.Printf("Error al sincronizar el archivo de offsets: %v", err)
			// 	}
			// }

		case err := <-partitionConsumer.Errors():
			log.Printf("Error al recibir el mensaje: %v", err)
		case <-signals:
			log.Println("Cerrando el programa...")
			return
		default:
		}
	}
}

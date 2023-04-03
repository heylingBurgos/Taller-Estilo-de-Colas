<template>
  <div class="User">
    <img alt="Medi plus" src="../assets/Logo.png" />
    <h1>BuroMedi</h1>
    <label for="id">Ingrese su cédula:</label>
    <input type="number" v-model="id" />
    <br />
    <br />
    <label for="name">Ingrese su nombre:</label>
    <input type="text" v-model="name" />
    <br />
    <br />
    <label for="cellphone">Ingrese su celular:</label>
    <input type="number" v-model="cellphone" />
    <br />
    <br />
    <br />
    <div v-if="datosTurno != null">
      <p>Bienvenido usuario {{ datosTurno.name }} su turno es</p>
      <h1>{{ datosTurno.turn }}</h1>
    </div>
    <br />
    <br />
    <button id="button" @click="turno()">registrar</button>
  </div>
</template>
<!-- Digiturno-->
<script>
import axios from "axios";
export default {
  name: "DigiturnoUser",
  data() {
    return {
      id: null, // Define la propiedad "id"
      name: null,
      cellphone: null,
      datosTurno: null,
    };
  },
  methods: {
    turno() {
      axios
        .post("http://localhost:4000/", {
          id: this.id,
          name: this.name,
          cellphone: this.cellphone,
        },
        console.log("Solicitando petición POST")
        )
        .then((response) => {
          console.log("Respuesta dada por el backend y recibida por DigiturnoUser: "+response.data);
          this.datosTurno = response.data;
        },
        )
        .catch((error) => {
          console.log(error);
        });
        this.id = null;
        this.name = null;
        this.cellphone = null;
    }
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
label {
  font-family: Avenir, Helvetica, Arial, sans-serif;
}

input {
  margin-left: 2%;
  border: none;
  padding: 1rem;
  border-radius: 1rem;
  background: #e8e8e8;
  box-shadow: 20px 20px 60px #ffff, -20px -20px 60px #ffffff;
  transition: 0.3s;
}

input:focus {
  outline-color: #e8e8e8;
  background: #e8e8e8;
  transition: 0.3s;
}

button {
  appearance: button;
  background-color: #1899d6;
  border: solid transparent;
  border-radius: 16px;
  border-width: 0 0 4px;
  box-sizing: border-box;
  color: #ffffff;
  cursor: pointer;
  display: inline-block;
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.8px;
  line-height: 20px;
  margin: 0;
  outline: none;
  overflow: visible;
  padding: 13px 19px;
  text-align: center;
  touch-action: manipulation;
  transform: translateZ(0);
  transition: filter 0.2s;
  user-select: none;
  -webkit-user-select: none;
  vertical-align: middle;
  white-space: nowrap;
}

button:after {
  background-clip: padding-box;
  background-color: #1cb0f6;
  border: solid transparent;
  border-radius: 16px;
  border-width: 0 0 4px;
  bottom: -4px;
  content: "";
  left: 0;
  position: absolute;
  right: 0;
  top: 0;
  z-index: -1;
}

button:main,
button:focus {
  user-select: auto;
}

button:hover:not(:disabled) {
  filter: brightness(1.1);
}

button:disabled {
  cursor: auto;
}

button:active:after {
  border-width: 0 0 0px;
}

button:active {
  padding-bottom: 10px;
}

div {
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
  padding: 20px;
  display: inline-block;
  max-width: 100%;
}

img {
  width: 35%;
  height: auto;
}

#routerlink {
  text-decoration: none;
  font-style: none;
}
</style>

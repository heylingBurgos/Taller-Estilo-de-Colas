<template>
  <div class="divi">
    <img alt="Medi plus" src="../assets/Logo.png">
    <h1>BuroMedi</h1>
    <div v-if="pacientesMostrados">
      <p v-if="mayor">No hay más turnos</p>
      <table>
      <thead>
        <tr>
          <th>Nombre</th>
          <th>Turno</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="paciente in pacientesMostrados" :key="paciente.id">
          <td>{{ paciente.name }}</td>
          <td>{{ paciente.turn }}</td>
        </tr>
      </tbody>
    </table>
    </div>
    <div v-else>
      <p>No hay nada que mostrar</p>
    </div>
    
    <br>
    <br>
    <button @click="mostrarSiguiente">Mostrar</button>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      pacientes: [],
      indiceMostrado: 0,
      pacientesPorPagina: 1,
      mayor: false
    };
  },
  computed: {
    pacientesMostrados() {
    if (this.pacientes) {
      return this.pacientes.slice(
        this.indiceMostrado,
        this.indiceMostrado + this.pacientesPorPagina
      );
    } else {
      return [];
    }
  }
  },
  methods: {
    mostrarSiguiente() {
      this.indiceMostrado += this.pacientesPorPagina;
      if (this.indiceMostrado >= this.pacientes.length) {
        this.mayor = true
      }else{
        this.mayor = false
      }
    },
    obtenerPacientes() {
      axios.get('http://localhost:4000/turn')
        .then(response => {
          console.log("Respuesta dada por el bakend y recibida por DigiturnoView: "+response.data);
          this.pacientes = response.data;
        })
        .catch(error => {
          console.error('Error al recibir los pacientes:', error);
        });
    }
  },
  mounted() {
    this.obtenerPacientes();
    setInterval(this.obtenerPacientes, 2000);
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
label{
  font-family: Avenir, Helvetica, Arial, sans-serif;  
}

input {
  margin-left: 2%;
  border: none;
  padding: 1rem;
  border-radius: 1rem;
  background: #e8e8e8;
  box-shadow: 20px 20px 60px #ffff,
      -20px -20px 60px #ffffff;
  transition: 0.3s;
}

input:focus {
  outline-color: #e8e8e8;
  background: #e8e8e8;
  transition: 0.3s;
}

button {
  appearance: button;
  background-color: #1899D6;
  border: solid transparent;
  border-radius: 16px;
  border-width: 0 0 4px;
  box-sizing: border-box;
  color: #FFFFFF;
  cursor: pointer;
  display: inline-block;
  font-size: 15px;
  font-weight: 700;
  letter-spacing: .8px;
  line-height: 20px;
  margin: 0;
  outline: none;
  overflow: visible;
  padding: 13px 19px;
  text-align: center;
  touch-action: manipulation;
  transform: translateZ(0);
  transition: filter .2s;
  user-select: none;
  -webkit-user-select: none;
  vertical-align: middle;
  white-space: nowrap;
}

button:after {
  background-clip: padding-box;
  background-color: #1CB0F6;
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

button:main, button:focus {
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

.divi{
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

table {
  width: 50%;
  margin: auto;
  text-align: center;
  border-collapse: collapse;
}

th, td {
  padding: 10px;
  border: 1px solid #ddd;
}

th {
  background-color: #f2f2f2;
}

</style>
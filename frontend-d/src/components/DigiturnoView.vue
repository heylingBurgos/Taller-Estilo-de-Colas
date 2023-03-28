<template>
    <div class="User">
      <img alt="Medi plus" src="../assets/Logo.png">
      <h1>BuroMedi</h1>
      <p>Bienvenido usuario {{name1}} es su turno</p>
      <br>
      <br>
      <br>
      <h1>1</h1>
      <button id="button">Regresar</button>
    </div>
</template>

<script>
export default {
  props: {
    id1: Number,
    name1: String,
    cellphone1: Number
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

div{
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

</style>
<template>
  <div>
    <Barra/>
    <ListaTurnos/>
  </div>
</template>

<script>
import Barra from './components/Barra.vue'
import ListaTurnos from './components/ListaTurnos.vue'

export default {
  name: 'App',
  components: {
    Barra,
    ListaTurnos
  }
}
</script>

